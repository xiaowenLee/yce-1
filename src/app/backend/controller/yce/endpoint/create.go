package endpoint

import (
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"app/backend/model/yce/endpoint"
	"strconv"
	"strings"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/api"
	yce "app/backend/controller/yce"
)

type CreateEndpointsController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
}


// Get ApiServer by dcId
func (cec *CreateEndpointsController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		cec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("CreateEndpointsController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

// Get ApiServer List for dcIdList
func (cec *CreateEndpointsController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := cec.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			log.Errorf("CreateEndpointsController getApiServerList Error")
			return
		}

		cec.apiServers = append(cec.apiServers, apiServer)
	}

	log.Infof("CreateEndpointsController getApiServerList success: len(apiServer)=%d", len(cec.apiServers))
	return
}

// Create k8sClient for every ApiServer
func (cec *CreateEndpointsController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	cec.k8sClients = make([]*client.Client, 0)


	for _, server := range cec.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			log.Errorf("CreateK8sClient Error: error=%s", err)
			cec.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		cec.k8sClients = append(cec.k8sClients, c)
		// why??
		//cec.apiServers = append(cec.apiServers, server)
		log.Infof("Append a new client to cec.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	log.Infof("CreateEndpointsController createK8sClient success: len(k8sclient)=%d", len(cec.k8sClients))
	return
}

// why need return value?
// Publish k8s.Service to every datacenter which in dcIdList
func (cec *CreateEndpointsController) createEndpoints(namespace string, endpoints *api.Endpoints) {
	// Foreach every K8sClient to create service
	for index, cli := range cec.k8sClients {
		_, err := cli.Endpoints(namespace).Create(endpoints)
		if err != nil {
			log.Errorf("createEndpoints Error: apiServer=%s, namespace=%s, error=%s", cec.apiServers[index], namespace, err)
			cec.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_ENDPOINTS, "")
			return
		}

		log.Infof("Create Endpoints successfully: name=%s, namespace=%s, apiServer=%s", endpoints.Name, namespace, cec.apiServers[index])
	}

	log.Infof("CreateEndpointsController createEndpoints success")
	return
}



func (cec CreateEndpointsController) Post() {
	sessionIdFromClient := cec.RequestHeader("Authorization")
	orgId := cec.Param("orgId")

	log.Debugf("CreateEndpointsController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)


	// Validate OrgId error
	cec.ValidateSession(sessionIdFromClient, orgId)
	if cec.CheckError() {
		return
	}

	// Parse data: service.
	ce := new(endpoint.CreateEndpoints)
	err := cec.ReadJSON(ce)
	if err != nil {
		log.Debugf("CreateEndpointsController ReadJSON Error: error=%s", err)
		cec.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cec.CheckError() {
		return
	}


	// Get DcIdList
	cec.getApiServerList(ce.DcIdList)
	if cec.CheckError() {
		return
	}

	// Get K8sClient
	cec.createK8sClients()
	if cec.CheckError() {
		return
	}

	// Publish server to every datacenter
	orgName := ce.OrgName
	cec.createEndpoints(orgName, &ce.Endpoints)
	if cec.CheckError() {
		return
	}

	cec.WriteOk("")
	log.Infoln("CreateEndpointsController over!")
	return
}