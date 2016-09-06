package deploy

import (
	"app/backend/common/util/Placeholder"
	mylog "app/backend/common/util/log"
	"app/backend/common/util/session"
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
)


const (

	ROLLING_TYPE = myoption.ROLLINGUPGRADE
	ROLLING_VERBE = "POST"
	ROLLING_URL = "/api/v1/organization/<orgId>/deployments/<deploymentName>/rolling"
	ROLLING_MAXUNAVAILABLE = 2
	ROLLING_MAXSURGE = 2

)

type RollingDeployController struct {
	*iris.Context
	k8sClients *client.Client
	apiServers string
	Ye         *myerror.YceError
}

func (rdc *RollingDeployController) WriteBack() {
	rdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("RollingDeployController Response YceError: controller=%p, code=%d, note=%s", rdc, rdc.Ye.Code, myerror.Errors[rdc.Ye.Code].LogMsg)
	rdc.Write(rdc.Ye.String())
}

// Validate Session
func (rdc *RollingDeployController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("RollingDeployment Validate Session error: sessionId=%s, error=%s", sessionId, err)
		rdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("RollingDeployment Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		rdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	mylog.Log.Infof("RollingDeployment validate sessionId: sessionId=%s, orgId=%s", sessionId, orgId)
	return
}

// Get ApiServer by dcId
func (rdc *RollingDeployController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("RollingDeployment getApiServerById QueryDataCenterById Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("RollingDeployment getApiServerByDcId: apiServer=%s, dcId=%d", apiServer, dcId)
	return apiServer
}

// Get ApiServer List for dcIdList
func (rdc *RollingDeployController) getApiServer(dcIdList []int32) {
	// Get ApiServer
	// must be one dcId
	dcId := dcIdList[0]
	apiServer := rdc.getApiServerByDcId(dcId)
	if strings.EqualFold(apiServer, "") {
		mylog.Log.Errorf("RollingDeployController getApiServerList Error")
		return
	}

	//rdc.apiServers = append(rdc.apiServers, apiServer)
	rdc.apiServers = apiServer
	mylog.Log.Infof("RollingDeployment getApiServer: len(rdc.apiServers)=%d", len(rdc.apiServers))
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
		mylog.Log.Errorf("createK8sClient Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	//rdc.k8sClients = append(rdc.k8sClients, c)
	//rdc.apiServers = append(rdc.apiServers, server)
	rdc.k8sClients = c
	mylog.Log.Infof("RollingDeployment CreateK8sClient: c=%p, apiServer=%s", c, server)

	return
}

func (rdc *RollingDeployController) RollingUpdate(namespace, deployment string, rd *deploy.RollingDeployment) (dp *extensions.Deployment) {

	cli := rdc.k8sClients
	dp, err := cli.Extensions().Deployments(namespace).Get(deployment)
	if err != nil {
		mylog.Log.Errorf("RollingDeployment RollingUpdate getDeployment Error: apiServer=%s, namespace=%s, err=%s", rdc.apiServers, namespace, err)
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
		mylog.Log.Errorf("Rolling Update Deployment Error: error=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_ROLLING_DEPLOYMENTS, "")
	}

	mylog.Log.Infof("Rolling Update deployment successfully: namespace=%s, apiserver=%s", namespace, rdc.apiServers)

	return dp

}

// Create Deployment(mysql) and insert it into db
func (rdc *RollingDeployController) createMysqlDeployment(success bool, name, json, reason, dcList, comments string, userId, orgId int32) error {

	uph := placeholder.NewPlaceHolder(ROLLING_URL)
	orgIdString := strconv.Itoa(int(orgId))
	actionUrl := uph.Replace("<orgId>", orgIdString, "<deploymentName>", name)
	actionOp := userId

	dp := mydeployment.NewDeployment(name, ROLLING_VERBE, actionUrl, dcList, reason, json, comments, int32(ROLLING_TYPE), actionOp, int32(1), orgId)
	err := dp.InsertDeployment()
	if err != nil {
		mylog.Log.Errorf("RollingDeployment CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		rdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	mylog.Log.Infof("RollingDeployment CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

//Encode dcIdList to JSON
func (rdc *RollingDeployController) encodeDcIdList(dcIdList []int32) string{
	dcIds := &deploy.DcIdListType{
		DcIdList:dcIdList,
	}

	data, _ := json.Marshal(dcIds)
	mylog.Log.Infof("RollingDeployController encodeDcIdList: dcIdList=%s", string(data))
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
	rdc.validateSession(sessionIdFromClient, orgId)

	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	rd := new(deploy.RollingDeployment)
	rdc.ReadJSON(rd)


	// Get DcIdList
	rdc.getApiServer(rd.DcIdList)
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	// create K8sClient
	rdc.createK8sClients()
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	// RollingUpdate the deployment
	dp := rdc.RollingUpdate(orgName, deploymentName, rd)
	if rdc.Ye != nil {
		rdc.WriteBack()
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
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	rdc.Ye = myerror.NewYceError(myerror.EOK, "")
	rdc.WriteBack()
	// TODO: 成功写回
	mylog.Log.Infoln("Rolling DeploymentController over!")
	return
}
