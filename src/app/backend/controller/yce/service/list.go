package service

import (
	mylog "app/backend/common/util/log"
	"github.com/kataras/iris"
	"app/backend/common/util/session"
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
)

type ListServiceController struct {
	*iris.Context
	apiServers []string
	k8sClients []*client.Client
	Ye *myerror.YceError
}

const (
	SELECT_USER = "SELECT name FROM user WHERE id=?"
)

func (lsc *ListServiceController) WriteBack() {
	lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("Create ListServiceController Response Error: controller=%p, code=%d, note=%s", lsc, lsc.Ye.Code, myerror.Errors[lsc.Ye.Code].LogMsg)
	lsc.Write(lsc.Ye.String())
}

func (lsc *ListServiceController) validateSessionId(sessionId, orgId string) {
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	// validate error
	if err != nil {
		mylog.Log.Errorf("Create ListServiceController Error: sessionId=%s, orgId=%s, error=%s", sessionId, orgId, err)
		lsc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// invalid sessionId
	if !ok {
		mylog.Log.Errorf("Create ListServiceController Failed: sessionId=%s, orgId=%s", sessionId, orgId)
		lsc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	mylog.Log.Infof("ListServiceController validate sessionId success")
	return
}


func (lsc *ListServiceController) getDatacentersByOrgId(sd *service.ListService, orgId string) {
	org, err := organization.GetOrganizationById(orgId)
	sd.Organization = org
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lsc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	dcList, err := organization.GetDataCentersByOrganization(sd.Organization)
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		lsc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return
	}

	sd.DcIdList = make([]int32, len(dcList))
	sd.DcName = make([]string, len(dcList))

	for index, dc := range dcList {
		sd.DcIdList[index] = dc.Id
		sd.DcName[index] = dc.Name
	}

	mylog.Log.Infof("ListServiceController getDatacentersByOrgId: len(dcIdList)=%d, len(dcName)=%d", len(sd.DcIdList), len(sd.DcName))
}


// Get ApiServer by dcId
func (lsc *ListServiceController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		lsc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

func (lsc *ListServiceController) getApiServerList(dcIdList []int32) {
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := lsc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("ListServiceController getApiServerList Error")
			return
		}

		lsc.apiServers = append(lsc.apiServers, apiServer)
	}

	mylog.Log.Infof("CreateServiceController getApiServerList: len(apiServer)=%d", len(lsc.apiServers))
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
			mylog.Log.Errorf("CreateK8sClient Error: error=%s", err)
			lsc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		lsc.k8sClients = append(lsc.k8sClients, c)
		// why??
		//lsc.apiServers = append(lsc.apiServers, server)
		mylog.Log.Infof("Append a new client to lsc.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	mylog.Log.Infof("CreateServiceController createK8sClients: len(k8sClients)=%d", len(lsc.k8sClients))
	return
}


func (lsc *ListServiceController) listService(namespace string, sd *service.ListService) (svcString string){
	svcList := make([]service.Service, len(lsc.apiServers))
	// Foreach every K8sClient to create service
	for index, cli := range lsc.k8sClients {
		svcs, err := cli.Services(namespace).List(api.ListOptions{})
		if err != nil {
			mylog.Log.Errorf("listService Error: apiServer=%s, namespace=%s, error=%s", lsc.apiServers[index], namespace, err)
			lsc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_SERVICE, "")
			return
		}

		//TODO: check consistency


		svcList[index].DcId = sd.DcIdList[index]
		svcList[index].DcName = sd.DcName[index]
		svcList[index].ServiceList = *svcs

		mylog.Log.Infof("listService successfully: namespace=%s, apiServer=%s", namespace, lsc.apiServers[index])

	}

	svcJson, err := json.Marshal(svcList)
	svcString = string(svcJson)
	if err != nil {
		mylog.Log.Errorf("listService Error: apiServer=%v, namespace=%s, error=%s", lsc.apiServers, namespace, err)
		lsc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_SERVICE, "")
		return
	}

	return svcString
}


//GET /api/v1/organizations/{orgId}/users/{userId}/endpoints
func (lsc ListServiceController) Get() {
	sessionIdFromClient := lsc.RequestHeader("Authorization")
	orgId := lsc.Param("orgId")

	mylog.Log.Debugf("ListServiceController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// validateSessionId
	lsc.validateSessionId(sessionIdFromClient, orgId)
	if lsc.Ye != nil {
		lsc.WriteBack()
		return
	}


	// Get Datacenters by organizations
	//ed :=  new(endpoint.ListEndpoints)
	sd := new(service.ListService)
	lsc.getDatacentersByOrgId(sd, orgId)
	if lsc.Ye != nil {
		lsc.WriteBack()
		return
	}


	// Get ApiServers by organizations
	lsc.getApiServerList(sd.DcIdList)
	if lsc.Ye != nil {
		lsc.WriteBack()
		return
	}

	// Get K8sClient
	lsc.createK8sClients()
	if lsc.Ye != nil {
		lsc.WriteBack()
		return
	}

	// List Endpoints
	orgName := sd.Organization.Name
	svcString := lsc.listService(orgName, sd)
	if lsc.Ye != nil {
		lsc.WriteBack()
		return
	}

	lsc.Ye = myerror.NewYceError(myerror.EOK, svcString)
	lsc.WriteBack()

	mylog.Log.Infoln("ListServiceController over!")

	return
}
