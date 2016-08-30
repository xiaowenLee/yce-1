package extensions

import (
	"github.com/kataras/iris"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	mylog "app/backend/common/util/log"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	"app/backend/common/yce/organization"
	"strings"
	"encoding/json"

	"app/backend/model/yce/extensions"
	"strconv"
	"k8s.io/kubernetes/pkg/api"
)

type ListExtensionsController struct {
	*iris.Context
	Ye *myerror.YceError
	apiServers []string
	k8sClients []*client.Client
}

func (lec *ListExtensionsController) WriteBack() {
	lec.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("Create ListExtensionsController Response Error: controller=%p, code=%d, note=%s", lec, lec.Ye.Code, myerror.Errors[lec.Ye.Code].LogMsg)
	lec.Write(lec.Ye.String())
}

func (lec *ListExtensionsController) validateSessionId(sessionId, orgId string) {

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


func (lec *ListExtensionsController) getDatacentersByOrgId(le *extensions.ListExtensions, orgId string) {
	org, err := organization.GetOrganizationById(orgId)
	le.Organization = org
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	dcList, err := organization.GetDataCentersByOrganization(le.Organization)
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return
	}

	le.DcIdList = make([]int32, len(dcList))
	le.DcName = make([]string, len(dcList))

	for index, dc := range dcList {
		le.DcIdList[index] = dc.Id
		le.DcName[index] = dc.Name
	}

}


// Get ApiServer by dcId
func (lec *ListExtensionsController) getApiServerByDcId(dcId int32) string {
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

func (lec *ListExtensionsController) getApiServerList(dcIdList []int32) {
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := lec.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("ListExtensionsController getApiServerList Error")
			return
		}

		lec.apiServers = append(lec.apiServers, apiServer)
	}

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
			mylog.Log.Errorf("CreateK8sClient Error: error=%s", err)
			lec.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		lec.k8sClients = append(lec.k8sClients, c)
		// why??
		//lec.apiServers = append(lec.apiServers, server)
		mylog.Log.Infof("Append a new client to lec.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	return
}

func (lec *ListExtensionsController) listServiceAndEndpoints(namespace string, le *extensions.ListExtensions) (extString string){
	extList := make([]extensions.Extensions, len(lec.apiServers))

	for index, cli := range lec.k8sClients {
		svcs, err := cli.Services(namespace).List(api.ListOptions{})
		if err != nil {
			mylog.Log.Errorf("listService Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
			myerror.NewYceError(myerror.EYCE_LIST_EXTENSIONS, "")
			return
		}

		eps, err := cli.Endpoints(namespace).List(api.ListOptions{})
		if err != nil {
			mylog.Log.Errorf("listEndpoints Error: apiServer=%s, namespace=%s, error=%s", lec.apiServers[index], namespace, err)
			myerror.NewYceError(myerror.EYCE_LIST_EXTENSIONS, "")
			return
		}

		extList[index].DcId = le.DcIdList[index]
		extList[index].DcName = le.DcName[index]
		extList[index].ServiceList = *svcs
		extList[index].EndpointList = *eps

		mylog.Log.Infof("list Service and Endpoints Successfully: namespace=%s, apiServer=%s", namespace, lec.apiServers)
	}

	extJson, err := json.Marshal(extList)
	extString = string(extJson)
	if err != nil {
		mylog.Log.Errorf("list Service and Endpoints Error: apiServer=%v, namespace=%s, error=%s", lec.apiServers, namespace, err)
		lec.Ye = myerror.NewYceError(myerror.EYCE_LIST_EXTENSIONS, "")
		return
	}

	return extString

}

func (lec ListExtensionsController) Get() {
	sessionIdFromClient := lec.RequestHeader("Authorization")
	orgId := lec.Param("orgId")

	// validateSessionId
	lec.validateSessionId(sessionIdFromClient, orgId)
	if lec.Ye != nil {
		lec.WriteBack()
		return
	}


	// Get Datacenters by organizations
	le := new(extensions.ListExtensions)
	lec.getDatacentersByOrgId(le, orgId)
	if lec.Ye != nil {
		lec.WriteBack()
		return
	}


	// Get ApiServers by organizations
	lec.getApiServerList(le.DcIdList)
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

	// Get ServiceList and Endpoint
	orgName := le.Organization.Name
	extString := lec.listServiceAndEndpoints(orgName, le)

	if lec.Ye != nil {
		lec.WriteBack()
		return
	}

	lec.Ye = myerror.NewYceError(myerror.EOK, extString)
	lec.WriteBack()

	mylog.Log.Infoln("ListExtensionsController over!")

	return
}

