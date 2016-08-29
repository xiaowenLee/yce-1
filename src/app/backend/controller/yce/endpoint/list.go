package endpoint

import (
	"github.com/kataras/iris"
	"app/backend/common/util/session"
	"app/backend/common/yce/organization"
	"app/backend/model/yce/endpoint"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	mylog "app/backend/common/util/log"
	"encoding/json"
	"strconv"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strings"
)

type ListEndpointController struct {
	*iris.Context
	apiServers []string
	k8sClients []client.Client
	Ye *myerror.YceError
}

func (lec *ListEndpointController) WriteBack() {
	lec.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("Create ListEndpointController Response Error: controller=%p, code=%d, note=%s", lec, lec.Ye.Code, myerror.Errors[lec.Ye.Code].LogMsg)
	lec.Write(lec.Ye.String())
}

func (lec *ListEndpointController) validateSessionId(sessionId, orgId string) {
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	// validate error
	if err != nil {
		mylog.Log.Errorf("Create ListEndpointController Error: sessionId=%s, orgId=%s, error=%s", sessionId, orgId, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// invalid sessionId
	if !ok {
		mylog.Log.Errorf("Create ListEndpoint Controller Failed: sessionId=%s, orgId=%s", sessionId, orgId)
		lec.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	return
}


func (lec *ListEndpointController) getDatacentersByOrgId(ed endpoint.ListEndpoints, orgId string) {
	org, err := organization.GetOrganizationById(orgId)
	ed.Organization = org
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	dcList, err := organization.GetDataCentersByOrganization(ed.Organization)
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return
	}

	for index, dc := range dcList {
		ed.DcIdList[index] = dc.Id
		ed.DcName[index] = dc.Name
	}

}


// Get ApiServer by dcId
func (lec *ListEndpointController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		lec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

func (lec *ListEndpointController) getApiServerList(dcIdList []int32) {
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := lec.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("ListEndpointController getApiServerList Error")
			return
		}

		lec.apiServers = append(lec.apiServers, apiServer)
	}

	return
}


func (lec *ListEndpointController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	lec.k8sClients := make([]*client.Client, len(lec.apiServers))


	for _, server := range lec.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			mylog.Log.Errorf("CreateK8sClient Error: error=%s", err)
			lec.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		lec.k8sClients = append(lec.k8sClients, c)
		// why??
		//lec.apiServers = append(lec.apiServers, server)
		mylog.Log.Infof("Append a new client to lec.K8sClients array: c=%p, apiServer=%s", c, server})
	}

	return
}

func (lec *ListEndpointController) listEndpoints(namespace string, ed endpoint.ListEndpoints) (epString string){
	epList := make([]endpoint.Endpoints, len(lec.apiServers))
	// Foreach every K8sClient to create service
	for index, cli := range lec.k8sClients {
		//_, err := cli.Services(namespace).Create(service)
		eps, err := cli.Endpoints(namespace).List(api.ListOptions{})
		if err != nil {
			mylog.Log.Errorf("listEndpoints Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
			lec.Ye = myerror.NewYceError(myerror.EKUBE_LIST_ENDPOINTS, "")
			return
		}

		//TODO: check consistency
		epList[index].DcId = ed.DcIdList[index]
		epList[index].DcName = ed.DcName[index]
		epList[index].EndpointsList = eps

		mylog.Log.Infof("listEndpoints successfully: namespace=%s, apiServer=%s", namespace, lec.apiServers[index])

	}

	epString, err := json.Marshal(epList)
	if err != nil {
		mylog.Log.Errorf("listEndpoints Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
		lec.Ye = myerror.NewYceError(myerror.EKUBE_LIST_ENDPOINTS, "")
		return
	}

	return
}


//GET /api/v1/organizations/{orgId}/users/{userId}/endpoints
func (lec ListEndpointController) Get() {
	sessionIdFromClient := lec.RequestHeader("Authorization")
	orgId := lec.Param("orgId")
	userId := lec.Params("userId")

	// validateSessionId
	lec.validateSessionId(sessionIdFromClient, orgId)
	if lec.Ye != nil {
		lec.WriteBack()
		return
	}


	// Get Datacenters by organizations
	ed :=  new(endpoint.ListEndpoints)
	lec.getDatacentersByOrgId(ed, orgId)
	if lec.Ye != nil {
		lec.WriteBack()
		return
	}


	// Get ApiServers by organizations
	lec.getApiServerList(ed.DcIdList)
	if lec.Ye != nil {
		lec.WriteBack()
		return
	}

	// Get K8sClient
	lec.createK8sClients()
	if lec.Ye != nil {
		lec.WriteBack()
		return
	}

	// List Endpoints
	orgName := ed.Organization.Name
	epString := lec.listEndpoints(orgName, ed)
	if lec.Ye != nil {
		lec.WriteBack()
		return
	}

	lec.Ye = myerror.NewYceError(myerror.EOK, "", epString)
	lec.WriteBack()

	mylog.Log.Infoln("ListEndpointController over!")

	return
}