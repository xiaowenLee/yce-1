package deploy

import (
	"app/backend/common/util/Placeholder"
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	mydeployment "app/backend/model/mysql/deployment"
	myoption "app/backend/model/mysql/option"
	"encoding/json"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	"app/backend/model/yce/deploy"
	yce "app/backend/controller/yce"
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
	AppName  string `json: "appName"`
	DcIdList []int32 `json: "dcIdList"`
	UserId   string `json: "userId"`
	Image    string `json: "image"`
	Revision string `json: "revision"`
	Comments string `json: "comments"`
}

// Create k8sClient for every ApiServer
func (rdc *RollbackDeploymentController) createK8sClients() {

	server := rdc.apiServer
	config := &restclient.Config{
		Host: server,
	}

	c, err := client.New(config)
	if err != nil {
		log.Errorf("RollbackDeployment createK8sClient Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	rdc.k8sClient = c
	log.Infof("RollbackDeployment createK8sClients: client=%p, apiServer=%s", c, server)

	return
}

// Get Deployment by deployment-name
func (rdc *RollbackDeploymentController) getDeploymentByName() {

	// Get namespace(org.Name) by orgId
	org, err := organization.GetOrganizationById(rdc.orgId)
	if err != nil {
		log.Errorf("RollbackDeployment getDatacentersByOrgId Error: orgId=%s, error=%s", rdc.orgId, err)
		rdc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	namespace := org.Name
	dp, err := rdc.k8sClient.Extensions().Deployments(namespace).Get(rdc.name)
	if err != nil {
		log.Errorf("RollbackDeployment getDeployByName Error: apiServer=%s, namespace=%s, deployment-name=%s, err=%s\n",
			rdc.apiServer, namespace, rdc.name, err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_GET_DEPLOYMENT, "")
		return
	}

	rdc.deployment = dp

	log.Infof("RollbackDeployment GetDeploymentByName over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		rdc.apiServer, namespace, rdc.name, dp)
}

// Get ApiServer by DcId
func (rdc *RollbackDeploymentController) getApiServerAndK8sClientByDcId() {

	// ApiServer
	dc := new(mydatacenter.DataCenter)

	//dcId, _ := strconv.Atoi(rdc.r.DcIdList.DcIdList[0])
	dcId := rdc.r.DcIdList[0]
	err := dc.QueryDataCenterById(dcId)
	//err := dc.QueryDataCenterById(int32(dcId))
	if err != nil {
		log.Errorf("RollbackDeployment getApiServerById QueryDataCenterById Error: dcId=%d, err=%s", dcId, err)
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
		log.Errorf("createK8sClient Error: err=%s", err)
		rdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	rdc.k8sClient = c
	log.Infof("RollbackDeployment GetApiServerAndK8sClientByDcId over: apiServer=%s, k8sClient=%p",
		rdc.apiServer, rdc.k8sClient)
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

// Encode DcIdList
func (rdc *RollbackDeploymentController) encodeDcIdList() string{
	dcIdList := &deploy.DcIdListType{
		DcIdList:rdc.r.DcIdList,
	}
	data, _ := json.Marshal(dcIdList)

	log.Infof("RollbackDeploymentController encodeDcIdList: dcIdList=%s", string(data))
	return string(data)
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

	// Get ApiServer and K8sClient
	rdc.name = rdc.r.AppName
	rdc.getApiServerAndK8sClientByDcId()
	if rdc.CheckError() {
		return
	}

	// Get Deployment by name
	rdc.getDeploymentByName()
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
	dcIdList := rdc.encodeDcIdList()

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
