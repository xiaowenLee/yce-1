package service

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mynodeport "app/backend/model/mysql/nodeport"
	"app/backend/model/yce/service"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
)

type CreateServiceController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string

	params *service.CreateService

	userId string
}


// Publish k8s.Service to every datacenter which in dcIdList
func (csc *CreateServiceController) createService(namespace string, service *api.Service) {
	// Foreach every K8sClient to create service
	for index, cli := range csc.k8sClients {
		_, err := cli.Services(namespace).Create(service)
		if err != nil {
			log.Errorf("createService Error: apiServer=%s, namespace=%s, error=%s", csc.apiServers[index], namespace, err)
			csc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_SERVICE, "")
			return
		}

		log.Infof("Create Service successfully: namespace=%s, apiServer=%s", namespace, csc.apiServers[index])
	}

	log.Infof("CreateServiceController createService success")
	return
}

// create NodePort(mysql) and insert it into db
func (csc *CreateServiceController) createMysqlNodePort(success bool, nodePort int32, dcIdList []int32, svcName string, op int32) {
	for _, dcId := range dcIdList {
		np := mynodeport.NewNodePort(nodePort, dcId, svcName, op)
		err := np.InsertNodePort(op)
		if err != nil {
			log.Errorf("createMysqlNodePort Error: nodeport=%d, dcId=%d, svcName=%s, error=%s", np.Port, np.DcId, np.SvcName, err)
			csc.Ye = myerror.NewYceError(myerror.EYCE_NODEPORT_EXIST, "")
			return
		}

		log.Infof("createMysqlNodePort Successfully: nodeport=%d, dcId=%d, svcName=%s", np.Port, np.DcId, np.SvcName)
	}

	return
}

func (csc *CreateServiceController) getParams() {
	err := csc.ReadJSON(csc.params)
	if err != nil {
		log.Errorf("CreateServiceController ReadJSON Error: error=%s", err)
		csc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}
}

func (csc *CreateServiceController) addNodePort() {
	op, _ := strconv.Atoi(csc.userId)
	cs := csc.params
	for _, port := range cs.Service.Spec.Ports {
		hasNodePort := mynodeport.PORT_START <= port.NodePort && port.NodePort <= mynodeport.PORT_LIMIT
		if hasNodePort {
			csc.createMysqlNodePort(hasNodePort, port.NodePort, cs.DcIdList, cs.Service.ObjectMeta.Name, int32(op))
			if csc.CheckError() {
				return
			}
		}
		log.Infof("CreateServiceController addNodePort : nodeport=%d", port.NodePort)
	}

}

// Post /api/v1/organizations/{:orgId}/user/{:userId}/services/new
func (csc CreateServiceController) Post() {
	sessionIdFromClient := csc.RequestHeader("Authorization")
	orgId := csc.Param("orgId")
	csc.userId = csc.Param("userId")
	log.Debugf("CreateServiceController Params: sessionId=%s, orgId=%s, userId=%s", sessionIdFromClient, orgId, csc.userId)

	csc.params = new(service.CreateService)

	// Validate OrgId error
	csc.ValidateSession(sessionIdFromClient, orgId)
	if csc.CheckError() {
		return
	}

	// get Params from client
	csc.getParams()
	if csc.CheckError() {
		return
	}

	// Get apiServers
	csc.apiServers, csc.Ye = yceutils.GetApiServerList(csc.params.DcIdList)
	if csc.CheckError() {
		return
	}

	// Get K8sClient
	csc.k8sClients, csc.Ye = yceutils.CreateK8sClientList(csc.apiServers)
	if csc.CheckError() {
		return
	}

	// Publish server to every datacenter
	csc.createService(csc.params.OrgName, &csc.params.Service)
	if csc.CheckError() {
		return
	}

	// Add NodePort to MySQL nodeport table
	csc.addNodePort()
	if csc.CheckError() {
		return
	}

	log.Infoln("CreateServiceController over!")
	csc.WriteOk("")
	return
}
