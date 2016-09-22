package endpoint

import (
	myerror "app/backend/common/yce/error"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"app/backend/model/yce/endpoint"
	"k8s.io/kubernetes/pkg/api"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type CreateEndpointsController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
}


// why need return value?
// Publish k8s.Service to every datacenter which in dcIdList
func (cec *CreateEndpointsController) createEndpoints(namespace string, endpoints *api.Endpoints) {
	// Foreach every K8sClient to create service
	for index, cli := range cec.k8sClients {
		_, err := cli.Endpoints(namespace).Create(endpoints)
		if err != nil {
			log.Errorf("createEndpoints Error: apiServer=%s, namespace=%s, error=%s", cec.apiServers[index], namespace, err)
			cec.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_ENDPOINTS, "")
			return
		}

		log.Infof("Create Endpoints successfully: name=%s, namespace=%s, apiServer=%s", endpoints.Name, namespace, cec.apiServers[index])
	}

	log.Infof("CreateEndpointsController createEndpoints success")
	return
}


func (cec CreateEndpointsController) Post() {
	sessionIdFromClient := cec.RequestHeader("Authorization")
	orgId := cec.Param("orgId")

	log.Debugf("CreateEndpointsController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)


	// Validate OrgId error
	cec.ValidateSession(sessionIdFromClient, orgId)
	if cec.CheckError() {
		return
	}

	// Parse data: service.
	ce := new(endpoint.CreateEndpoints)
	err := cec.ReadJSON(ce)
	if err != nil {
		log.Debugf("CreateEndpointsController ReadJSON Error: error=%s", err)
		cec.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cec.CheckError() {
		return
	}

	// Get DcIdList
	if len(ce.DcIdList) == 0 {
		log.Errorln("Empty DcIdList!")
		cec.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	if cec.CheckError() {
		return
	}

	// Get ApiServer List
	cec.apiServers, cec.Ye = yceutils.GetApiServerList(ce.DcIdList)
	if cec.CheckError() {
		return
	}

	// Create K8sClient List
	cec.k8sClients, cec.Ye = yceutils.CreateK8sClientList(cec.apiServers)
	if cec.CheckError() {
		return
	}

	// Publish server to every datacenter
	orgName := ce.OrgName
	cec.createEndpoints(orgName, &ce.Endpoints)
	if cec.CheckError() {
		return
	}

	cec.WriteOk("")
	log.Infoln("CreateEndpointsController over!")
	return
}