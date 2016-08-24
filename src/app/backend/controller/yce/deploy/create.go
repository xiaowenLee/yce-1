package deploy

import (
	"app/backend/common/util/placeholder"
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
	mylog "app/backend/common/util/log"
	"strconv"
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
}

// Validate Session
func (cdc *CreateDeployController) validateSession(sessionId, orgId string) (*myerror.YceError, error) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		cdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		cdc.Write(errJson)
		return ye, err
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		cdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		cdc.Write(errJson)
		return ye, err
	}

	return nil, nil
}

// Get ApiServer by dcId
func (cdc *CreateDeployController) getApiServerByDcId(dcId int32) (string, error) {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		return "", err
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateDeployController getApiServerByDcId: apiServer=%s, dcId=%d", apiServer, dcId)
	return apiServer, nil
}

// Get ApiServer List for dcIdList
func (cdc *CreateDeployController) getApiServerList(dcIdList []int32) error {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer, err := cdc.getApiServerByDcId(dcId)
		if err != nil {
			mylog.Log.Errorf("getApiServerList error: err=%s", err)
			return err
		}

		cdc.apiServers = append(cdc.apiServers, apiServer)
	}
	return nil
}

// Create k8sClients for every ApiServer
func (cdc *CreateDeployController) createK8sClients() error {

	// Foreach every ApiServer to create it's k8sClient
	cdc.k8sClients = make([]*client.Client, 0)

	for _, server := range cdc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)

		if err != nil {
			mylog.Log.Errorf("createK8sClient Error: err=%s", err)
			return err
		}

		cdc.k8sClients = append(cdc.k8sClients, c)
		cdc.apiServers = append(cdc.apiServers, server)
		mylog.Log.Infof("Append a new client to cdc.k8sClients array: c=%p, apiServer=%s", c, server)
	}

	return nil
}

// Publish k8s.Deployment to every datacenter which in dcIdList
func (cdc *CreateDeployController) createDeployment(namespace string, deployment *extensions.Deployment) error {

	// Foreach every k8sClient to create deployment
	for index, cli := range cdc.k8sClients {
		_, err := cli.Extensions().Deployments(namespace).Create(deployment)
		if err != nil {
			mylog.Log.Errorf("createDeployment Error: apiServer=%s, namespace=%s, err=%s",
				cdc.apiServers[index], namespace, err)
			return err
		}

		mylog.Log.Infof("Create deployment successfully: namespace=%s, apiserver=%s", namespace, cdc.apiServers[index])
	}

	return nil
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
	ye, err := cdc.validateSession(sessionIdFromClient, orgId)

	if ye != nil || err != nil {
		mylog.Log.Warnf("CreateDeployController validateSession: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		errJson, _ := ye.EncodeJson()
		cdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		cdc.Write(errJson)
		return
	}

	// Parse data: deploy.CreateDeployment
	cd := new(deploy.CreateDeployment)
	cdc.ReadJSON(cd)

	// Get DcIdList
	err = cdc.getApiServerList(cd.DcIdList)
	if err != nil {
		mylog.Log.Errorf("CreateDeployController getApiServerList: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1401, "ERR", "Get Api Server List Error")
		errJson, _ := ye.EncodeJson()
		cdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		cdc.Write(errJson)
		return
	}

	// Create k8s clients
	err = cdc.createK8sClients()
	if err != nil {
		mylog.Log.Errorf("CreateDeployController createK8sClients: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1402, "ERR", "create K8s Client Error")
		errJson, _ := ye.EncodeJson()
		cdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		cdc.Write(errJson)
		return
	}

	// Publish deployment to every datacenter
	orgName := cd.OrgName
	err = cdc.createDeployment(orgName, &cd.Deployment)
	if err != nil {
		mylog.Log.Errorf("CreateDeployController createDeployment: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1403, "ERR", "Publish K8s Deployment Error")
		errJson, _ := ye.EncodeJson()
		cdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		cdc.Write(errJson)
		return
	}

	// Encode cd.DcIdList to json
	dcl, _ := json.Marshal(cd.DcIdList)

	// Encode k8s.deployment to json
	kd, _ := json.Marshal(cd.Deployment)

	// Insert into mysql.Deployment
	err = cdc.createMysqlDeployment(true, cd.AppName, orgId, userId, string(kd), "", string(dcl))

	if err != nil {
		mylog.Log.Errorf("CreateDeployController createDeployment: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1404, "ERR", "Insert into MySQL Error")
		errJson, _ := ye.EncodeJson()
		cdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		cdc.Write(errJson)
		return
	}

	// ToDo: 数据库中两个dcList的格式不一致,要改过来,统一叫DcIdList
	// ToDo: 发布出错时也要插入数据库

	mylog.Log.Infoln("CreateDeploymentController over!")

}
