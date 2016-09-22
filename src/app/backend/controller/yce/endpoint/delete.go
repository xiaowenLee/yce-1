package endpoint

import (
	myerror "app/backend/common/yce/error"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	//"strconv"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type DeleteEndpointsController struct {
	yce.Controller
	k8sClient *client.Client
	apiServer string

	params DeleteEndpointsParam
}


type DeleteEndpointsParam struct {
	DcId int32 `json:dcId`
}

// Publish k8s.Endpoint to every datacenter which in dcIdList
func (dec *DeleteEndpointsController) deleteEndpoints(namespace, epName string) {
	// Foreach every K8sClient to create service

	//for index, cli := range dec.k8sClients {
	err := dec.k8sClient.Endpoints(namespace).Delete(epName)
	if err != nil {
		log.Errorf("deleteEndpoint Error: apiServer=%s, namespace=%s, error=%s", dec.apiServer, namespace, err)
		//dec.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_SERVICE, "")
		dec.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_ENDPOINT, "")
		return
	}

	log.Infof("DeleteEndpointsController Delete Endpoints successfully: namespace=%s, apiServer=%s", namespace, dec.apiServer)
	return
}

// get Params
func (dec *DeleteEndpointsController) getParams() {
	err := dec.ReadJSON(dec.params)
	if err != nil {
		log.Errorf("DeleteEndpointsController getParams Error: error=%s", err)
		dec.Ye = myerror.NewYceError(myerror.EYCE_DELETE_ENDPOINTS, "")
		return
	}

	log.Errorf("DeleteEndpointsController getParams success: dcId=%d", dec.params.DcId)
}

//func (dec DeleteEndpointsController) Delete() {
//POST /api/v1/organizations/{:orgId}/endpoints/{:epName}
func (dec DeleteEndpointsController) Post() {
	sessionIdFromClient := dec.RequestHeader("Authorization")
	orgId := dec.Param("orgId")
	//dcId := dec.Param("dcId")
	epName := dec.Param("epName")

	dec.getParams()
	if dec.CheckError() {
		return
	}

	log.Debugf("DeleteEndpontsController Params: sessionId=%s, orgId=%s, dec.params.dcId=%d, epName=%s", sessionIdFromClient, orgId, dec.params.DcId, epName)

	// Validate OrgId error
	dec.ValidateSession(sessionIdFromClient, orgId)
	if dec.CheckError() {
		return
	}

	// Get ApiServer List
	dec.apiServer, dec.Ye = yceutils.GetApiServerByDcId(dec.params.DcId)
	if dec.CheckError() {
		return
	}

	// Create K8sClient List
	dec.k8sClient, dec.Ye = yceutils.CreateK8sClient(dec.apiServer)
	if dec.CheckError() {
		return
	}

	// Publish server to every datacenter
	orgName, ye := yceutils.GetOrgNameByOrgId(orgId)
	if ye != nil {
		dec.Ye = ye
	}
	if dec.CheckError() {
		return
	}

	dec.deleteEndpoints(orgName, epName)
	if dec.CheckError() {
		return
	}

	dec.WriteOk("")
	log.Infoln("DeleteEndpointsController over!")
	return
}
