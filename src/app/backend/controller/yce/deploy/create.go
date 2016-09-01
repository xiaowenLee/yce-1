package deploy

import (
	mylog "app/backend/common/util/log"
	"app/backend/common/util/Placeholder"
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
	"strconv"
	"strings"
)

const (
	ACTION_TYPE  = myoption.ONLINE
	ACTION_VERBE = "POST"
	ACTION_URL   = "/api/v1/organization/<orgId>/users/<userId>/deployments"
)

type CreateDeployController struct {
	*iris.Context
	k8sClients []*client.Client
	apiServers []string
	Ye         *myerror.YceError
}

func (cdc *CreateDeployController) WriteBack() {
	cdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateDeployController Response YceError: controller=%p, code=%d, note=%s", cdc, cdc.Ye.Code, myerror.Errors[cdc.Ye.Code].LogMsg)
	cdc.Write(cdc.Ye.String())
}

// Validate Session
func (cdc *CreateDeployController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		cdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		cdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	return
}

// Get ApiServer by dcId
func (cdc *CreateDeployController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		cdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateDeployController getApiServerByDcId: apiServer=%s, dcId=%d", apiServer, dcId)
	return apiServer
}

// Get ApiServer List for dcIdList
func (cdc *CreateDeployController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := cdc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("CreateDeployController getApiServerList Error")
			return
		}

		cdc.apiServers = append(cdc.apiServers, apiServer)
	}
	return
}

// Create k8sClients for every ApiServer
func (cdc *CreateDeployController) createK8sClients() {

	// Foreach every ApiServer to create it's k8sClient
	cdc.k8sClients = make([]*client.Client, 0)

	for _, server := range cdc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			mylog.Log.Errorf("createK8sClient Error: err=%s", err)
			cdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		cdc.k8sClients = append(cdc.k8sClients, c)
		cdc.apiServers = append(cdc.apiServers, server)
		mylog.Log.Infof("Append a new client to cdc.k8sClients array: c=%p, apiServer=%s", c, server)
	}

	return
}

// Publish k8s.Deployment to every datacenter which in dcIdList
func (cdc *CreateDeployController) createDeployment(namespace string, deployment *extensions.Deployment) {

	// Foreach every k8sClient to create deployment
	for index, cli := range cdc.k8sClients {
		_, err := cli.Extensions().Deployments(namespace).Create(deployment)
		if err != nil {
			mylog.Log.Errorf("createDeployment Error: apiServer=%s, namespace=%s, err=%s",
				cdc.apiServers[index], namespace, err)
			cdc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_DEPLOYMENT, "")
			return
		}

		mylog.Log.Infof("Create deployment successfully: namespace=%s, apiserver=%s", namespace, cdc.apiServers[index])
	}
}

// Create Deployment(mysql) and insert it into db
func (cdc *CreateDeployController) createMysqlDeployment(success bool, name, orgId, userId, json, reason, dcList string) error {

	uph := placeholder.NewPlaceHolder(ACTION_URL)
	actionUrl := uph.Replace("<orgId>", orgId, "<userId>", userId)
	actionOp, _ := strconv.Atoi(userId)
	dp := mydeployment.NewDeployment(name, ACTION_VERBE, actionUrl, dcList, reason, json, "Create a Deployment", ACTION_TYPE, int32(actionOp), int32(1))
	err := dp.InsertDeployment()
	if err != nil {
		mylog.Log.Errorf("CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		cdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	mylog.Log.Infof("CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

// POST /api/v1/organizations/{orgId}/users/{userId}/deployments
func (cdc CreateDeployController) Post() {
	sessionIdFromClient := cdc.RequestHeader("Authorization")
	orgId := cdc.Param("orgId")
	userId := cdc.Param("userId")

	// Validate OrgId error
	cdc.validateSession(sessionIdFromClient, orgId)

	if cdc.Ye != nil {
		cdc.WriteBack()
		return
	}

	// Parse data: deploy.CreateDeployment
	cd := new(deploy.CreateDeployment)
	cdc.ReadJSON(cd)

	// Get DcIdList
	cdc.getApiServerList(cd.DcIdList)
	if cdc.Ye != nil {
		cdc.WriteBack()
		return
	}

	// Create k8s clients
	cdc.createK8sClients()
	if cdc.Ye != nil {
		cdc.WriteBack()
		return
	}

	// Publish deployment to every datacenter
	orgName := cd.OrgName
	cdc.createDeployment(orgName, &cd.Deployment)
	if cdc.Ye != nil {
		cdc.WriteBack()
		return
	}

	// Encode cd.DcIdList to json
	dcl, _ := json.Marshal(cd.DcIdList)

	// Encode k8s.deployment to json
	kd, _ := json.Marshal(cd.Deployment)

	// Insert into mysql.Deployment
	cdc.createMysqlDeployment(true, cd.AppName, orgId, userId, string(kd), "", string(dcl))
	if cdc.Ye != nil {
		cdc.WriteBack()
		return
	}

	// ToDo: 数据库中两个dcList的格式不一致,要改过来,统一叫DcIdList
	// ToDo: 发布出错时也要插入数据库

	cdc.Ye = myerror.NewYceError(myerror.EOK, "")
	cdc.WriteBack()
	// TODO: 成功写回
	mylog.Log.Infoln("CreateDeploymentController over!")
	return
}
