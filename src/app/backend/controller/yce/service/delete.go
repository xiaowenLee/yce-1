package service

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mynodeport "app/backend/model/mysql/nodeport"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

type DeleteServiceController struct {
	yce.Controller
	k8sClient *client.Client
	apiServer string

	params *DeleteServiceParam
}

type DeleteServiceParam struct {
	UserId   int32 `json:"userId"`
	DcId     int32 `json:"dcId"`
	NodePort int32 `json:"nodePort"`
}

// Publish k8s.Service to every datacenter which in dcIdList
func (dsc *DeleteServiceController) deleteService(namespace, svcName string) {
	err := dsc.k8sClient.Services(namespace).Delete(svcName)
	if err != nil {
		log.Errorf("deleteService Error: apiServer=%s, namespace=%s, error=%s", dsc.apiServer, namespace, err)
		dsc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_SERVICE, "")
		return
	}

	log.Infof("Delete Service successfully: namespace=%s, apiServer=%s", namespace, dsc.apiServer)

	return
}

// create NodePort(mysql) and insert it into db
func (dsc *DeleteServiceController) deleteMysqlNodePort() {
	dcId := dsc.params.DcId
	op := dsc.params.UserId
	nodePort := dsc.params.NodePort
	np := &mynodeport.NodePort{
		Port: nodePort,
		DcId: dcId,
	}

	//err := np.DeleteNodePort(op)
	err := np.ReleaseNodePort(op)
	if err != nil {
		log.Errorf("DeleteMysqlNodePort Error: nodeport=%d, dcId=%d, svcName=%s, error=%s", np.Port, np.DcId, np.SvcName, err)
		dsc.Ye = myerror.NewYceError(myerror.EYCE_DELETE_NODEPORT, "")
		return
	}

	log.Infof("DeleteMysqlNodePort Successfully: nodeport=%d, dcId=%d, svcName=%s", np.Port, np.DcId, np.SvcName)

	return
}

func (dsc DeleteServiceController) getParams() {
	err := dsc.ReadJSON(dsc.params)
	if err != nil {
		log.Errorf("DeleteServiceController getParams Error: error=%s", err)
		dsc.Ye = myerror.NewYceError(myerror.EYCE_DELETE_SERVICE, "")
		return
	}
	log.Debugf("DeleteServiceController getParams success: userId=%d, dcId=%d, nodePort=%d", dsc.params.UserId, dsc.params.DcId, dsc.params.NodePort)
}

// Post /api/v1/organizations/{:orgId}/services/{:svcName}
func (dsc DeleteServiceController) Post() {
	sessionIdFromClient := dsc.RequestHeader("Authorization")
	orgId := dsc.Param("orgId")
	svcName := dsc.Param("svcName")
	dsc.params = new(DeleteServiceParam)

	//TODO: frontend change params json, url and method
	dsc.getParams()
	if dsc.CheckError() {
		return
	}

	log.Debugf("DeleteServiceController Params: sessionId=%s, orgId=%s, dcId=%d, userId=%d, nodePort=%d, svcName=%s", sessionIdFromClient, orgId, dsc.params.DcId, dsc.params.UserId, dsc.params.NodePort, svcName)

	// Validate OrgId error
	dsc.ValidateSession(sessionIdFromClient, orgId)
	if dsc.CheckError() {
		return
	}

	// Get ApiServer
	dsc.apiServer, dsc.Ye = yceutils.GetApiServerByDcId(dsc.params.DcId)
	if dsc.CheckError() {
		return
	}

	// Create K8sClient
	dsc.k8sClient, dsc.Ye = yceutils.CreateK8sClient(dsc.apiServer)
	if dsc.CheckError() {
		return
	}

	// Publish server to every datacenter
	var orgName string
	orgName, dsc.Ye = yceutils.GetOrgNameByOrgId(orgId)
	if dsc.CheckError() {
		return
	}

	dsc.deleteService(orgName, svcName)
	if dsc.CheckError() {
		return
	}

	// Update NodePort Status to MySQL nodeport table
	dsc.deleteMysqlNodePort()
	if dsc.CheckError() {
		return
	}

	dsc.WriteOk("")
	log.Infoln("DeleteServiceController over!")
	return
}
