package deploy

import (
	"app/backend/common/util/Placeholder"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
	mydeployment "app/backend/model/mysql/deployment"
	myoption "app/backend/model/mysql/option"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/util/intstr"
	"strconv"
	"strings"
	yce "app/backend/controller/yce"
)

type RollingDeployController struct {
	yce.Controller
	k8sClients *client.Client
	apiServers string
}

// Get ApiServer by dcId
func (rdc *RollingDeployController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("RollingDeployment getApiServerById QueryDataCenterById Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("RollingDeployment getApiServerByDcId: apiServer=%s, dcId=%d", apiServer, dcId)
	return apiServer
}

// Get ApiServer List for dcIdList
func (rdc *RollingDeployController) getApiServer(dcIdList []int32) {
	// Get ApiServer
	// must be one dcId
	var dcId int32
	if len(dcIdList) > 0 {
		dcId = dcIdList[0]
	} else {
		log.Errorf("RollingDeployController get DcIdList error: len(dcIdList)=%d, error=no value in DcIdList, index out of range", len(dcIdList))
		rdc.Ye = myerror.NewYceError(myerror.EOOM, "")
		return

	}

	apiServer := rdc.getApiServerByDcId(dcId)
	if strings.EqualFold(apiServer, "") {
		log.Errorf("RollingDeployController getApiServerList Error")
		rdc.Ye = myerror.NewYceError(myerror.EOOM, "")
		return
	}

	//rdc.apiServers = append(rdc.apiServers, apiServer)
	rdc.apiServers = apiServer
	log.Infof("RollingDeployment getApiServer: len(rdc.apiServers)=%d", len(rdc.apiServers))
	return
}

// Create k8sClients for every ApiServer
func (rdc *RollingDeployController) createK8sClients() {

	server := rdc.apiServers
	config := &restclient.Config{
		Host: server,
	}

	c, err := client.New(config)
	if err != nil {
		log.Errorf("createK8sClient Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	//rdc.k8sClients = append(rdc.k8sClients, c)
	//rdc.apiServers = append(rdc.apiServers, server)
	rdc.k8sClients = c
	log.Infof("RollingDeployment CreateK8sClient: c=%p, apiServer=%s", c, server)

	return
}

func (rdc *RollingDeployController) rollingUpdate(namespace, deployment string, rd *deploy.RollingDeployment) (dp *extensions.Deployment) {

	cli := rdc.k8sClients
	dp, err := cli.Extensions().Deployments(namespace).Get(deployment)
	if err != nil {
		log.Errorf("RollingDeployment RollingUpdate getDeployment Error: apiServer=%s, namespace=%s, err=%s", rdc.apiServers, namespace, err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
		return
	}

	ds := new(extensions.DeploymentStrategy)
	ds.Type = extensions.RollingUpdateDeploymentStrategyType
	ds.RollingUpdate = new(extensions.RollingUpdateDeployment)
	ds.RollingUpdate.MaxUnavailable = intstr.FromInt(int(ROLLING_MAXUNAVAILABLE))
	ds.RollingUpdate.MaxSurge = intstr.FromInt(int(ROLLING_MAXSURGE))

	dp.Spec.Strategy = *ds
	dp.Spec.Template.Spec.Containers[0].Image = rd.Strategy.Image

	//rolling update interval
	//rd.Strategy.UpdateInterval

	dp.Annotations["userId"] = rd.UserId
	dp.Annotations["image"] = rd.Strategy.Image
	dp.Annotations["kubernetes.io/change-cause"] = rd.Comments

	_, err = cli.Extensions().Deployments(namespace).Update(dp)
	if err != nil {
		log.Errorf("Rolling Update Deployment Error: error=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_ROLLING_DEPLOYMENTS, "")
		return
	}

	log.Infof("Rolling Update deployment successfully: namespace=%s, apiserver=%s", namespace, rdc.apiServers)

	return dp

}

// Create Deployment(mysql) and insert it into db
func (rdc *RollingDeployController) createMysqlDeployment(success bool, name, json, reason, dcList, comments string, userId, orgId int32) error {

	uph := placeholder.NewPlaceHolder(ROLLING_URL)
	orgIdString := strconv.Itoa(int(orgId))
	actionUrl := uph.Replace("<orgId>", orgIdString, "<deploymentName>", name)
	actionOp := userId
	log.Debugf("RollingDeployController createMySQLDeployment: actionOp=%d", actionOp)
	dp := mydeployment.NewDeployment(name, ROLLING_VERBE, actionUrl, dcList, reason, json, comments, int32(ROLLING_TYPE), actionOp, int32(1), orgId)
	err := dp.InsertDeployment()
	if err != nil {
		log.Errorf("RollingDeployment CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		rdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	log.Infof("RollingDeployment CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

//Encode dcIdList to JSON
func (rdc *RollingDeployController) encodeDcIdList(dcIdList []int32) string{
	dcIds := &deploy.DcIdListType{
		DcIdList:dcIdList,
	}

	data, _ := json.Marshal(dcIds)
	log.Infof("RollingDeployController encodeDcIdList: dcIdList=%s", string(data))
	return string(data)
}

func (rdc RollingDeployController) Post() {
	orgId := rdc.Param("orgId")
	deploymentName := rdc.Param("deploymentName")
	org := new(myorganization.Organization)


	// Get orgName by orgId
	orgIdInt, _ := strconv.Atoi(orgId)
	org.QueryOrganizationById(int32(orgIdInt))
	orgName := org.Name

	sessionIdFromClient := rdc.RequestHeader("Authorization")
	log.Debugf("RolllingDeployController Params: sessionId=%s, orgId=%s, orgName=%s", sessionIdFromClient, orgId, orgName)

	rdc.ValidateSession(sessionIdFromClient, orgId)
	if rdc.CheckError() {
		return
	}

	rd := new(deploy.RollingDeployment)
	err := rdc.ReadJSON(rd)
	if err != nil {
		log.Errorf("RollingUpDeployController ReadJSON Error: error=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}

	if rdc.CheckError() {
		return
	}

	// Get DcIdList
	rdc.getApiServer(rd.DcIdList)
	if rdc.CheckError() {
		return
	}

	// create K8sClient
	rdc.createK8sClients()
	if rdc.CheckError() {
		return
	}

	// RollingUpdate the deployment
	dp := rdc.rollingUpdate(orgName, deploymentName, rd)
	if rdc.CheckError() {
		return
	}

	// Encode cd.DcIdList to json
	//dcl, _ := json.Marshal(rd.DcIdList)
	dcIdList := rdc.encodeDcIdList(rd.DcIdList)

	// Encode k8s.deployment to json
	kd, _ := json.Marshal(dp)

	oId, _ := strconv.Atoi(orgId)


	// Insert into mysql.Deployment
	userId, _ :=  strconv.Atoi(rd.UserId)
	rdc.createMysqlDeployment(true, rd.AppName,  string(kd), "", dcIdList, rd.Comments, int32(userId), int32(oId))
	if rdc.CheckError() {
		return
	}

	rdc.WriteOk("")
	log.Infoln("Rolling DeploymentController over!")

	return
}
