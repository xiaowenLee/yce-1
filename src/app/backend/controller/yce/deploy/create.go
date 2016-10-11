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
	"strconv"
)

type CreateDeploymentController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
}

// Publish k8s.Deployment to every datacenter which in dcIdList
func (cdc *CreateDeploymentController) createDeployment(namespace string, deployment *extensions.Deployment) {

	// Foreach every k8sClient to create deployment
	for index, cli := range cdc.k8sClients {
		_, err := cli.Extensions().Deployments(namespace).Create(deployment)
		if err != nil {
			log.Errorf("createDeployment Error: apiServer=%s, namespace=%s, err=%s",
				cdc.apiServers[index], namespace, err)
			cdc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_DEPLOYMENT, "")
			return
		}

		log.Infof("Create deployment successfully: namespace=%s, apiserver=%s", namespace, cdc.apiServers[index])
	}
}

// Create Deployment(mysql) and insert it into db
func (cdc *CreateDeploymentController) createMysqlDeployment(success bool, name, userId, json, reason, dcList string, orgId int32) error {

	uph := placeholder.NewPlaceHolder(CREATE_URL)
	orgIdString := strconv.Itoa(int(orgId))
	actionUrl := uph.Replace("<orgId>", orgIdString, "<userId>", userId)
	actionOp, _ := strconv.Atoi(userId)
	log.Debugf("CreateDeploymentController createMySQLDeployment: actionUrl=%s, actionOp=%d", actionUrl, actionOp)

	dp := mydeployment.NewDeployment(name, CREATE_VERBE, actionUrl, dcList, reason, json, "Create a Deployment", CREATE_TYPE, int32(actionOp), int32(1), orgId)
	err := dp.InsertDeployment()
	if err != nil {
		log.Errorf("CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		cdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	log.Infof("CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

// POST /api/v1/organizations/{orgId}/users/{userId}/deployments
func (cdc CreateDeploymentController) Post() {
	sessionIdFromClient := cdc.RequestHeader("Authorization")
	orgId := cdc.Param("orgId")
	userId := cdc.Param("userId")

	log.Debugf("CreateDeploymentController get Params:  sessionIdFromClient=%s, orgId=%s, userId=%s", sessionIdFromClient, orgId, userId)

	// Validate OrgId error
	cdc.ValidateSession(sessionIdFromClient, orgId)
	if cdc.CheckError() {
		return
	}

	// Parse data: deploy.CreateDeployment
	cd := new(deploy.CreateDeployment)
	err := cdc.ReadJSON(cd)
	if err != nil {
		log.Errorf("CreateDeploymentController ReadJSON error: error=%s", err)
		cdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}

	if cdc.CheckError() {
		return
	}

	log.Infof("CreateDeploymentController ReadJSON success: cd=%p", cd)

	// Get DcIdList
	cdc.apiServers, cdc.Ye = yceutils.GetApiServerList(cd.DcIdList)
	if cdc.CheckError() {
		return
	}

	// Create k8s client list
	cdc.k8sClients, cdc.Ye = yceutils.CreateK8sClientList(cdc.apiServers)
	if cdc.CheckError() {
		return
	}

	// Publish deployment to every datacenter
	orgName := cd.OrgName
	cdc.createDeployment(orgName, &cd.Deployment)
	if cdc.CheckError() {
		return
	}

	// Encode cd.DcIdList to json
	dcIdList, _ := yceutils.EncodeDcIdList(cd.DcIdList)

	// Encode k8s.deployment to json
	kd, _ := json.Marshal(cd.Deployment)
	oId, _ := strconv.Atoi(orgId)

	// Insert into mysql.Deployment
	cdc.createMysqlDeployment(true, cd.AppName, userId, string(kd), "", dcIdList, int32(oId))
	if cdc.CheckError() {
		return
	}

	cdc.WriteOk("")
	log.Infoln("CreateDeploymentController over!")
	return
}
