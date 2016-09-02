package deploy

import (
	"github.com/kataras/iris"
	"app/backend/common/util/session"
	"app/backend/common/yce/organization"
	"app/backend/model/yce/deploy"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	mylog "app/backend/common/util/log"
	"encoding/json"
	"strconv"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strings"
	deploymentutil "k8s.io/kubernetes/pkg/controller/deployment/util"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type ListDeployController struct {
	*iris.Context
	apiServers []string
	k8sClients []*client.Client
	Ye *myerror.YceError
}

func (ldc *ListDeployController) WriteBack() {
	ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("Create ListDeployController Response Error: controller=%p, code=%d, note=%s", ldc, ldc.Ye.Code, myerror.Errors[ldc.Ye.Code].LogMsg)
	ldc.Write(ldc.Ye.String())
}

// Validate SessionId with OrgId
func (ldc *ListDeployController) validateSessionId(sessionId, orgId string) {
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	// validate error
	if err != nil {
		mylog.Log.Errorf("Create ListDeployController Error: sessionId=%s, orgId=%s, error=%s", sessionId, orgId, err)
		ldc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// invalid sessionId
	if !ok {
		mylog.Log.Errorf("Create ListDeployController Failed: sessionId=%s, orgId=%s", sessionId, orgId)
		ldc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	return
}


// get Datacenters owned by this Organization via OrgId
func (ldc *ListDeployController) getDatacentersByOrgId(ld *deploy.ListDeployment, orgId string) {
	org, err := organization.GetOrganizationById(orgId)
	ld.Organization = org
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		ldc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	dcList, err := organization.GetDataCentersByOrganization(ld.Organization)
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", orgId, err)
		ldc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return
	}

	ld.DcIdList = make([]int32, len(dcList))
	ld.DcName = make([]string, len(dcList))

	for index, dc := range dcList {
		ld.DcIdList[index] = dc.Id
		ld.DcName[index] = dc.Name
	}

	mylog.Log.Infof("CreateServiceController getDatacentersByOrgId: dcList=%s", dcList)
}


// Get ApiServer(k8s cluster host) of this datacenter
func (ldc *ListDeployController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		ldc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateServiceController getApiServerByDcId: apiServer=%s", apiServer)

	return apiServer


}

// Get ApiServer(k8s cluster host) of every datacenter
func (ldc *ListDeployController) getApiServerList(dcIdList []int32) {
	for _, dcId := range dcIdList {
		// Get ApiServer
		apiServer := ldc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("ListDeployController getApiServerList Error")
			return
		}

		ldc.apiServers = append(ldc.apiServers, apiServer)
	}

	mylog.Log.Infof("CreateServiceController getApiServerList successfully")
	return
}

// Create K8s client according to the apiServers
func (ldc *ListDeployController) createK8sClients() {
	// Foreach every ApiServer to create it's k8sClient
	ldc.k8sClients = make([]*client.Client, 0)


	for _, server := range ldc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			mylog.Log.Errorf("CreateK8sClient Error: error=%s", err)
			ldc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		ldc.k8sClients = append(ldc.k8sClients, c)
		// why??
		//ldc.apiServers = append(ldc.apiServers, server)
		mylog.Log.Infof("Append a new client to ldc.K8sClients array: c=%p, apiServer=%s", c, server)
	}

	return
}

// Get PodList by ReplicaSet
func (ldc *ListDeployController)getPodsByReplicaSet(c *client.Client, rs *extensions.ReplicaSet) (*api.PodList) {
	namespace := rs.Namespace
	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		mylog.Log.Errorf("getPodsByReplicaSet Error: error=%s", err)
		ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
		return nil
	}
	options := api.ListOptions{LabelSelector: selector}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		mylog.Log.Errorf("getPodsByReplicaSet Error: error=%s", err)
		ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
		return nil
	}
	mylog.Log.Infof("Get PodList by ReplicaSet successfully: podList=%p", &podList)

	return podList
}

// Get ReplicaSetList by Deployment
func (ldc *ListDeployController)getReplicaSetsByDeployment(c *client.Client, deployment *extensions.Deployment) ([]extensions.ReplicaSet) {

	namespace := deployment.Namespace
	selector, err := unver.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil
		mylog.Log.Errorf("getReplicaSetsByDeployment Error: error=%s", err)
	}
	options := api.ListOptions{LabelSelector: selector}
	rsList, err := c.Extensions().ReplicaSets(namespace).List(options)

	mylog.Log.Infof("Get ReplicaSetList by Deployment successfully: ReplicaSetList=%p", &rsList)
	return rsList.Items
}

func (ldc *ListDeployController) getDeployAndPodList(cli *client.Client, deploymentList *extensions.DeploymentList) (dap []deploy.DeployAndPodList){

	dap = make([]deploy.DeployAndPodList, 0)

	for _, dploy := range deploymentList.Items {

		dp := new(deploy.DeployAndPodList)
		dp.Deploy = new(extensions.Deployment)

		*dp.Deploy = dploy

		rsList := ldc.getReplicaSetsByDeployment(cli, dp.Deploy)
		newRs, err := deploymentutil.FindNewReplicaSet(dp.Deploy, rsList)

		PodList := ldc.getPodsByReplicaSet(cli, newRs)
		if err != nil {
			mylog.Log.Errorf("FindNewReplicaSet Error: error=%s", err)
			ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
			return nil
		}

		dp.PodList = new(api.PodList)
		dp.PodList = PodList

		dap = append(dap, *dp)

	}
	return dap

}

// List all deployments in this namespace
func (ldc *ListDeployController) listDeployments(namespace string, ld *deploy.ListDeployment) (dpString string){
	dpList := make([]deploy.Deployment, 0)

	// Foreach every K8sClient to get DeploymentsList
	for index, cli := range ldc.k8sClients {

		dps, err := cli.Deployments(namespace).List(api.ListOptions{})
		if err != nil {
			mylog.Log.Errorf("listDeployments Error: apiServer=%s, namespace=%s, error=%s", ldc.apiServers[index], namespace, err)
			ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
			return
		}

		//TODO: check consistency


		dp := new(deploy.Deployment)
		dp.DcId = ld.DcIdList[index]
		dp.DcName = ld.DcName[index]
		dp.Deployments = ldc.getDeployAndPodList(cli, dps)

		dpList = append(dpList, *dp)

		mylog.Log.Infof("listDeployments successfully: namespace=%s, apiServer=%s", namespace, ldc.apiServers[index])

	}

	dpJson, err := json.Marshal(dpList)
	dpString = string(dpJson)
	if err != nil {
		mylog.Log.Errorf("listDeployments Error: apiServer=%v, namespace=%s, error=%s", ldc.apiServers, namespace, err)
		ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
		return
	}

	return dpString
}


//GET /api/v1/organizations/{orgId}/users/{userId}/deployments
func (ldc ListDeployController) Get() {
	sessionIdFromClient := ldc.RequestHeader("Authorization")
	orgId := ldc.Param("orgId")

	// validateSessionId
	ldc.validateSessionId(sessionIdFromClient, orgId)
	if ldc.Ye != nil {
		ldc.WriteBack()
		return
	}


	// Get Datacenters by organizs
	ld :=  new(deploy.ListDeployment)
	ldc.getDatacentersByOrgId(ld, orgId)
	if ldc.Ye != nil {
		ldc.WriteBack()
		return
	}


	// Get ApiServers by organizations
	ldc.getApiServerList(ld.DcIdList)
	if ldc.Ye != nil {
		ldc.WriteBack()
		return
	}

	// Get K8sClient
	ldc.createK8sClients()
	if ldc.Ye != nil {
		ldc.WriteBack()
		return
	}

	// List deployments
	orgName := ld.Organization.Name
	dpString := ldc.listDeployments(orgName, ld)
	if ldc.Ye != nil {
		ldc.WriteBack()
		return
	}

	ldc.Ye = myerror.NewYceError(myerror.EOK, dpString)
	ldc.WriteBack()

	mylog.Log.Infoln("ListDeployController over!")

	return
}