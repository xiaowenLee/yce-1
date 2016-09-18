package deploy

import (
	"app/backend/common/util/Placeholder"
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	mydeployment "app/backend/model/mysql/deployment"
	myoption "app/backend/model/mysql/option"
	myorganization "app/backend/model/mysql/organization"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"app/backend/model/yce/deploy"
	"strconv"
	//"app/backend/common/yce/datacenter"
	yce "app/backend/controller/yce"
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

// getDatacenter by DcId
func (ddc *DeleteDeploymentController) getDatacenterByDcId(dcId int32) *mydatacenter.DataCenter {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("DeleteDeploymentController getDatacenter QueryDataCenterById Error: dcId=%d, error=%s", dcId, err)
		ddc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return nil
	}

	log.Infof("DeleteDeploymentController getDatacenterByDcId successfully: name=%s, id=%d", dc.Name, dc.Id)
	return dc
}

// getApiServer
func (ddc *DeleteDeploymentController) getApiServer() {
	dcId := ddc.getDcId()
	if ddc.Ye != nil {
		return
	}

	datacenter := ddc.getDatacenterByDcId(dcId)
	if ddc.Ye != nil {
		return
	}

	host := datacenter.Host
	port := strconv.Itoa(int(datacenter.Port))
	apiServer := host + ":" + port

	ddc.apiServer = apiServer

	log.Infof("DeleteDeploymentController getApiServer successfully: apiServer=%s, dcId=%d", apiServer, dcId)

}

// create single k8s client by dcId
func (ddc *DeleteDeploymentController) createK8sClient() *client.Client {
	config := &restclient.Config{
		Host: ddc.apiServer,
	}

	c, err := client.New(config)
	if err != nil {
		log.Errorf("DeleteDeploymentController createK8sClient Error: apiServer=%s, error=%s", ddc.apiServer, err)
		ddc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return nil
	}

	log.Debugf("DeleteDeploymentController createK8sClient successfully: apiServer=%s, k8sClient=%p", ddc.apiServer, c)
	return c
}

// get k8sclient
func (ddc *DeleteDeploymentController) getK8sClient() {
	ddc.k8sClient = ddc.createK8sClient()
	if ddc.k8sClient == nil {
		log.Errorf("DeleteDeploymentController createK8sClient Error: apiServer=%s, error=createK8sClient failed", ddc.apiServer)
		ddc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}
	log.Infof("DeleteDeploymentController getK8sClient successfully: k8sClient=%p, apiServer=%s", ddc.k8sClient, ddc.apiServer)
}

// get OrgNameByOrgId
func (ddc *DeleteDeploymentController) getOrgNameByOrgId() string {
	organization := new(myorganization.Organization)

	orgId, _ := strconv.Atoi(ddc.orgId)
	organization.QueryOrganizationById(int32(orgId))
	log.Infof("DeleteDeploymentController getOrgNameByOrgId successfully: orgName=%s, orgId=%d", organization.Name, orgId)
	return organization.Name
}

// getDeploymentByName
func (ddc *DeleteDeploymentController) getDeploymentByName() {

	namespace := ddc.getOrgNameByOrgId()
	var err error
	ddc.deployment, err = ddc.k8sClient.Extensions().Deployments(namespace).Get(ddc.deploymentName)
	if err != nil {
		log.Errorf("DeleteDeploymentController getDeploymentByName Error: apiServer=%s, namespace=%s, err=%s", ddc.apiServer, namespace, err)
		ddc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
		return
	}
	log.Infof("DeleteDeploymentController getDeploymentByName successfully: name=%s, createTime=%s", ddc.deployment.Name, ddc.deployment.CreationTimestamp)
}

// getReplicaSetListByDeployment
func (ddc *DeleteDeploymentController) getReplicaSetListByDeployment() {
	selector, err := unver.LabelSelectorAsSelector(ddc.deployment.Spec.Selector)
	if err != nil {
		log.Errorf("DeleteDeploymentController getReplicaSetListByDeployment Error: name=%s, err=%s", ddc.deployment.Name, err)
		ddc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_REPLICASET, "")
		return
	}

	options := api.ListOptions{LabelSelector: selector}

	rsList, err := ddc.k8sClient.Extensions().ReplicaSets(ddc.deployment.Namespace).List(options)
	if err != nil {
		log.Errorf("DeleteDeploymentController getReplicaSetListByDeployment Error: name=%s, err=%s", ddc.deployment.Name, err)
		ddc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_REPLICASET, "")
		return
	}

	ddc.replicaSets = rsList.Items

	log.Infof("DeleteDeploymentController getReplicaSetListByDeployment successfully: name=%s, len(replicaSet)=%d", ddc.deployment.Name, len(ddc.replicaSets))
}

// deleteReplicaSet
func (ddc *DeleteDeploymentController) deleteReplicaSet() {
	for _, rs := range ddc.replicaSets {
		falseVar := false
		deleteOptions := &api.DeleteOptions{OrphanDependents: &falseVar}

		log.Debugf("DeleteDeploymentController ReplicaSet Name: replicaSetName=%s", rs.Name)
		err := ddc.k8sClient.Extensions().ReplicaSets(ddc.deployment.Namespace).Delete(rs.Name, deleteOptions)
		if err != nil {
			log.Errorf("DeleteDeploymentController deleteReplicaSet Error: name=%s, err=%s", rs.Name, err)
			ddc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_REPLICASET, "")
			return
		}

	}

	log.Infof("DeleteDeploymentController deleteReplicaSet successfully")
}

// getPodsByReplicaSet
func (ddc *DeleteDeploymentController) getPodsByReplicaSet() {
	for _, rs := range ddc.replicaSets {
		selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
		if err != nil {
			log.Errorf("DeleteDeploymentController getPodsByReplicaSet Error: rsName=%s, error=%s", rs.Name, err)
			ddc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
			return
		}

		options := api.ListOptions{LabelSelector: selector}

		podList, err := ddc.k8sClient.Pods(rs.Namespace).List(options)
		if err != nil {
			log.Errorf("DeleteDeploymentController getPodsByReplicaSet Error: rsName=%s, error=%s", rs.Name, err)
			ddc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
			return
		}

		for _, pod := range podList.Items {
			ddc.pods = append(ddc.pods, pod)
		}

		log.Infof("DeleteDeploymentController append pods: len(podList.Items)=%d, len(pods)=%d", len(podList.Items), len(ddc.pods))
	}

	log.Infof("DeleteDeploymentController getPodsByReplicaSet successfully: len(pods)=%d", len(ddc.pods))

}

// delete Pods
func (ddc *DeleteDeploymentController) deletePods() {
	for _, pod := range ddc.pods {
		falseVar := false
		deleteOptions := &api.DeleteOptions{OrphanDependents: &falseVar}

		log.Infof("DeleteDeploymentController deletePods: podName=%s", pod.Name)
		err := ddc.k8sClient.Pods(pod.Namespace).Delete(pod.Name, deleteOptions)

		if err != nil {
			log.Errorf("DeleteDeploymentController deletePods: Error: name=%s, err=%s", pod.Name, err)
			ddc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_POD, "")
			return
		}

	}

	log.Infof("DeleteDeploymentController delete pods successfully")
}

// delete Deployment
func (ddc *DeleteDeploymentController) deleteDeployment() {
	err := ddc.k8sClient.Extensions().Deployments(ddc.deployment.Namespace).Delete(ddc.deployment.Name, nil)
	if err != nil {
		log.Errorf("DeleteDeploymentController deleteDeployment Error: name=%s, err=%s", ddc.deployment.Name, err)
		ddc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_DEPLOYMENT, "")
		return
	}

	log.Infof("DeleteDeploymentController deleteDeployment successfully: name=%s", ddc.deploymentName)
}

func (ddc *DeleteDeploymentController) encodeDcIdList() string {
	/*
		dcIdList := new(deploy.DcIdListType)
		dcIdList.DcIdList = make([]int32, 0)
		dcIdList.DcIdList = append(dcIdList.DcIdList, ddc.params.DcIdList)
	*/
	dcIdList := new(deploy.DcIdListType)
	dcIdList.DcIdList = ddc.params.DcIdList
	//dcIdListJson, _ := json.Marshal(ddc.params.DcIdList)
	dcIdListJson, _ := json.Marshal(dcIdList)

	dcIdListString := string(dcIdListJson)

	log.Infof("DeleteDeploymentController encodeDcIdList successfully: dcIdList=%s", dcIdListString)
	return dcIdListString
}

// create MySQL Deployment of Delete
func (ddc *DeleteDeploymentController) createMysqlDeployment() {

	uph := placeholder.NewPlaceHolder(DELETE_URL)
	actionUrl := uph.Replace("<orgId>", ddc.orgId, "<deploymentName>", ddc.deploymentName)
	actionOp, _ := strconv.Atoi(ddc.params.UserId)
	log.Debugf("DeleteDeploymentController createMySQLDeployment: actionOp=%d, actionUrl=%s", actionOp, actionUrl)

	//dcIdList := strconv.Itoa(int(ddc.params.DcId))
	dcIdList := ddc.encodeDcIdList()
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
	ddc.getDeploymentByName()
	if ddc.Ye != nil {
		ddc.WriteBack()
		return
	}

	// gerReplicaSet List referred to this Deployment
	ddc.getReplicaSetListByDeployment()
	if ddc.Ye != nil {
		ddc.WriteBack()
		return
	}

	// getPods referred to every replicase
	ddc.getPodsByReplicaSet()
	if ddc.Ye != nil {
		ddc.WriteBack()
		return
	}

	// delete Deployment
	ddc.deleteDeployment()
	if ddc.Ye != nil {
		ddc.WriteBack()
		return
	}

	// delete ReplicaSet
	ddc.deleteReplicaSet()
	if ddc.Ye != nil {
		ddc.WriteBack()
		return
	}

	// delete Pods
	ddc.deletePods()
	if ddc.Ye != nil {
		ddc.WriteBack()
		return
	}

	// write delete event to mysql
	ddc.createMysqlDeployment()
	if ddc.Ye != nil {
		ddc.WriteBack()
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
	ddc.getApiServer()
	if ddc.CheckError() {
		return
	}

	// getK8sClient
	ddc.getK8sClient()
	if ddc.CheckError() {
		return
	}

	// deleteDeployment
	ddc.delete()
	if ddc.CheckError() {
		return
	}

	ddc.WriteOk()
	log.Infoln("Delete Deployment over!")
	return

}
