package deploy

import (
	"app/backend/common/util/Placeholder"
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mydeployment "app/backend/model/mysql/deployment"
	"encoding/json"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
)

type RollbackDeploymentController struct {
	yce.Controller
	k8sClient  *client.Client
	apiServer  string
	orgId      string
	name       string
	r          *RollbackDeployParam
	deployment *extensions.Deployment
}

type RollbackDeployParam struct {
	AppName  string  `json: "appName"`
	DcIdList []int32 `json: "dcIdList"`
	UserId   string  `json: "userId"`
	Image    string  `json: "image"`
	Revision string  `json: "revision"`
	Comments string  `json: "comments"`
}

// Rollback
func (rdc *RollbackDeploymentController) rollback() {

	dr := new(extensions.DeploymentRollback)
	dr.Name = rdc.deployment.Name
	dr.UpdatedAnnotations = make(map[string]string, 0)
	dr.UpdatedAnnotations[ROLLBACK_IMAGE] = rdc.r.Image
	dr.UpdatedAnnotations[ROLLBACK_USERID] = rdc.r.UserId
	dr.UpdatedAnnotations[ROLLBACK_REVISION_ANNOTATION] = rdc.r.Comments
	dr.UpdatedAnnotations[ROLLBACK_CHANGE_CAUSE] = rdc.r.Comments

	// Convert revision from string to int64
	revision, _ := strconv.ParseInt(rdc.r.Revision, 10, 64)
	dr.RollbackTo = extensions.RollbackConfig{Revision: revision}

	// Rollback
	//err := rdc.k8sClient.Extensions().Deployments(rdc.deployment.Name).Rollback(dr)
	err := rdc.k8sClient.Extensions().Deployments(rdc.deployment.Namespace).Rollback(dr)
	if err != nil {
		log.Errorf("Deployment Rollback Error: err=%s\n", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_ROLLBACK_DEPLOYMENT, "")
		return
	}

	log.Infof("RollbackDeployment over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		rdc.apiServer, rdc.deployment.Namespace, rdc.deployment.Name, rdc.deployment)
}

// Create Deployment(mysql) and insert it into db
func (rdc *RollbackDeploymentController) createMysqlDeployment(success bool, name, json, reason, dcList string, userId, orgId int32) error {

	uph := placeholder.NewPlaceHolder(ROLLBACK_ACTION_URL)
	orgIdString := strconv.Itoa(int(orgId))
	actionUrl := uph.Replace("<orgId>", orgIdString, "<name>", name)
	actionOp := userId
	dp := mydeployment.NewDeployment(name, ROLLBACK_ACTION_VERBE, actionUrl, dcList, reason, json, "Rolilng Update a Deployment", int32(ROLLBACK_ACTION_TYPE), actionOp, int32(1), orgId)
	err := dp.InsertDeployment()
	if err != nil {
		log.Errorf("CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		rdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	log.Infof("RollbackDeployment CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

// POST /api/v1/organizations/{orgId}/deployments/{name}/rollback
func (rdc RollbackDeploymentController) Post() {
	rdc.orgId = rdc.Param("orgId")
	rdc.name = rdc.Param("name")
	sessionIdFromClient := rdc.RequestHeader("Authorization")

	log.Debugf("RollbackDeploymentController Params: sessionId=%s, orgId=%s, name=%s", sessionIdFromClient, rdc.orgId, rdc.name)

	// Validate the session
	rdc.ValidateSession(sessionIdFromClient, rdc.orgId)
	if rdc.CheckError() {
		return
	}

	// RollbackDeployment Param
	rdc.r = new(RollbackDeployParam)
	err := rdc.ReadJSON(rdc.r)
	if err != nil {
		log.Errorf("RollbackDeploymentController ReadJSON Error: error=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}

	if rdc.CheckError() {
		return
	}

	rdc.name = rdc.r.AppName

	if len(rdc.r.DcIdList) == 0 {
		log.Errorln("Empty DcIdList!")
		rdc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}

	if rdc.CheckError() {
		return
	}

	dcId := rdc.r.DcIdList[0]

	// Get ApiServer
	rdc.apiServer, rdc.Ye = yceutils.GetApiServerByDcId(dcId)
	if rdc.CheckError() {
		return
	}

	// Create K8sClient
	rdc.k8sClient, rdc.Ye = yceutils.CreateK8sClient(rdc.apiServer)
	if rdc.CheckError() {
		return
	}

	// Get Namespace
	namespace, ye := yceutils.GetOrgNameByOrgId(rdc.orgId)
	if ye != nil {
		rdc.Ye = ye
	}
	if rdc.CheckError() {
		return
	}

	// Get Deployment by name
	rdc.deployment, rdc.Ye = yceutils.GetDeploymentByNameAndNamespace(rdc.k8sClient, rdc.name, namespace)
	if rdc.CheckError() {
		return
	}

	// RollBack
	rdc.rollback()
	if rdc.CheckError() {
		return
	}

	// Encode deployment to string
	dd, _ := json.Marshal(rdc.deployment)

	// Encode DcIdList to string
	dcIdList, ye := yceutils.EncodeDcIdList(rdc.r.DcIdList)
	if ye != nil {
		rdc.Ye = ye
	}
	if rdc.CheckError() {
		return
	}

	// Convert UserId from string to int32
	userId, _ := strconv.Atoi(rdc.r.UserId)
	oId, _ := strconv.Atoi(rdc.orgId)

	// Insert into MySQL.Deployment
	rdc.createMysqlDeployment(true, rdc.r.AppName, string(dd), rdc.r.Comments, dcIdList, int32(userId), int32(oId))
	if rdc.CheckError() {
		return
	}

	rdc.WriteOk("")
	log.Infoln("Rollback DeploymentController over!")
	return
}
