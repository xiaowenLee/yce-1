package deploy

import (
	"app/backend/common/util/Placeholder"
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mydeployment "app/backend/model/mysql/deployment"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/util/intstr"
	"strconv"
)

type RollingDeploymentController struct {
	yce.Controller
	k8sClients *client.Client
	apiServers string
}

func (rdc *RollingDeploymentController) rollingUpdate(namespace, deployment string, rd *deploy.RollingDeployment) (dp *extensions.Deployment) {

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
func (rdc *RollingDeploymentController) createMysqlDeployment(success bool, name, json, reason, dcList, comments string, userId, orgId int32) error {

	uph := placeholder.NewPlaceHolder(ROLLING_URL)
	orgIdString := strconv.Itoa(int(orgId))
	actionUrl := uph.Replace("<orgId>", orgIdString, "<deploymentName>", name)
	actionOp := userId
	log.Debugf("RollingDeploymentController createMySQLDeployment: actionOp=%d", actionOp)
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

func (rdc RollingDeploymentController) Post() {
	orgId := rdc.Param("orgId")
	deploymentName := rdc.Param("deploymentName")

	sessionIdFromClient := rdc.RequestHeader("Authorization")

	// Get orgName by orgId
	orgName, ye := yceutils.GetOrgNameByOrgId(orgId)
	if ye != nil {
		rdc.Ye = ye
	}
	if rdc.CheckError() {
		return
	}
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
	if len(rd.DcIdList) == 0 {
		rdc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	if rdc.CheckError() {
		return
	}

	// Get ApiServer
	dcId := rd.DcIdList[0]
	rdc.apiServers, rdc.Ye = yceutils.GetApiServerByDcId(dcId)
	if rdc.CheckError() {
		return
	}

	// Create K8sClient
	rdc.k8sClients, rdc.Ye = yceutils.CreateK8sClient(rdc.apiServers)
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
	dcIdList, ye := yceutils.EncodeDcIdList(rd.DcIdList)
	if ye != nil {
		rdc.Ye = ye
	}
	if rdc.CheckError() {
		return
	}

	// Encode k8s.deployment to json
	kd, _ := json.Marshal(dp)
	oId, _ := strconv.Atoi(orgId)

	// Insert into mysql.Deployment
	userId, _ := strconv.Atoi(rd.UserId)
	rdc.createMysqlDeployment(true, rd.AppName, string(kd), "", dcIdList, rd.Comments, int32(userId), int32(oId))

	if rdc.CheckError() {
		return
	}

	rdc.WriteOk("")
	log.Infoln("Rolling DeploymentController over!")

	return
}
