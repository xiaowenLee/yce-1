package endpoint

import (
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	"github.com/kataras/iris"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"app/backend/common/util/session"
	"app/backend/model/yce/endpoint"
	"strconv"
	"strings"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/api"
)

type CreateEndpointsController struct {
	*iris.Context
	Ye *myerror.YceError
	k8sClients []*client.Client
	apiServers []string
}

func (cec *CreateEndpointsController) WriteBack() {
	cec.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateEndpointsController Response YceError: controller=%p, code=%d, note=%s", cec, cec.Ye.Code, myerror.Errors[cec.Ye.Code].LogMsg)
	cec.Write(cec.Ye.String())
}

// Validate Session
func (cec *CreateEndpointsController) validateSession(sessionIdFromClient, orgId string) {
	ss := session.SessionStoreInstance()

	// validate the session
	ok, err := ss.ValidateOrgId(sessionIdFromClient, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session Error: sessionId=%s, error=%s", sessionIdFromClient, err)
		cec.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Invalidate SessionId
	if !ok {
		mylog.Log.Errorf("Validate Session Failed: sessionId=%s, error=%s", sessionIdFromClient, err)
		cec.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
	}

	mylog.Log.Infof("CreateEndpointsController validate sessionId successfully: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	return
}

// Get ApiServer by dcId
func (cec *CreateEndpointsController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		cec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateEndpointsController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

// Get ApiServer List for dcIdList
func (cec *CreateEndpointsController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := cec.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("CreateEndpointsController getApiServerList Error")
			return
		}

		cec.apiServers = append(cec.apiServers, apiServer)
	}

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
			mylog.Log.Errorf("CreateK8sClient Error: error=%s", err)
			cec.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		cec.k8sClients = append(cec.k8sClients, c)
		// why??
		//cec.apiServers = append(cec.apiServers, server)
		mylog.Log.Infof("Append a new client to cec.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	return
}

// why need return value?
// Publish k8s.Service to every datacenter which in dcIdList
func (cec *CreateEndpointsController) createEndpoints(namespace string, endpoints *api.Endpoints) {
	// Foreach every K8sClient to create service
	for index, cli := range cec.k8sClients {
		_, err := cli.Endpoints(namespace).Create(endpoints)
		if err != nil {
			mylog.Log.Errorf("createEndpoints Error: apiServer=%s, namespace=%s, error=%s", cec.apiServers[index], namespace, err)
			cec.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_ENDPOINTS, "")
			return
		}

		mylog.Log.Infof("Create Endpoints successfully: namespace=%s, apiServer=%s", namespace, cec.apiServers[index])
	}

	return
}



func (cec CreateEndpointsController) Post() {
	sessionIdFromClient := cec.RequestHeader("Authorization")
	orgId := cec.Param("orgId")

	mylog.Log.Debugf("CreateEndpointsController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)


	// Validate OrgId error
	cec.validateSession(sessionIdFromClient, orgId)
	if cec.Ye != nil {
		cec.WriteBack()
		return
	}

	// Parse data: service.
	ce := new(endpoint.CreateEndpoints)
	cec.ReadJSON(ce)


	// Get DcIdList
	cec.getApiServerList(ce.DcIdList)
	if cec.Ye != nil {
		cec.WriteBack()
		return
	}

	// Get K8sClient
	cec.createK8sClients()
	if cec.Ye != nil {
		cec.WriteBack()
		return
	}

	// Publish server to every datacenter
	orgName := ce.OrgName
	cec.createEndpoints(orgName, &ce.Endpoints)
	if cec.Ye != nil {
		cec.WriteBack()
		return
	}



	cec.Ye = myerror.NewYceError(myerror.EOK, "")
	mylog.Log.Infoln("CreateEndpointsController over!")
	cec.WriteBack()
	return
}