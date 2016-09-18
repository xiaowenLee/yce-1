package deploy

import (
	mylog "app/backend/common/util/log"
	"app/backend/common/util/Placeholder"
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
	yce "app/backend/controller/yce"
)


type CreateDeploymentController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
}

// Get ApiServer by dcId
func (cdc *CreateDeploymentController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		cdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("CreateDeploymentController getApiServerByDcId: apiServer=%s, dcId=%d", apiServer, dcId)
	return apiServer
}

// Get ApiServer List for dcIdList
func (cdc *CreateDeploymentController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := cdc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			log.Errorf("CreateDeploymentController getApiServerList Error")
			cdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
			return
		}

		cdc.apiServers = append(cdc.apiServers, apiServer)
	}

	log.Infof("CreateDeploymentController getApiServerList success: len(apiSeverList)=%d", len(cdc.apiServers))
	return
}

// Create k8sClients for every ApiServer
func (cdc *CreateDeploymentController) createK8sClients() {

	// Foreach every ApiServer to create it's k8sClient
	cdc.k8sClients = make([]*client.Client, 0)

	for _, server := range cdc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			log.Errorf("createK8sClient Error: err=%s", err)
			cdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		cdc.k8sClients = append(cdc.k8sClients, c)
		cdc.apiServers = append(cdc.apiServers, server)
		log.Infof("Append a new client to cdc.k8sClients array: c=%p, apiServer=%s", c, server)
	}
	log.Infof("CreateDeploymentController createK8sClients success: len(k8sCLients)=%d", len(cdc.k8sClients))
	return
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

// Encode JSON of dcIdList
func (cdc *CreateDeploymentController) encodeDcIdList(dcIdList []int32) string {
	dcIds := &deploy.DcIdListType{
		DcIdList:dcIdList,
	}

	data, _ := json.Marshal(dcIds)

	log.Infof("CreateDeploymentController encodeDcIdList: dcIdList=%s", string(data))
	return string(data)
}

// POST /api/v1/organizations/{orgId}/users/{userId}/deployments
func (cdc CreateDeploymentController) Post() {
	sessionIdFromClient := cdc.RequestHeader("Authorization")
	orgId := cdc.Param("orgId")
	userId := cdc.Param("userId")

	log.Debugf("CreateDeploymentController get Params:  sessionIdFromClient=%s, orgId=%s, userId=%s", sessionIdFromClient, orgId, userId)

	// Validate OrgId error
	cdc.ValidateSession(sessionIdFromClient, orgId)
	if cdc.Ye != nil {
		cdc.WriteBack()
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
	cdc.getApiServerList(cd.DcIdList)
	if cdc.CheckError() {
		return
	}


	// Create k8s clients
	cdc.createK8sClients()
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
	dcIdList := cdc.encodeDcIdList(cd.DcIdList)

	// Encode k8s.deployment to json
	kd, _ := json.Marshal(cd.Deployment)

	oId, _ := strconv.Atoi(orgId)

	// Insert into mysql.Deployment
	cdc.createMysqlDeployment(true, cd.AppName, userId, string(kd), "", dcIdList, int32(oId))
	if cdc.CheckError() {
		return
	}

	// ToDo: 数据库中两个dcList的格式不一致,要改过来,统一叫DcIdList
	// ToDo: 发布出错时也要插入数据库
	cdc.WriteOk("")
	log.Infoln("CreateDeploymentController over!")
	return
}
