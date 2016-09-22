package service

import (
	"app/backend/model/yce/service"
	myerror "app/backend/common/yce/error"
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type ListServiceController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client
}


func (lsc *ListServiceController) listService(namespace string, sd *service.ListService) (svcString string){
	//TODO: use slice instead of array
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


//GET /api/v1/organizations/{orgId}/users/{userId}/services
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
	sd := new(service.ListService)

	dcList, ye := yceutils.GetDatacenterListByOrgId(orgId)
	if ye != nil {
		lsc.Ye = ye
	}
	if lsc.CheckError() {
		return
	}


	sd.DcIdList = dcList.DcIdList
	sd.DcName = dcList.DcName

	// Get ApiServers by DcIdList
	lsc.apiServers, lsc.Ye = yceutils.GetApiServerList(sd.DcIdList)
	if lsc.CheckError() {
		return
	}

	// Get K8sClient
	lsc.k8sClients, lsc.Ye = yceutils.CreateK8sClientList(lsc.apiServers)
	if lsc.CheckError() {
		return
	}

	//get OrgName by OrgId
	orgName, ye := yceutils.GetOrgNameByOrgId(orgId)
	if ye != nil {
		lsc.Ye = ye
	}
	if lsc.CheckError() {
		return
	}

	// List Services
	svcString := lsc.listService(orgName, sd)
	if lsc.CheckError() {
		return
	}

	lsc.WriteOk(svcString)
	log.Infoln("ListServiceController over!")
	return
}
