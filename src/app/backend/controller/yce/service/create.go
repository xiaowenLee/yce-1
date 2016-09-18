package service

import (
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	mynodeport "app/backend/model/mysql/nodeport"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"app/backend/model/yce/service"
	"strconv"
	"strings"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/api"
	yce "app/backend/controller/yce"
)

type CreateServiceController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
}

// Get ApiServer by dcId
func (csc *CreateServiceController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		csc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("CreateServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

// Get ApiServer List for dcIdList
func (csc *CreateServiceController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := csc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			log.Errorf("CreateServiceController getApiServerList Error")
			return
		}

		csc.apiServers = append(csc.apiServers, apiServer)
	}

	log.Infof("CreateServiceController getApiServerList: len(apiServer)=%d", len(csc.apiServers))
	return
}

// Create k8sClient for every ApiServer
func (csc *CreateServiceController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	//csc.k8sClients = make([]*client.Client, len(csc.apiServers))
	csc.k8sClients = make([]*client.Client, 0)


	for _, server := range csc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			log.Errorf("CreateK8sClient Error: error=%s", err)
			csc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		csc.k8sClients = append(csc.k8sClients, c)
		// why??
		//csc.apiServers = append(csc.apiServers, server)
		log.Infof("Append a new client to csc.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	log.Infof("CreateServiceController createK8sClients: len(k8sClients)=%d", len(csc.k8sClients))
	return
}

// why need return value?
// Publish k8s.Service to every datacenter which in dcIdList
func (csc *CreateServiceController) createService(namespace string, service *api.Service) {
	// Foreach every K8sClient to create service
	for index, cli := range csc.k8sClients {
		_, err := cli.Services(namespace).Create(service)
		if err != nil {
			log.Errorf("createService Error: apiServer=%s, namespace=%s, error=%s", csc.apiServers[index], namespace, err)
			csc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_SERVICE, "")
			return
		}

		log.Infof("Create Service successfully: namespace=%s, apiServer=%s", namespace, csc.apiServers[index])
	}

	log.Infof("CreateServiceController createService success")
	return
}

// create NodePort(mysql) and insert it into db
func (csc *CreateServiceController)createMysqlNodePort(success bool, nodePort int32, dcIdList []int32, svcName string, op int32) {
	for _, dcId := range dcIdList {
		np := mynodeport.NewNodePort(nodePort, dcId, svcName, op)
		err := np.InsertNodePort(op)
		if err != nil {
			log.Errorf("createMysqlNodePort Error: nodeport=%d, dcId=%d, svcName=%s, error=%s", np.Port, np.DcId, np.SvcName, err)
			csc.Ye = myerror.NewYceError(myerror.EYCE_NODEPORT_EXIST, "")
			return
		}

		log.Infof("createMysqlNodePort Successfully: nodeport=%d, dcId=%d, svcName=%s", np.Port, np.DcId, np.SvcName)
	}

	return
}


func (csc CreateServiceController) Post() {
	sessionIdFromClient := csc.RequestHeader("Authorization")
	orgId := csc.Param("orgId")
	userId := csc.Param("userId")
	log.Debugf("CreateServiceController Params: sessionId=%s, orgId=%s, userId=%s", sessionIdFromClient, orgId, userId)


	// Validate OrgId error
	csc.ValidateSession(sessionIdFromClient, orgId)
	if csc.CheckError() {
		return
	}

	// Parse data: service.
	cs := new(service.CreateService)
	err := csc.ReadJSON(cs)
	if err != nil {
		log.Errorf("CreateServiceController ReadJSON Error: error=%s", err)
		csc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if csc.CheckError() {
		return
	}

	// Get DcIdList
	csc.getApiServerList(cs.DcIdList)
	if csc.CheckError() {
		return
	}

	// Get K8sClient
	csc.createK8sClients()
	if csc.CheckError() {
		return
	}

	// Publish server to every datacenter
	orgName := cs.OrgName
	csc.createService(orgName, &cs.Service)
	if csc.CheckError() {
		return
	}

	// And NodePort to MySQL nodeport table
	op, _ := strconv.Atoi(userId)
	for _, v := range cs.Service.Spec.Ports {
		hasNodePort := mynodeport.PORT_START <= v.NodePort && v.NodePort <= mynodeport.PORT_LIMIT
		if hasNodePort {
			csc.createMysqlNodePort(hasNodePort, v.NodePort, cs.DcIdList, cs.Service.ObjectMeta.Name, int32(op))
			if csc.CheckError() {
				return
			}
		}
	}

	log.Infoln("CreateServiceController over!")
	csc.WriteOk("")
	return
}
