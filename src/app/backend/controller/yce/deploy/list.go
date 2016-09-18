package deploy

import (
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	deploymentutil "k8s.io/kubernetes/pkg/controller/deployment/util"
	"app/backend/common/util/mysql"
	"strconv"
	"strings"
	yce "app/backend/controller/yce"
)

type ListDeploymentController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client
}

// get Datacenters owned by this Organization via OrgId
func (ldc *ListDeploymentController) getDatacentersByOrgId(ld *deploy.ListDeployment, orgId string) {
	org, err := organization.GetOrganizationById(orgId)
	ld.Organization = org
	if err != nil {
		log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		ldc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	dcList, err := organization.GetDataCentersByOrganization(ld.Organization)
	if err != nil {
		log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		ldc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return
	}

	ld.DcIdList = make([]int32, 0)
	ld.DcName = make([]string, 0)

	for _, dc := range dcList {
		ld.DcIdList = append(ld.DcIdList, dc.Id)
		ld.DcName = append(ld.DcName, dc.Name)
	}

	log.Infof("CreateServiceController getDatacentersByOrgId: len(DcIdList)=%d, len(DcName)=%d", len(ld.DcIdList), len(ld.DcName))
}

// Get ApiServer(k8s cluster host) of this datacenter
func (ldc *ListDeploymentController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		ldc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("CreateServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer

}

// Get ApiServer(k8s cluster host) of every datacenter
func (ldc *ListDeploymentController) getApiServerList(dcIdList []int32) {
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := ldc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			log.Errorf("ListDeploymentController getApiServerList Error")
			return
		}

		ldc.apiServers = append(ldc.apiServers, apiServer)
	}

	log.Infof("CreateServiceController getApiServerList success: len(apiServers)=%d", len(ldc.apiServers))
	return
}

// Create K8s client according to the apiServers
func (ldc *ListDeploymentController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	ldc.k8sClients = make([]*client.Client, 0)

	for _, server := range ldc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			log.Errorf("CreateK8sClient Error: error=%s", err)
			ldc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		ldc.k8sClients = append(ldc.k8sClients, c)
		// why??
		//ldc.apiServers = append(ldc.apiServers, server)
		log.Infof("ListDeploymentController create K8sClients successfully: client=%p, apiServer=%s", c, server)
	}

	return
}

// Get PodList by ReplicaSet
func (ldc *ListDeploymentController) getPodsByReplicaSet(c *client.Client, rs *extensions.ReplicaSet) *api.PodList {
	namespace := rs.Namespace
	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		log.Errorf("getPodsByReplicaSet Error: error=%s", err)
		ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
		return nil
	}
	options := api.ListOptions{LabelSelector: selector}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		log.Errorf("getPodsByReplicaSet Error: error=%s", err)
		ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
		return nil
	}
	log.Infof("Get PodList by ReplicaSet successfully: podList=%p", &podList)

	return podList
}

// Get ReplicaSetList by Deployment
func (ldc *ListDeploymentController) getReplicaSetsByDeployment(c *client.Client, deployment *extensions.Deployment) []extensions.ReplicaSet {

	namespace := deployment.Namespace
	selector, err := unver.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil
		log.Errorf("getReplicaSetsByDeployment Error: error=%s", err)
	}
	options := api.ListOptions{LabelSelector: selector}
	rsList, err := c.Extensions().ReplicaSets(namespace).List(options)

	log.Infof("Get ReplicaSetList by Deployment successfully: ReplicaSetList=%p", &rsList)
	return rsList.Items
}

// Query UserName by UserId
func (ldc *ListDeploymentController) queryUserNameByUserId(userId int32) (name string) {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(SELECT_USER)
	if err != nil {
		log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		ldc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId).Scan(&name)
	if err != nil {
		log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		ldc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	log.Infof("queryUserNameByUserId successfully")
	return name
}

// Get DeployAndPodList Pair by deploymentList
func (ldc *ListDeploymentController) getDeployAndPodList(userId int32, cli *client.Client, deploymentList *extensions.DeploymentList) (dap []deploy.DeployAndPodList) {

	dap = make([]deploy.DeployAndPodList, 0)

	for _, deployment := range deploymentList.Items {

		dp := new(deploy.DeployAndPodList)

		dp.UserName = ldc.queryUserNameByUserId(userId)



		dp.Deploy = new(extensions.Deployment)

		*dp.Deploy = deployment

		rsList := ldc.getReplicaSetsByDeployment(cli, dp.Deploy)
		newRs, err := deploymentutil.FindNewReplicaSet(dp.Deploy, rsList)

		PodList := ldc.getPodsByReplicaSet(cli, newRs)
		if err != nil {
			log.Errorf("FindNewReplicaSet Error: error=%s", err)
			ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
			return nil
		}

		dp.PodList = new(api.PodList)
		dp.PodList = PodList

		dap = append(dap, *dp)

	}
	log.Infof("ListDeploymentController getDeployAndPodList successfully")
	return dap

}

// List all deployments in this namespace
func (ldc *ListDeploymentController) listDeployments(userId int32, namespace string, ld *deploy.ListDeployment) (dpString string) {
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
		deployment.DcId = ld.DcIdList[index]
		deployment.DcName = ld.DcName[index]
		deployment.Deployments = ldc.getDeployAndPodList(userId, cli, deploymentList)

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

	log.Infoln("ListDeploymentController listDeployments success: len(Deployment)=%d", len(dpList))
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
	ld := new(deploy.ListDeployment)
	ldc.getDatacentersByOrgId(ld, orgId)
	if ldc.CheckError() {
		return
	}

	// Get ApiServers by organizations
	ldc.getApiServerList(ld.DcIdList)
	if ldc.CheckError() {
		return
	}

	// Get K8sClient
	ldc.createK8sClients()
	if ldc.CheckError() {
		return
	}

	// List deployments
	orgName := ld.Organization.Name
	uId, _ := strconv.Atoi(userId)
	dpString := ldc.listDeployments(int32(uId), orgName, ld)
	if ldc.CheckError() {
		return
	}

	ldc.WriteOk(dpString)
	log.Infoln("ListDeploymentController Get over!")

	return
}
