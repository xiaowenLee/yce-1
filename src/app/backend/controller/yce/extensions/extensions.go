package extensions

import (
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	"app/backend/common/yce/organization"
	"strings"
	"encoding/json"
	"app/backend/model/yce/extensions"
	"strconv"
	"k8s.io/kubernetes/pkg/api"
	yce "app/backend/controller/yce"
)

type ListExtensionsController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client
}


func (lec *ListExtensionsController) getDatacentersByOrgId(le *extensions.ListExtensions, orgId string) {
	org, err := organization.GetOrganizationById(orgId)
	le.Organization = org
	if err != nil {
		log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	dcList, err := organization.GetDataCentersByOrganization(le.Organization)
	if err != nil {
		log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return
	}

	le.DcIdList = make([]int32, len(dcList))
	le.DcName = make([]string, len(dcList))

	for index, dc := range dcList {
		le.DcIdList[index] = dc.Id
		le.DcName[index] = dc.Name
	}

	log.Debugf("ListExtensionsController getDatacentersByOrgId: len(dcIdList)=%d, len(dcName)=%d", len(le.DcIdList), len(le.DcName))

}


// Get ApiServer by dcId
func (lec *ListExtensionsController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		lec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("CreateServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

func (lec *ListExtensionsController) getApiServerList(dcIdList []int32) {
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := lec.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			log.Errorf("ListExtensionsController getApiServerList Error")
			return
		}

		lec.apiServers = append(lec.apiServers, apiServer)
	}

	log.Infof("ListExtensionsController getApiServerList: len(apiServer)=%d", len(lec.apiServers))
	return
}


func (lec *ListExtensionsController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	//lec.k8sClients = make([]*client.Client, len(lec.apiServers))
	lec.k8sClients = make([]*client.Client, 0)


	for _, server := range lec.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			log.Errorf("CreateK8sClient Error: error=%s", err)
			lec.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		lec.k8sClients = append(lec.k8sClients, c)
		// why??
		//lec.apiServers = append(lec.apiServers, server)
		log.Infof("Append a new client to lec.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	log.Infof("ListExtensionsController createK8sClient: len(k8sClient)=%d", len(lec.k8sClients))
	return
}

func (lec *ListExtensionsController) listServiceAndEndpoints(namespace string, le *extensions.ListExtensions) (extString string){
	extList := make([]extensions.Extensions, len(lec.apiServers))

	for index, cli := range lec.k8sClients {
		svcs, err := cli.Services(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("listService Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
			myerror.NewYceError(myerror.EYCE_LIST_EXTENSIONS, "")
			return
		}

		eps, err := cli.Endpoints(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("listEndpoints Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
			myerror.NewYceError(myerror.EYCE_LIST_EXTENSIONS, "")
			return
		}

		extList[index].DcId = le.DcIdList[index]
		extList[index].DcName = le.DcName[index]
		extList[index].ServiceList = *svcs
		extList[index].EndpointList = *eps

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

	// Get Datacenters by organizations
	le := new(extensions.ListExtensions)
	lec.getDatacentersByOrgId(le, orgId)
	if lec.CheckError() {
		return
	}


	// Get ApiServers by organizations
	lec.getApiServerList(le.DcIdList)
	if lec.CheckError() {
		return
	}

	// Get K8sClient
	lec.createK8sClients()
	if lec.CheckError() {
		return
	}

	// Get ServiceList and Endpoint
	orgName := le.Organization.Name
	extString := lec.listServiceAndEndpoints(orgName, le)
	if lec.CheckError() {
		return
	}

	lec.WriteOk(extString)
	log.Infoln("ListExtensionsController over!")

	return
}

