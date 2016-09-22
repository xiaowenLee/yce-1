package extensions

import (
	myerror "app/backend/common/yce/error"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"encoding/json"
	"app/backend/model/yce/extensions"
	"k8s.io/kubernetes/pkg/api"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type ListExtensionsController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client
}

//TODO: for condition that we cann't delete the endpoints which was created by the service with selector, now cancel the endpoints deletion
func (lec *ListExtensionsController) listServiceAndEndpoints(namespace string, le *yceutils.DatacenterList) (extString string){
	extList := make([]extensions.Extensions, len(lec.apiServers))

	for index, cli := range lec.k8sClients {
		svcs, err := cli.Services(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("listService Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
			myerror.NewYceError(myerror.EYCE_LIST_EXTENSIONS, "")
			return
		}

		/*
		eps, err := cli.Endpoints(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("listEndpoints Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
			myerror.NewYceError(myerror.EYCE_LIST_EXTENSIONS, "")
			return
		}
		*/

		extList[index].DcId = le.DcIdList[index]
		extList[index].DcName = le.DcName[index]
		extList[index].ServiceList = *svcs
		//extList[index].EndpointList = *eps

		log.Infof("list Service and Endpoints Successfully: namespace=%s, apiServer=%s", namespace, lec.apiServers)
	}

	extJson, err := json.Marshal(extList)
	extString = string(extJson)
	if err != nil {
		log.Errorf("list Service and Endpoints Error: apiServer=%v, namespace=%s, error=%s", lec.apiServers, namespace, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_LIST_EXTENSIONS, "")
		return
	}

	return extString

}

func (lec ListExtensionsController) Get() {
	sessionIdFromClient := lec.RequestHeader("Authorization")
	orgId := lec.Param("orgId")
	log.Debugf("ListExtensionsController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// ValidateSessionId
	lec.ValidateSession(sessionIdFromClient, orgId)
	if lec.CheckError() {
		return
	}

	// Get Datacenter List by orgId
	le, ye := yceutils.GetDatacenterListByOrgId(orgId)
	if ye != nil {
		lec.Ye = ye
	}

	if lec.CheckError() {
		return
	}

	// Get ApiServers by DcIdList
	lec.apiServers, lec.Ye = yceutils.GetApiServerList(le.DcIdList)
	if lec.CheckError() {
		return
	}

	// Create K8sClient
	lec.k8sClients, lec.Ye = yceutils.CreateK8sClientList(lec.apiServers)
	if lec.CheckError() {
		return
	}

	// Get ServiceList and Endpoint
	orgName, ye := yceutils.GetOrgNameByOrgId(orgId)
	if ye != nil {
		lec.Ye = ye
	}
	extString := lec.listServiceAndEndpoints(orgName, le)
	if lec.CheckError() {
		return
	}

	lec.WriteOk(extString)
	log.Infoln("ListExtensionsController over!")

	return
}

