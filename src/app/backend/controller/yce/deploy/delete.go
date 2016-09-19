package deploy

import (
	"app/backend/common/util/Placeholder"
	myerror "app/backend/common/yce/error"
	mydeployment "app/backend/model/mysql/deployment"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type DeleteDeploymentController struct {
	yce.Controller

	// must
	apiServer string
	k8sClient *client.Client

	// url param
	orgId          string
	deploymentName string

	// json param
	params *DeleteDeploymentParam

	// k8s objects
	deployment  *extensions.Deployment
	replicaSets []extensions.ReplicaSet
	pods        []api.Pod
}

// json from client
type DeleteDeploymentParam struct {
	UserId   string  `json:"userId"`
	DcIdList []int32 `json:"dcIdList"`
	Comments string  `json:"comments"`
}


// parse params from json
func (ddc *DeleteDeploymentController) getParams() {
	err := ddc.ReadJSON(ddc.params)
	if err != nil {
		log.Errorf("DeleteDeploymentController getParams Error: error=%s", err)
		ddc.Ye = myerror.NewYceError(myerror.EYCE_DELETE_DEPLOYMENT, "")
		return
	}
	log.Debugf("DeleteDeploymentController getParams successfully: dcId=%d, userId=%s", ddc.params.DcIdList[0], ddc.params.UserId)
}

// get dcId int32
func (ddc *DeleteDeploymentController) getDcId() int32 {
	//dcId, _ := strconv.Itoi(ddc.params.DcId)
	//return int32(dcId)
	//TODO: unsupported multi deletion

	if len(ddc.params.DcIdList) > 0 {
		return ddc.params.DcIdList[0]
	} else {
		log.Errorf("DeleteDeploymentController getDcId Error: len(DcIdList)=%d, err=no value in DcIdList, Index out of range", len(ddc.params.DcIdList))
		ddc.Ye = myerror.NewYceError(myerror.EOOM, "")
		return 0
	}
}

// create MySQL Deployment of Delete
func (ddc *DeleteDeploymentController) createMysqlDeployment() {

	uph := placeholder.NewPlaceHolder(DELETE_URL)
	actionUrl := uph.Replace("<orgId>", ddc.orgId, "<deploymentName>", ddc.deploymentName)
	actionOp, _ := strconv.Atoi(ddc.params.UserId)
	log.Debugf("DeleteDeploymentController createMySQLDeployment: actionOp=%d, actionUrl=%s", actionOp, actionUrl)

	//dcIdList := strconv.Itoa(int(ddc.params.DcId))
	dcIdList, ye := yceutils.EncodeDcIdList(ddc.params.DcIdList)
	if ye != nil {
		ddc.Ye = ye
		return
	}

	orgIdInt, _ := strconv.Atoi(ddc.orgId)
	dp := mydeployment.NewDeployment(ddc.deploymentName, DELETE_VERBE, actionUrl, dcIdList, "", "", ddc.params.Comments, int32(DELETE_TYPE), int32(actionOp), int32(1), int32(orgIdInt))

	err := dp.InsertDeployment()
	if err != nil {
		log.Errorf("DeleteDeploymentController CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcIdList, err)
		ddc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}

	log.Infof("DeleteDeploymentController CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcIdList)
}

// delete all
func (ddc *DeleteDeploymentController) delete() {
	// getDeployment By Name and DcId and namespace
	namespace, ye := yceutils.GetOrgNameByOrgId(ddc.orgId)

	if ye != nil {
		ddc.Ye = ye
		return
	}

	ddc.deployment, ddc.Ye = yceutils.GetDeploymentByNameAndNamespace(ddc.k8sClient, ddc.deploymentName, namespace)

	if ddc.CheckError() {
		return
	}

	// gerReplicaSet List referred to this Deployment
	// ddc.getReplicaSetListByDeployment()
	ddc.replicaSets, ddc.Ye = yceutils.GetReplicaSetsByDeployment(ddc.k8sClient, ddc.deployment)

	if ddc.CheckError() {
		return
	}

	// getPods referred to every replicasets
	ddc.pods, ddc.Ye = yceutils.GetPodsByReplicaSets(ddc.k8sClient, ddc.replicaSets)
	if ddc.CheckError() {
		return
	}

	// delete Deployment
	ddc.Ye = yceutils.DeleteDeployment(ddc.k8sClient, ddc.deployment)
	if ddc.CheckError() {
		return
	}

	// delete ReplicaSet
	ddc.Ye = yceutils.DeleteReplicaSets(ddc.k8sClient, ddc.replicaSets)
	if ddc.CheckError() {
		return
	}

	// delete Pods
	ddc.Ye = yceutils.DeletePods(ddc.k8sClient, ddc.pods)
	if ddc.CheckError() {
		return
	}

	// write delete event to mysql
	ddc.createMysqlDeployment()
	if ddc.CheckError() {
		return
	}

	log.Infof("DeleteDeploymentController delete replicaset and deployment and create mysql deployment successfully")
}

// main
func (ddc DeleteDeploymentController) Post() {

	ddc.params = new(DeleteDeploymentParam)

	sessionIdFromClient := ddc.RequestHeader("Authorization")
	ddc.orgId = ddc.Param("orgId")
	ddc.deploymentName = ddc.Param("deploymentName")

	log.Debugf("DeleteDeploymentController Params: sessionId=%s, orgId=%s, deploymentName=%s", sessionIdFromClient, ddc.orgId, ddc.deploymentName)


	// validate sessionId
	ddc.ValidateSession(sessionIdFromClient, ddc.orgId)
	if ddc.CheckError() {
		return
	}

	//Parse Param
	ddc.getParams()
	if ddc.CheckError() {
		return
	}

	// getApiServer
	dcId := ddc.getDcId()
	ddc.apiServer, ddc.Ye = yceutils.GetApiServerByDcId(dcId)
	if ddc.CheckError() {
		return
	}

	// getK8sClient
	ddc.k8sClient, ddc.Ye = yceutils.CreateK8sClient(ddc.apiServer)
	if ddc.CheckError() {
		return
	}

	// deleteDeployment
	ddc.delete()
	if ddc.CheckError() {
		return
	}

	ddc.WriteOk("")
	log.Infoln("Delete Deployment over!")
	return

}
