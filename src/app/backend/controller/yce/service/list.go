package service

import (
	"app/backend/common/yce/organization"
	"app/backend/model/yce/service"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	"encoding/json"
	"strconv"
	"strings"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	yce "app/backend/controller/yce"
)

type ListServiceController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client
}

func (lsc *ListServiceController) getDatacentersByOrgId(sd *service.ListService, orgId string) {
	org, err := organization.GetOrganizationById(orgId)
	sd.Organization = org
	if err != nil {
		log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lsc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	dcList, err := organization.GetDataCentersByOrganization(sd.Organization)
	if err != nil {
		log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lsc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return
	}

	sd.DcIdList = make([]int32, len(dcList))
	sd.DcName = make([]string, len(dcList))

	for index, dc := range dcList {
		sd.DcIdList[index] = dc.Id
		sd.DcName[index] = dc.Name
	}

	log.Infof("ListServiceController getDatacentersByOrgId: len(dcIdList)=%d, len(dcName)=%d", len(sd.DcIdList), len(sd.DcName))
}


// Get ApiServer by dcId
func (lsc *ListServiceController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		lsc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("CreateServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

func (lsc *ListServiceController) getApiServerList(dcIdList []int32) {
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := lsc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			log.Errorf("ListServiceController getApiServerList Error")
			return
		}

		lsc.apiServers = append(lsc.apiServers, apiServer)
	}

	log.Infof("CreateServiceController getApiServerList: len(apiServer)=%d", len(lsc.apiServers))
	return
}


func (lsc *ListServiceController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	//lsc.k8sClients = make([]*client.Client, len(lsc.apiServers))
	lsc.k8sClients = make([]*client.Client, 0)


	for _, server := range lsc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			log.Errorf("CreateK8sClient Error: error=%s", err)
			lsc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		lsc.k8sClients = append(lsc.k8sClients, c)
		// why??
		//lsc.apiServers = append(lsc.apiServers, server)
		log.Infof("Append a new client to lsc.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	log.Infof("CreateServiceController createK8sClients: len(k8sClients)=%d", len(lsc.k8sClients))
	return
}


func (lsc *ListServiceController) listService(namespace string, sd *service.ListService) (svcString string){
	svcList := make([]service.Service, len(lsc.apiServers))
	// Foreach every K8sClient to create service
	for index, cli := range lsc.k8sClients {
		svcs, err := cli.Services(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("listService Error: apiServer=%s, namespace=%s, error=%s", lsc.apiServers[index], namespace, err)
			lsc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_SERVICE, "")
			return
		}

		//TODO: check consistency


		svcList[index].DcId = sd.DcIdList[index]
		svcList[index].DcName = sd.DcName[index]
		svcList[index].ServiceList = *svcs

		log.Infof("listService successfully: namespace=%s, apiServer=%s", namespace, lsc.apiServers[index])

	}

	svcJson, err := json.Marshal(svcList)
	svcString = string(svcJson)
	if err != nil {
		log.Errorf("listService Error: apiServer=%v, namespace=%s, error=%s", lsc.apiServers, namespace, err)
		lsc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_SERVICE, "")
		return
	}

	return svcString
}


//GET /api/v1/organizations/{orgId}/users/{userId}/endpoints
func (lsc ListServiceController) Get() {
	sessionIdFromClient := lsc.RequestHeader("Authorization")
	orgId := lsc.Param("orgId")

	log.Debugf("ListServiceController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// ValidateSession
	lsc.ValidateSession(sessionIdFromClient, orgId)
	if lsc.CheckError() {
		return
	}


	// Get Datacenters by organizations
	//ed :=  new(endpoint.ListEndpoints)
	sd := new(service.ListService)
	lsc.getDatacentersByOrgId(sd, orgId)
	if lsc.CheckError() {
		return
	}


	// Get ApiServers by organizations
	lsc.getApiServerList(sd.DcIdList)
	if lsc.CheckError() {
		return
	}

	// Get K8sClient
	lsc.createK8sClients()
	if lsc.CheckError() {
		return
	}

	// List Endpoints
	orgName := sd.Organization.Name
	svcString := lsc.listService(orgName, sd)
	if lsc.CheckError() {
		return
	}

	lsc.WriteOk(svcString)
	log.Infoln("ListServiceController over!")
	return
}
