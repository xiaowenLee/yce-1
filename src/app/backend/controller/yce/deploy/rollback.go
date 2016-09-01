package deploy

import (
	"app/backend/common/util/Placeholder"
	mylog "app/backend/common/util/log"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
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
	ROLLBACK_ACTION_TYPE  = myoption.ROLLINGBACK
	ROLLBACK_ACTION_VERBE = "POST"
	ROLLBACK_ACTION_URL   = "/api/v1/organization/<orgId>/deployments/<name>/rollback"
	ROLLBACK_REVISION_ANNOTATION string = "deployment.kubernetes.io/revision"
	ROLLBACK_IMAGE = "image"
	ROLLBACK_USERID = "userId"
)

type RollbackDeployController struct {
	*iris.Context
	k8sClient *client.Client
	apiServer string
	Ye         *myerror.YceError
	orgId string
	name string
	r *RollbackDeployParam
	deployment *extensions.Deployment
}

type RollbackDeployParam struct {
	AppName string `json: "appName"`
	DcId string `json: "dcId"`
	UserId string `json: "userId"`
	Image string `json: "image"`
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

	return
}


// Get ApiServer by dcId
func (rdc *RollbackDeployController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
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
func (rdc *RollbackDeployController) getApiServer(dcId int32) {
	// Get ApiServer
	apiServer := rdc.getApiServerByDcId(dcId)
	if strings.EqualFold(apiServer, "") {
		mylog.Log.Errorf("RollbackDeployController getApiServerList Error")
		return
	}

	//rdc.apiServer = append(rdc.apiServer, apiServer)
	rdc.apiServer = apiServer
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
		mylog.Log.Errorf("createK8sClient Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	rdc.k8sClient = c
	mylog.Log.Infof("Append a new client to rdc.k8sClient array: c=%p, apiServer=%s", c, server)

	return
}

// Get Deployment by deployment-name
func (rdc *RollbackDeployController) getDeploymentByName() {

	// Get namespace(org.Name) by orgId
	org, err := organization.GetOrganizationById(rdc.orgId)
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", rdc.orgId, err)
		rdc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	namespace := org.Name
	dp, err := rdc.k8sClient.Extensions().Deployments(namespace).Get(rdc.name)
	if err != nil {
		mylog.Log.Errorf("getDeployByName Error: apiServer=%s, namespace=%s, deployment-name=%s, err=%s\n",
			rdc.apiServer, namespace, rdc.name, err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_GET_DEPLOYMENT, "")
		return
	}

	rdc.deployment = dp

	mylog.Log.Infof("GetDeploymentByName over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		rdc.apiServer, namespace, rdc.name, dp)
}

// Get ApiServer by DcId
func (rdc *RollbackDeployController) getApiServerAndK8sClientByDcId() {

	// ApiServer
	dc := new(mydatacenter.DataCenter)
	dcId, _ := strconv.Atoi(rdc.dcId)
	err := dc.QueryDataCenterById(int32(dcId))
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: dcId=%d, err=%s", dcId, err)
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
	mylog.Log.Infof("GetApiServerAndK8sClientByDcId over: apiServer=%s, k8sClient=%p",
		rdc.apiServer, rdc.k8sClient)
}


// Rollback
func (rdc *RollbackDeployController) rollback() {

	dr := new(extensions.DeploymentRollback)
	dr.Name = rdc.deployment.Name
	dr.UpdatedAnnotations[ROLLBACK_IMAGE] = rdc.r.Image
	dr.UpdatedAnnotations[ROLLBACK_USERID] = rdc.r.UserId
	dr.UpdatedAnnotations[ROLLBACK_REVISION_ANNOTATION] = rdc.r.Comments

	// Convert revision from string to int64
	revision, _ := strconv.ParseInt(rdc.r.Revision, 10, 64)
	dr.RollbackTo = extensions.RollbackConfig{Revision: revision}

	// Rollback
	err := rdc.k8sClient.Extensions().Deployments(rdc.deployment.Name).Rollback(dr)
	if err != nil {
		mylog.Log.Errorf("Deployment Rollback Error: err=%s\n", err)
		hdc.Ye = myerror.NewYceError(myerror.EKUBE_ROLLBACK_DEPLOYMENT, "")
		return
	}

	mylog.Log.Infof("Deployment Rollback over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		rdc.apiServer, rdc.deployment.Namespace, rdc.deployment.name, rdc.deployment)
}

// Create Deployment(mysql) and insert it into db
func (rdc *RollbackDeployController) createMysqlDeployment(success bool, name, orgId, json, reason, dcList string, userId int32) error {

	uph := placeholder.NewPlaceHolder(ROLLBACK_ACTION_URL)
	actionUrl := uph.Replace("<orgId>", orgId, "<name>", name)
	//actionOp, _ := strconv.Atoi(userId)
	actionOp := userId
	dp := mydeployment.NewDeployment(name, ROLLBACK_ACTION_VERBE, actionUrl, dcList, reason, json, "Rolilng Update a Deployment", int32(ROLLBACK_ACTION_TYPE), actionOp, int32(1))
	err := dp.InsertDeployment()
	if err != nil {
		mylog.Log.Errorf("CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		rdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	mylog.Log.Infof("CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

// 在查看历史中,要显示在那个数据中心做的回滚,为了保证接口一致,存储成DcIdList的Json格式
// Encode DcIdList
func (rdc *RollbackDeployController) encodeDcIdList() string {
	dcIdList := new(deploy.DcIdListType)
	dcId, _  := strconv.Atoi(rdc.r.DcId)
	*dcIdList = append(*dcIdList, int32(dcId))

	data, _ := json.Marshal(dcIdList)
	return string(data)
}


// POST /api/v1/organizations/{orgId}/deployments/{name}/rollback
func (rdc *RollbackDeployController) Post() {
	rdc.orgId = rdc.Param("orgId")
	rdc.name = rdc.Param("name")

	// RollbackDeployment Param
	rdc.r = new(RollbackDeployParam)
	rdc.ReadJSON(rdc.r)

	// Validate the session
	sessionIdFromClient := rdc.RequestHeader("Authorization")
	rdc.validateSession(sessionIdFromClient, orgId)
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	// Get ApiServer and K8sClient
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

	// Encode DcIdList
	dcIdList := rdc.encodeDcIdList()

	// Insert into MySQL.Deployment
	rdc.createMysqlDeployment(true, rdc.r.AppName, rdc.orgId, string(dd), rdc.r.Comments, dcIdList, rdc.r.UserId)
	if rdc.Ye != nil {
		rdc.WriteBack()
		return
	}

	rdc.Ye = myerror.NewYceError(myerror.EOK, "")
	rdc.WriteBack()
	mylog.Log.Infoln("Rolling DeploymentController over!")
	return
}