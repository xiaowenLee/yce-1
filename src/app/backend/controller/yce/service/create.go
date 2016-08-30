package service

import (
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	mynodeport "app/backend/model/mysql/nodeport"
	"github.com/kataras/iris"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"app/backend/common/util/session"
	"app/backend/model/yce/service"
	"strconv"
	"strings"
	"github.com/kubernetes/kubernetes/pkg/client/restclient"
)

type CreateServiceController struct {
	*iris.Context
	Ye *myerror.YceError
	k8sClients []client.Client
	apiServers []string
}

func (csc *CreateServiceController) WriteBack() {
	csc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateServiceController Response YceError: controller=%p, code=%d, note=%s", csc, csc.Ye.Code, myerror.Errors[csc.Ye.Code].LogMsg)
	csc.Write(csc.Ye.String())
}

// Validate Session
func (csc *CreateServiceController) validateSession(sessionIdFromClient, orgId string) {
	ss := session.SessionStoreInstance()

	// validate the session
	ok, err := ss.ValidateOrgId(sessionIdFromClient, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session Error: sessionId=%s, error=%s", sessionIdFromClient, err)
		csc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Invalidate SessionId
	if !ok {
		mylog.Log.Errorf("Validate Session Failed: sessionId=%s, error=%s", sessionIdFromClient, err)
		csc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
	}

	return
}

// Get ApiServer by dcId
func (csc *CreateServiceController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		csc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

// Get ApiServer List for dcIdList
func (csc *CreateServiceController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := csc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("CreateServiceController getApiServerList Error")
			return
		}

		csc.apiServers = append(csc.apiServers, apiServer)
	}

	return
}

// Create k8sClient for every ApiServer
func (csc *CreateServiceController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	csc.k8sClients := make([]*client.Client, len(csc.apiServers))


	for _, server := range csc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			mylog.Log.Errorf("CreateK8sClient Error: error=%s", err)
			csc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		csc.k8sClients = append(csc.k8sClients, c)
		// why??
		//csc.apiServers = append(csc.apiServers, server)
		mylog.Log.Infof("Append a new client to csc.K8sClients array: c=%p, apiServer=%s", c, server})
	}

	return
}

// why need return value?
// Publish k8s.Service to every datacenter which in dcIdList
func (csc *CreateServiceController) createService(namespace string, service *service.Service) {
	// Foreach every K8sClient to create service
	for index, cli := range csc.k8sClients {
		_, err := cli.Services(namespace).Create(service)
		if err != nil {
			mylog.Log.Errorf("createService Error: apiServer=%s, namespace=%s, error=%s", csc.apiServers[index], namespace, err)
			csc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_SERVICE, "")
			return
		}

		mylog.Log.Infof("Create Service successfully: namespace=%s, apiServer=%s", namespace, csc.apiServers[index])
	}

	return
}

// create NodePort(mysql) and insert it into db
func (csc *CreateServiceController)createMysqlNodePort(success bool, nodePort int32, dcIdList []int32, svcName string, op string) {
	for _, dcId := range dcIdList {
		np := mynodeport.NewNodePort(nodePort, dcId, svcName, op)
		err := np.InsertNodePort(op)
		if err != nil {
			mylog.Log.Errorf("createMysqlNodePort Error: nodeport=%d, dcId=%d, svcName=%s, error=%s", np.Port, np.DcId, np.SvcName, err)
			csc.Ye = myerror.NewYceError(myerror.EYCE_NODEPORT_EXIST, "")
			return
		}

		mylog.Log.Infof("createMysqlNodePort Successfully: nodeport=%d, dcId=%d, svcName=%s", np.Port, np.DcId, np.SvcName)
	}

	return
}


func (csc *CreateServiceController) Post() {
	sessionIdFromClient := csc.RequestHeader("Authorization")
	orgId := csc.Params("orgId")
	userId := csc.Params("userId")

	// Validate OrgId error
	csc.validateSession(sessionIdFromClient, orgId)
	if csc.Ye != nil {
		csc.WriteBack()
		return
	}

	// Parse data: service.
	sd := new(service.CreateService)
	csc.ReadJSON(sd)

	// Get DcIdList
	csc.getApiServerList(sd.DcIdList)
	if csc.Ye != nil {
		csc.WriteBack()
		return
	}

	// Get K8sClient
	csc.createK8sClients()
	if csc.Ye != nil {
		csc.WriteBack()
		return
	}

	// Publish server to every datacenter
	orgName := sd.OrgName
	csc.createService(orgName, &sd.Service)
	if csc.Ye != nil {
		csc.WriteBack()
		return
	}

	// And NodePort to MySQL nodeport table
	// assuming single port service
	// or spec.ports[0].nodeport == 0 for judging
	hasNodePort := strings.EqualFold(sd.Service.Spec.Type, "NodePort")
	op, _ := strconv.Atoi(userId)
	csc.createMysqlNodePort(hasNodePort, sd.Service.Spec.Ports[0].NodePort, sd.DcIdList, sd.Service.ObjectMeta.Name, op)
	if csc.Ye != nil {
		csc.WriteBack()
		return
	}

	csc.Ye = myerror.NewYceError(myerror.EOK, "")
	mylog.Log.Infoln("CreateServiceController over!")
	return

