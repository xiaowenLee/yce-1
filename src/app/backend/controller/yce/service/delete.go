package service

import (
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	mynodeport "app/backend/model/mysql/nodeport"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	"strings"
	"app/backend/model/yce/service"
	yce "app/backend/controller/yce"
)

type DeleteServiceController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
}

// Get ApiServer by dcId
func (dsc *DeleteServiceController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		dsc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("DeleteServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer

}

// Get ApiServer List for dcIdList
func (dsc *DeleteServiceController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := dsc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			log.Errorf("DeleteServiceController getApiServerList Error")
			return
		}

		dsc.apiServers = append(dsc.apiServers, apiServer)
	}
	log.Infof("DeleteServiceController getApiServerList successfully: len(apiServers)=%d", len(dsc.apiServers))
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
			log.Errorf("CreateK8sClient Error: error=%s", err)
			dsc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		dsc.k8sClients = append(dsc.k8sClients, c)
		// why??
		//dsc.apiServers = append(dsc.apiServers, server)
		log.Infof("Append a new client to dsc.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	log.Infof("DeleteServiceController createK8sClient: len(k8sClients)=%d", len(dsc.k8sClients))
	return
}

// Publish k8s.Service to every datacenter which in dcIdList
func (dsc *DeleteServiceController) deleteService(namespace, svcName string) {
	// Foreach every K8sClient to create service

	for index, cli := range dsc.k8sClients {
		err := cli.Services(namespace).Delete(svcName)
		if err != nil {
			log.Errorf("deleteService Error: apiServer=%s, namespace=%s, error=%s", dsc.apiServers[index], namespace, err)
			//dsc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_SERVICE, "")
			dsc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_SERVICE, "")
			return
		}

		log.Infof("Delete Service successfully: namespace=%s, apiServer=%s", namespace, dsc.apiServers[index])
	}

	return
}

func (dsc *DeleteServiceController) getOrgNameByOrgId(orgId string) string {
	org, err := myorganization.GetOrganizationById(orgId)
	if err != nil {
		log.Errorf("deleteService Error: error=%s",err)
		dsc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_SERVICE, "")
		return ""
	}
	log.Infof("DeleteServiceController getOrgNameByOrgId successfully: orgName=%s, orgId=%d", org.Name, orgId)
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
			log.Errorf("DeleteMysqlNodePort Error: nodeport=%d, dcId=%d, svcName=%s, error=%s", np.Port, np.DcId, np.SvcName, err)
			dsc.Ye = myerror.NewYceError(myerror.EYCE_DELETE_NODEPORT, "")
			return
		}

		log.Infof("DeleteMysqlNodePort Successfully: nodeport=%d, dcId=%d, svcName=%s", np.Port, np.DcId, np.SvcName)
	}

	return
}

func (dsc DeleteServiceController) Delete() {
	sessionIdFromClient := dsc.RequestHeader("Authorization")
	orgId := dsc.Param("orgId")
	dcId := dsc.Param("dcId")
	userId := dsc.Param("userId")
	svcName := dsc.Param("svcName")
	log.Debugf("DeleteServiceController Params: sessionId=%s, orgId=%s, dcId=%s, userId=%s, svcName=%s", sessionIdFromClient, orgId, dcId, svcName)


	nodePort := new(service.NodePortType)
	err := dsc.ReadJSON(nodePort)
	if err != nil {
		log.Errorf("DeleteServiceController ReadJSON Error: error=%s", err)
		dsc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if dsc.CheckError() {
		return
	}

	log.Debugf("DeleteServiceController ReadJSON: nodePort=%d", nodePort.NodePort)

	// Validate OrgId error
	dsc.ValidateSession(sessionIdFromClient, orgId)
	if dsc.CheckError() {
		return
	}

	// Get DcIdList
	dcIdList := make([]int32, 0)
	datacenterId, _ := strconv.Atoi(dcId)
	dcIdList = append(dcIdList, int32(datacenterId))
	log.Debugf("DeleteServiceController Params: len(dcIdList)=%d", len(dcIdList))

	dsc.getApiServerList(dcIdList)
	if dsc.CheckError() {
		return
	}

	// Get K8sClient
	dsc.createK8sClients()
	if dsc.CheckError() {
		return
	}

	// Publish server to every datacenter
	orgName := dsc.getOrgNameByOrgId(orgId)
	if dsc.CheckError() {
		return
	}

	dsc.deleteService(orgName, svcName)
	if dsc.CheckError() {
		return
	}

	// Update NodePort Status to MySQL nodeport table
	op, _ := strconv.Atoi(userId)
	dsc.deleteMysqlNodePort(dcIdList, nodePort.NodePort, int32(op))
	if dsc.CheckError() {
		return
	}

	dsc.WriteOk("")
	log.Infoln("DeleteServiceController over!")
	return
}
