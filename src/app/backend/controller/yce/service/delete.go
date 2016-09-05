package service

import (
	mylog "app/backend/common/util/log"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	mynodeport "app/backend/model/mysql/nodeport"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	"strings"
)

type DeleteServiceController struct {
	*iris.Context
	Ye         *myerror.YceError
	k8sClients []*client.Client
	apiServers []string
}

func (dsc *DeleteServiceController) WriteBack() {
	dsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("DeleteServiceController Response YceError: controller=%p, code=%d, note=%s", dsc, dsc.Ye.Code, myerror.Errors[dsc.Ye.Code].LogMsg)
	dsc.Write(dsc.Ye.String())
}

// Validate Session
func (dsc *DeleteServiceController) validateSession(sessionIdFromClient, orgId string) {
	ss := session.SessionStoreInstance()

	// validate the session
	ok, err := ss.ValidateOrgId(sessionIdFromClient, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session Error: sessionId=%s, error=%s", sessionIdFromClient, err)
		dsc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Invalidate SessionId
	if !ok {
		mylog.Log.Errorf("Validate Session Failed: sessionId=%s, error=%s", sessionIdFromClient, err)
		dsc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
	}

	mylog.Log.Infof("DeleteServiceController validate sessionId successfully: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	return
}

// Get ApiServer by dcId
func (dsc *DeleteServiceController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		dsc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("DeleteServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer

}

// Get ApiServer List for dcIdList
func (dsc *DeleteServiceController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := dsc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("DeleteServiceController getApiServerList Error")
			return
		}

		dsc.apiServers = append(dsc.apiServers, apiServer)
	}
	mylog.Log.Infof("DeleteServiceController getApiServerList successfully: len(apiServers)=%d", len(dsc.apiServers))
	return
}

// Create k8sClient for every ApiServer
func (dsc *DeleteServiceController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	//dsc.k8sClients = make([]*client.Client, len(dsc.apiServers))
	dsc.k8sClients = make([]*client.Client, 0)

	for _, server := range dsc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			mylog.Log.Errorf("CreateK8sClient Error: error=%s", err)
			dsc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		dsc.k8sClients = append(dsc.k8sClients, c)
		// why??
		//dsc.apiServers = append(dsc.apiServers, server)
		mylog.Log.Infof("Append a new client to dsc.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	return
}

// Publish k8s.Service to every datacenter which in dcIdList
func (dsc *DeleteServiceController) deleteService(namespace, svcName string) {
	// Foreach every K8sClient to create service

	for index, cli := range dsc.k8sClients {
		err := cli.Services(namespace).Delete(svcName)
		if err != nil {
			mylog.Log.Errorf("deleteService Error: apiServer=%s, namespace=%s, error=%s", dsc.apiServers[index], namespace, err)
			//dsc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_SERVICE, "")
			dsc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_SERVICE, "")
			return
		}

		mylog.Log.Infof("Delete Service successfully: namespace=%s, apiServer=%s", namespace, dsc.apiServers[index])
	}

	return
}

func (dsc *DeleteServiceController) getOrgNameByOrgId(orgId string) string {
	org, err := myorganization.GetOrganizationById(orgId)
	if err != nil {
		mylog.Log.Errorf("deleteService Error: error=%s",err)
		dsc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_SERVICE, "")
		return ""
	}
	mylog.Log.Infof("DeleteServiceController getOrgNameByOrgId successfully: orgName=%s, orgId=%d", org.Name, orgId)
	return org.Name
}

// create NodePort(mysql) and insert it into db
func (dsc *DeleteServiceController) deleteMysqlNodePort(dcIdList []int32, nodePort, op int32) {
	for _, dcId := range dcIdList {


		np := &mynodeport.NodePort{
			Port:nodePort,
			DcId:dcId,
		}


		err := np.DeleteNodePort(op)
		if err != nil {
			mylog.Log.Errorf("DeleteMysqlNodePort Error: nodeport=%d, dcId=%d, svcName=%s, error=%s", np.Port, np.DcId, np.SvcName, err)
			dsc.Ye = myerror.NewYceError(myerror.EYCE_DELETE_NODEPORT, "")
			return
		}

		mylog.Log.Infof("DeleteMysqlNodePort Successfully: nodeport=%d, dcId=%d, svcName=%s", np.Port, np.DcId, np.SvcName)
	}

	return
}

func (dsc DeleteServiceController) Delete() {
	sessionIdFromClient := dsc.RequestHeader("Authorization")
	orgId := dsc.Param("orgId")
	dcId := dsc.Param("dcId")
	userId := dsc.Param("userId")
	svcName := dsc.Param("svcName")
	nodePort := dsc.RequestHeader("NodePort")

	// Validate OrgId error
	dsc.validateSession(sessionIdFromClient, orgId)
	if dsc.Ye != nil {
		dsc.WriteBack()
		return
	}

	// Get DcIdList
	dcIdList := make([]int32, 0)
	datacenterId, _ := strconv.Atoi(dcId)
	dcIdList = append(dcIdList, int32(datacenterId))

	dsc.getApiServerList(dcIdList)
	if dsc.Ye != nil {
		dsc.WriteBack()
		return
	}

	// Get K8sClient
	dsc.createK8sClients()
	if dsc.Ye != nil {
		dsc.WriteBack()
		return
	}

	// Publish server to every datacenter
	orgName := dsc.getOrgNameByOrgId(orgId)
	dsc.deleteService(orgName, svcName)
	if dsc.Ye != nil {
		dsc.WriteBack()
		return
	}

	// Update NodePort Status to MySQL nodeport table
	op, _ := strconv.Atoi(userId)
	port, _ := strconv.Atoi(nodePort)
	dsc.deleteMysqlNodePort(dcIdList, int32(port), int32(op))
	if dsc.Ye != nil {
		dsc.WriteBack()
		return
	}

	dsc.Ye = myerror.NewYceError(myerror.EOK, "")
	mylog.Log.Infoln("DeleteServiceController over!")
	dsc.WriteBack()
	return
}
