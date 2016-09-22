package deploy

import (
	myerror "app/backend/common/yce/error"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type ListDeploymentController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client
}

// List all deployments in this namespace
func (ldc *ListDeploymentController) listDeployments(userId int32, namespace string, dcList *yceutils.DatacenterList) (dpString string) {
	dpList := make([]deploy.Deployment, 0)

	// Foreach every K8sClient to get DeploymentsList
	for index, cli := range ldc.k8sClients {

		deploymentList, err := cli.Deployments(namespace).List(api.ListOptions{})
		if err != nil {
			log.Errorf("listDeployments Error: apiServer=%s, namespace=%s, error=%s", ldc.apiServers[index], namespace, err)
			ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
			return
		}

		//TODO: check consistency

		deployment := new(deploy.Deployment)
		deployment.DcId = dcList.DcIdList[index]
		deployment.DcName = dcList.DcName[index]
		deployment.Deployments, ldc.Ye = yceutils.GetDeployAndPodList(userId, cli, deploymentList)

		dpList = append(dpList, *deployment)

		log.Infof("listDeployments successfully: namespace=%s, apiServer=%s", namespace, ldc.apiServers[index])

	}

	dpJson, err := json.Marshal(dpList)
	dpString = string(dpJson)
	if err != nil {
		log.Errorf("listDeployments Error: apiServer=%v, namespace=%s, error=%s", ldc.apiServers, namespace, err)
		ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
		return
	}

	log.Infof("ListDeploymentController listDeployments success: len(Deployment)=%d", len(dpList))
	return dpString
}

//GET /api/v1/organizations/{orgId}/users/{userId}/deployments
func (ldc ListDeploymentController) Get() {
	sessionIdFromClient := ldc.RequestHeader("Authorization")
	orgId := ldc.Param("orgId")
	userId := ldc.Param("userId")

	log.Debugf("ListDeploymentController Params: sessionId=%s, orgid=%s, userId=%s", sessionIdFromClient, orgId, userId)


	// ValidateSession
	ldc.ValidateSession(sessionIdFromClient, orgId)
	if ldc.CheckError() {
		return
	}

	// Get Datacenters by organizs
	dcList, ye := yceutils.GetDatacenterListByOrgId(orgId)
	if ye != nil {
		ldc.Ye = ye
	}
	if ldc.CheckError() {
		return
	}

	// Get ApiServers by organizations
	ldc.apiServers, ldc.Ye = yceutils.GetApiServerList(dcList.DcIdList)
	if ldc.CheckError() {
		return
	}

	// Get K8sClient
	ldc.k8sClients, ldc.Ye = yceutils.CreateK8sClientList(ldc.apiServers)
	if ldc.CheckError() {
		return
	}

	// List deployments
	orgName, ye := yceutils.GetOrgNameByOrgId(orgId)
	if ye != nil {
		ldc.Ye = ye
	}
	if ldc.CheckError() {
		return
	}

	// List Deployments and marshal it to json
	uId, _ := strconv.Atoi(userId)
	dpString := ldc.listDeployments(int32(uId), orgName, dcList)
	if ldc.CheckError() {
		return
	}

	ldc.WriteOk(dpString)
	log.Infoln("ListDeploymentController Get over!")

	return
}
