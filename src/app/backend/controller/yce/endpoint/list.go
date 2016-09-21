package endpoint

import (
	"app/backend/model/yce/endpoint"
	myerror "app/backend/common/yce/error"
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type ListEndpointsController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client
}


func (lec *ListEndpointsController) listEndpoints(namespace string, ed *endpoint.ListEndpoints) (epString string){
	epList := make([]endpoint.Endpoints, len(lec.apiServers))
	// Foreach every K8sClient to create service
	for index, cli := range lec.k8sClients {
		//_, err := cli.Services(namespace).Create(service)
		eps, err := cli.Endpoints(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("listEndpoints Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
			lec.Ye = myerror.NewYceError(myerror.EKUBE_LIST_ENDPOINTS, "")
			return
		}

		//TODO: check consistency
		epList[index].DcId = ed.DcIdList[index]
		epList[index].DcName = ed.DcName[index]
		epList[index].EndpointsList = *eps

		log.Infof("listEndpoints successfully: namespace=%s, apiServer=%s", namespace, lec.apiServers[index])
	}

	epJson, err := json.Marshal(epList)
	epString = string(epJson)
	if err != nil {
		log.Errorf("listEndpoints Error: apiServer=%v, namespace=%s, error=%s", lec.apiServers, namespace, err)
		lec.Ye = myerror.NewYceError(myerror.EKUBE_LIST_ENDPOINTS, "")
		return
	}

	return epString
}


//GET /api/v1/organizations/{orgId}/users/{userId}/endpoints
func (lec ListEndpointsController) Get() {
	sessionIdFromClient := lec.RequestHeader("Authorization")
	orgId := lec.Param("orgId")
	log.Debugf("ListEndpointsController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// ValidateSessionId
	lec.ValidateSession(sessionIdFromClient, orgId)
	if lec.CheckError() {
		return
	}

	// Get Datacenter List by orgId
	ed :=  new(endpoint.ListEndpoints)
	dcList, ye := yceutils.GetDatacenterListByOrgId(orgId)

	if ye != nil {
		lec.Ye = ye
	}

	ed.DcIdList = dcList.DcIdList
	ed.DcName = dcList.DcName
	if lec.CheckError() {
		return
	}

	// Get ApiServer List
	lec.apiServers, lec.Ye = yceutils.GetApiServerList(ed.DcIdList)
	if lec.CheckError() {
		return
	}

	// Create K8sClient
	lec.k8sClients, lec.Ye = yceutils.CreateK8sClientList(lec.apiServers)
	if lec.CheckError() {
		return
	}

	// List Endpoints
	orgName, ye := yceutils.GetOrgNameByOrgId(orgId)
	if ye != nil {
		lec.Ye = ye
	}

	epString := lec.listEndpoints(orgName, ed)
	if lec.CheckError() {
		return
	}

	lec.WriteOk(epString)
	log.Infoln("ListEndpointsController over!")
	return
}