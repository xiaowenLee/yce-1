package endpoint
import (
	mylog "app/backend/common/util/log"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	"strings"
)

type DeleteEndpointsController struct {
	*iris.Context
	Ye         *myerror.YceError
	k8sClients []*client.Client
	apiServers []string
}

func (dec *DeleteEndpointsController) WriteBack() {
	dec.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("DeleteEndpointsController Response YceError: controller=%p, code=%d, note=%s", dec, dec.Ye.Code, myerror.Errors[dec.Ye.Code].LogMsg)
	dec.Write(dec.Ye.String())
}

// Validate Session
func (dec *DeleteEndpointsController) validateSession(sessionIdFromClient, orgId string) {
	ss := session.SessionStoreInstance()

	// validate the session
	ok, err := ss.ValidateOrgId(sessionIdFromClient, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session Error: sessionId=%s, error=%s", sessionIdFromClient, err)
		dec.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Invalidate SessionId
	if !ok {
		mylog.Log.Errorf("Validate Session Failed: sessionId=%s, error=%s", sessionIdFromClient, err)
		dec.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
	}

	mylog.Log.Infof("DeleteEndpointsController validate sessionId successfully: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	return
}

// Get ApiServer by dcId
func (dec *DeleteEndpointsController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		dec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("DeleteEndpointsController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer

}

// Get ApiServer List for dcIdList
func (dec *DeleteEndpointsController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := dec.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("DeleteEndpointsController getApiServerList Error")
			return
		}

		dec.apiServers = append(dec.apiServers, apiServer)
	}
	mylog.Log.Infof("DeleteEndpointsController getApiServerList successfully: len(apiServers)=%d", len(dec.apiServers))
	return
}

// Create k8sClient for every ApiServer
func (dec *DeleteEndpointsController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	//dec.k8sClients = make([]*client.Client, len(dec.apiServers))
	dec.k8sClients = make([]*client.Client, 0)

	for _, server := range dec.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			mylog.Log.Errorf("CreateK8sClient Error: error=%s", err)
			dec.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		dec.k8sClients = append(dec.k8sClients, c)
		// why??
		//dec.apiServers = append(dec.apiServers, server)
		mylog.Log.Infof("Append a new client to dec.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	mylog.Log.Infof("DeleteEndpointsController createK8sClient success: len(createK8sClient)=%d", len(dec.k8sClients))
	return
}

// Publish k8s.Endpoint to every datacenter which in dcIdList
func (dec *DeleteEndpointsController) deleteEndpoints(namespace, epName string) {
	// Foreach every K8sClient to create service

	for index, cli := range dec.k8sClients {
		err := cli.Endpoints(namespace).Delete(epName)
		if err != nil {
			mylog.Log.Errorf("deleteEndpoint Error: apiServer=%s, namespace=%s, error=%s", dec.apiServers[index], namespace, err)
			//dec.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_SERVICE, "")
			dec.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_ENDPOINT, "")
			return
		}

		mylog.Log.Infof("Delete Endpoint successfully: namespace=%s, apiServer=%s", namespace, dec.apiServers[index])
	}
	mylog.Log.Infof("DeleteEndpointsController Delete Endpoints success")
	return
}

func (dec *DeleteEndpointsController) getOrgNameByOrgId(orgId string) string {
	org, err := myorganization.GetOrganizationById(orgId)
	if err != nil {
		mylog.Log.Errorf("deleteEndpoint Error: error=%s",err)
		dec.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_SERVICE, "")
		return ""
	}
	mylog.Log.Infof("DeleteEndpointsController getOrgNameByOrgId successfully: orgName=%s, orgId=%d", org.Name, orgId)
	return org.Name
}



func (dec DeleteEndpointsController) Delete() {
	sessionIdFromClient := dec.RequestHeader("Authorization")
	orgId := dec.Param("orgId")
	dcId := dec.Param("dcId")
	epName := dec.Param("epName")

	mylog.Log.Debugf("DeleteEndpontsController Params: sessionId=%s, orgId=%s, dcId=%s, epName=%s", sessionIdFromClient, orgId, dcId, epName)


	// Validate OrgId error
	dec.validateSession(sessionIdFromClient, orgId)
	if dec.Ye != nil {
		dec.WriteBack()
		return
	}

	// Get DcIdList
	dcIdList := make([]int32, 0)
	datacenterId, _ := strconv.Atoi(dcId)
	dcIdList = append(dcIdList, int32(datacenterId))
	mylog.Log.Debugf("DeleteEndpointController len(DcIdList)=%d", len(dcIdList))

	dec.getApiServerList(dcIdList)
	if dec.Ye != nil {
		dec.WriteBack()
		return
	}

	// Get K8sClient
	dec.createK8sClients()
	if dec.Ye != nil {
		dec.WriteBack()
		return
	}

	// Publish server to every datacenter
	orgName := dec.getOrgNameByOrgId(orgId)
	if dec.Ye != nil {
		dec.WriteBack()
		return
	}

	dec.deleteEndpoints(orgName, epName)
	if dec.Ye != nil {
		dec.WriteBack()
		return
	}

	dec.Ye = myerror.NewYceError(myerror.EOK, "")
	mylog.Log.Infoln("DeleteEndpointsController over!")
	dec.WriteBack()
	return
}
