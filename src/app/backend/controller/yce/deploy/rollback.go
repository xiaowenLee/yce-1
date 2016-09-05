package deploy

import (
	"app/backend/common/util/Placeholder"
	mylog "app/backend/common/util/log"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	mydeployment "app/backend/model/mysql/deployment"
	myoption "app/backend/model/mysql/option"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	"app/backend/model/yce/deploy"
)

const (
	ROLLBACK_ACTION_TYPE                = myoption.ROLLINGBACK
	ROLLBACK_ACTION_VERBE               = "POST"
	ROLLBACK_ACTION_URL                 = "/api/v1/organization/<orgId>/deployments/<name>/rollback"
	ROLLBACK_REVISION_ANNOTATION string = "deployment.kubernetes.io/revision"
	ROLLBACK_IMAGE                      = "image"
	ROLLBACK_USERID                     = "userId"
	ROLLBACK_CHANGE_CAUSE string = "kubernetes.io/change-cause"
)

type RollbackDeployController struct {
	*iris.Context
	k8sClient  *client.Client
	apiServer  string
	Ye         *myerror.YceError
	orgId      string
	name       string
	r          *RollbackDeployParam
	deployment *extensions.Deployment
}

type RollbackDeployParam struct {
	AppName  string `json: "appName"`
	DcIdList []int32 `json: "dcIdList"`
	UserId   string `json: "userId"`
	Image    string `json: "image"`
	Revision string `json: "revision"`
	Comments string `json: "comments"`
}

func (rdc *RollbackDeployController) WriteBack() {
	rdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("RollbackDeployController Response YceError: controller=%p, code=%d, note=%s", rdc, rdc.Ye.Code, myerror.Errors[rdc.Ye.Code].LogMsg)
	rdc.Write(rdc.Ye.String())
}

// Validate Session
func (rdc *RollbackDeployController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		rdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		rdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	mylog.Log.Infof("RollbackDeployment sessionId successfully: sessionId=%s, orgId=%d", sessionId, orgId)
	return
}


// Create k8sClient for every ApiServer
func (rdc *RollbackDeployController) createK8sClients() {

	server := rdc.apiServer
	config := &restclient.Config{
		Host: server,
	}

	c, err := client.New(config)
	if err != nil {
		mylog.Log.Errorf("RollbackDeployment createK8sClient Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	rdc.k8sClient = c
	mylog.Log.Infof("RollbackDeployment createK8sClients: client=%p, apiServer=%s", c, server)

	return
}

// Get Deployment by deployment-name
func (rdc *RollbackDeployController) getDeploymentByName() {

	// Get namespace(org.Name) by orgId
	org, err := organization.GetOrganizationById(rdc.orgId)
	if err != nil {
		mylog.Log.Errorf("RollbackDeployment getDatacentersByOrgId Error: orgId=%s, error=%s", rdc.orgId, err)
		rdc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	namespace := org.Name
	dp, err := rdc.k8sClient.Extensions().Deployments(namespace).Get(rdc.name)
	if err != nil {
		mylog.Log.Errorf("RollbackDeployment getDeployByName Error: apiServer=%s, namespace=%s, deployment-name=%s, err=%s\n",
			rdc.apiServer, namespace, rdc.name, err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_GET_DEPLOYMENT, "")
		return
	}

	rdc.deployment = dp

	mylog.Log.Infof("RollbackDeployment GetDeploymentByName over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		rdc.apiServer, namespace, rdc.name, dp)
}

// Get ApiServer by DcId
func (rdc *RollbackDeployController) getApiServerAndK8sClientByDcId() {

	// ApiServer
	dc := new(mydatacenter.DataCenter)

	//dcId, _ := strconv.Atoi(rdc.r.DcIdList.DcIdList[0])
	dcId := rdc.r.DcIdList[0]
	err := dc.QueryDataCenterById(dcId)
	//err := dc.QueryDataCenterById(int32(dcId))
	if err != nil {
		mylog.Log.Errorf("RollbackDeployment getApiServerById QueryDataCenterById Error: dcId=%d, err=%s", dcId, err)
		rdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	rdc.apiServer = host + ":" + port

	// K8sClient
	config := &restclient.Config{
		Host: rdc.apiServer,
	}

	c, err := client.New(config)
	if err != nil {
		mylog.Log.Errorf("createK8sClient Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	rdc.k8sClient = c
	mylog.Log.Infof("RollbackDeployment GetApiServerAndK8sClientByDcId over: apiServer=%s, k8sClient=%p",
		rdc.apiServer, rdc.k8sClient)
}

// Rollback
func (rdc *RollbackDeployController) rollback() {

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
		mylog.Log.Errorf("Deployment Rollback Error: err=%s\n", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_ROLLBACK_DEPLOYMENT, "")
		return
	}

	mylog.Log.Infof("RollbackDeployment over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		rdc.apiServer, rdc.deployment.Namespace, rdc.deployment.Name, rdc.deployment)
}

// Create Deployment(mysql) and insert it into db
func (rdc *RollbackDeployController) createMysqlDeployment(success bool, name, json, reason, dcList string, userId, orgId int32) error {

	uph := placeholder.NewPlaceHolder(ROLLBACK_ACTION_URL)
	orgIdString := strconv.Itoa(int(orgId))
	actionUrl := uph.Replace("<orgId>", orgIdString, "<name>", name)
	actionOp := userId
	dp := mydeployment.NewDeployment(name, ROLLBACK_ACTION_VERBE, actionUrl, dcList, reason, json, "Rolilng Update a Deployment", int32(ROLLBACK_ACTION_TYPE), actionOp, int32(1), orgId)
	err := dp.InsertDeployment()
	if err != nil {
		mylog.Log.Errorf("CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		rdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	mylog.Log.Infof("RollbackDeployment CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

// Encode DcIdList
func (rdc *RollbackDeployController) encodeDcIdList() string{
	dcIdList := &deploy.DcIdListType{
		DcIdList:rdc.r.DcIdList,
	}
	data, _ := json.Marshal(dcIdList)

	mylog.Log.Infof("RollbackDeployController encodeDcIdList: dcIdList=%s", string(data))
	return string(data)
}


// POST /api/v1/organizations/{orgId}/deployments/{name}/rollback
func (rdc RollbackDeployController) Post() {
	rdc.orgId = rdc.Param("orgId")
	rdc.name = rdc.Param("name")


	// Validate the session
	sessionIdFromClient := rdc.RequestHeader("Authorization")
	rdc.validateSession(sessionIdFromClient, rdc.orgId)
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	// RollbackDeployment Param
	rdc.r = new(RollbackDeployParam)
	rdc.ReadJSON(rdc.r)


	// Get ApiServer and K8sClient
	rdc.name = rdc.r.AppName
	rdc.getApiServerAndK8sClientByDcId()
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	// Get Deployment by name
	rdc.getDeploymentByName()
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	// RollBack
	rdc.rollback()
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	// Encode deployment to string
	dd, _ := json.Marshal(rdc.deployment)

	// Encode DcIdList to string
	dcIdList := rdc.encodeDcIdList()

	// Convert UserId from string to int32
	userId, _ := strconv.Atoi(rdc.r.UserId)
	oId, _ := strconv.Atoi(rdc.orgId)

	// Insert into MySQL.Deployment
	rdc.createMysqlDeployment(true, rdc.r.AppName, string(dd), rdc.r.Comments, dcIdList, int32(userId), int32(oId))
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	rdc.Ye = myerror.NewYceError(myerror.EOK, "")
	rdc.WriteBack()
	mylog.Log.Infoln("Rollback DeploymentController over!")
	return
}
