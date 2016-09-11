package topology

import (
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	"app/backend/common/util/session"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	"github.com/kataras/iris"
	"strconv"
	"strings"
)

type DcList struct {
	Items []string `json:"dcList"`
}

type TopologyController struct {
	*iris.Context
	k8sClients []*client.Client
	apiServers []string
	Ye *myerror.YceError
	orgName string
	orgId int32
	topology *Topology
	dcIdList []int32
}

/*==========================================================================
 Definations
==========================================================================*/
type PodType struct {
	api.Pod
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

type ServiceType struct {
	api.Service
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

type ReplicaSetType struct {
	extensions.ReplicaSet
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

type NodeType struct {
	api.Node
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

type ItemType map[string]interface{}

type RelationsType struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Topology struct {
	Items     ItemType        `json:"items"`
	Relations []RelationsType `json:"relations"`
}


func (tc *TopologyController) WriteBack() {
	tc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateDeployController Response YceError: controller=%p, code=%d, note=%s", tc, tc.Ye.Code, myerror.Errors[tc.Ye.Code].LogMsg)
	tc.Write(tc.Ye.String())
}

// Validate Session
func (tc *TopologyController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		tc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		tc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	return
}

// Get OrgName by orgId
func (tc *TopologyController) getOrgNameByOrgId() {
	org := new(myorganization.Organization)
	err := org.QueryOrganizationById(tc.orgId)
	if err != nil {
		tc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	tc.orgName = org.Name
	return
}

// Get DcIdList by OrgId
func (tc *TopologyController) getDcIdListByOrgId() {
	org := new(myorganization.Organization)
	err := org.QueryOrganizationById(tc.orgId)
	if err != nil {
		tc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	dcList := DcList{}
	err = json.Unmarshal([]byte(org.DcList), &dcList)
	if err != nil {
		tc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	for _, dcId := range dcList.Items {
		id, _ := strconv.Atoi(dcId)
		tc.dcIdList = append(tc.dcIdList, int32(id))
	}

	// Decode to DcIdList
	return
}

// Get ApiServer by dcId
func (tc *TopologyController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		tc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateDeployController getApiServerByDcId: apiServer=%s, dcId=%d", apiServer, dcId)
	return apiServer
}

// Get ApiServer List for dcIdList
func (tc *TopologyController) getApiServerList() {
	// Foreach dcIdList
	for _, dcId := range tc.dcIdList {
		// Get ApiServer
		apiServer := tc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("CreateDeployController getApiServerList Error")
			return
		}

		tc.apiServers = append(tc.apiServers, apiServer)
	}
	return
}

// Create k8sClients for every ApiServer
func (tc *TopologyController) createK8sClients() {

	// Foreach every ApiServer to create it's k8sClient
	tc.k8sClients = make([]*client.Client, 0)

	for _, server := range tc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			mylog.Log.Errorf("createK8sClient Error: err=%s", err)
			tc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		tc.k8sClients = append(tc.k8sClients, c)
		tc.apiServers = append(tc.apiServers, server)
		mylog.Log.Infof("Append a new client to tc.k8sClients array: c=%p, apiServer=%s", c, server)
	}

	return
}

/*==========================================================================
 Helper funcs
==========================================================================*/

func getDeploymentsByNamespace(c *client.Client, namespace string) ([]extensions.Deployment, error)  {

	dps, err := c.Extensions().Deployments(namespace).List(api.ListOptions{})
	if err != nil {
		mylog.Log.Errorf("getDeploymentsByNamespace Error: err=%s\n", err)
		return nil, err
	}

	return dps.Items, nil
}

func getReplicaSetsByDeployment(c *client.Client, deployment *extensions.Deployment) ([]extensions.ReplicaSet, error) {

	namespace := deployment.Namespace
	selector, err := unver.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		mylog.Log.Errorf("getReplicaSetsByDeployment Error: err=%s\n", err)
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}
	rsList, err := c.Extensions().ReplicaSets(namespace).List(options)

	mylog.Log.Infof("getReplicaSetsByDeployment: dp.Name=%s, len(rs.Items)=%d\n", deployment.Name, len(rsList.Items))

	return rsList.Items, nil
}

func getPodsByReplicaSet(c *client.Client, namespace string, rs *extensions.ReplicaSet) ([]api.Pod, error) {
	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		mylog.Log.Infof("getPodsByReplicaSet Error: err=%s\n", err)
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		mylog.Log.Errorf("getPodsByReplicaSet Error: err=%s\n", err)
		return nil, err
	}

	mylog.Log.Infof("getPodsByReplicaSet: rs.Name=%s, len(rs.Items)=%d\n", rs.Name, len(podList.Items))

	return podList.Items, nil
}

func getNodeByPod(c *client.Client, pod *api.Pod) (*api.Node, error) {
	name := pod.Spec.NodeName
	node, err := c.Nodes().Get(name)
	if err != nil {
		mylog.Log.Infof("getNodeByPod Error: err=%s\n", err)
		return nil, err
	}
	return node, nil
}

func getServicesByNamespace(c *client.Client, namespace string) ([]api.Service, error) {
	svcs, err := c.Services(namespace).List(api.ListOptions{})
	if err != nil {
		mylog.Log.Infof("getServicesByNamespace Error: err=%s\n", err)
		return nil, err
	}

	return svcs.Items, nil
}

func getPodByService(c *client.Client, namespace string, svc *api.Service) ([]api.Pod, error) {
	selector := new(unver.LabelSelector)
	selector.MatchLabels = svc.Spec.Selector

	s, err := unver.LabelSelectorAsSelector(selector)
	if err != nil {
		mylog.Log.Infof("getPodByService Error: err=%s\n", err)
		return nil, err
	}

	options := api.ListOptions{LabelSelector: s}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		mylog.Log.Fatalf("getPodByService Error: err=%s\n", err)
		return nil, err
	}

	return podList.Items, nil
}

/*==========================================================================
 Topology
==========================================================================*/
/*

begin:
	ops --> Deployments.List
	Foreach deployment in Deployments.List
		rs := findNewReplicaSet()
		rs --> Select Pods.List
		Foreach pod in Pods.List(){
			Pod.Name --> Node: pod <---> node
			rs <---> pod
		}
	ops --> Services.List
	Foreach service in Services.List
		service --> Select Pods.List
		service <---> pod
:end
*/
func (tc *TopologyController) getTopology(c *client.Client, namespace string) bool {
	// Get Deployments.List
	dpList, err := getDeploymentsByNamespace(c, namespace)
	if err != nil {
		return false
	}

	// Foreach Deployments.List
	for _, dp := range dpList {
		rsList, err := getReplicaSetsByDeployment(c, &dp)
		if err != nil {
			mylog.Log.Errorf("getTopology Error: err=%s\n", err)
			tc.Ye = myerror.NewYceError(myerror.EKUBE_GET_RS_BY_DEPLOYMENT, "")
			return false
		}

		for _, rs := range rsList {
			// Topology.Items
			uid := string(rs.UID)
			tc.topology.Items[uid] = ReplicaSetType{
				Kind: "ReplicaSet",
				ApiVersion: "v1beta2",
				ReplicaSet: rs,
			}

			podList, err := getPodsByReplicaSet(c, namespace, &rs)
			if err != nil {
				mylog.Log.Errorf("getPodsByReplicaSet Error", err)
				tc.Ye = myerror.NewYceError(myerror.EKUBE_GET_PODS_BY_RS, "")
				return false
			}
			for _, pod := range podList {
				uid = string(pod.UID)
				tc.topology.Items[uid] = PodType{
					Kind: "Pod",
					ApiVersion: "v1beta2",
					Pod: pod,
				}

				relation := RelationsType {
					Source: string(rs.UID),
					Target: string(pod.UID),
				}

				tc.topology.Relations = append(tc.topology.Relations, relation)

				node, err := getNodeByPod(c, &pod)
				if err != nil {
					tc.Ye = myerror.NewYceError(myerror.EKUBE_GET_NODE_BY_POD, "")
					return false
				}

				uid = string(node.UID)
				tc.topology.Items[uid] = NodeType{
					Kind: "Node",
					ApiVersion: "v1beata2",
					Node: *node,
				}

				relation = RelationsType {
					Source: string(node.UID),
					Target: string(pod.UID),
				}
				tc.topology.Relations = append(tc.topology.Relations, relation)
			}
		}
	}

	// Get Services.List
	svcList, err := getServicesByNamespace(c, namespace)
	if err != nil {
		mylog.Log.Errorf("getTopology Error: client=%p, namespace=%s, err=%s\n", c, namespace, err)
		tc.Ye = myerror.NewYceError(myerror.EKUBE_GET_SERVICES_BY_NAMESPACE, "")
		return false
	}

	for _, svc := range svcList {
		uid := string(svc.UID)
		tc.topology.Items[uid] = ServiceType{
			Kind: "Service",
			ApiVersion: "v1beta1",
			Service: svc,
		}

		podList, err := getPodByService(c, namespace, &svc)
		if err != nil {
			mylog.Log.Fatalf("getTopology Error: client=%p, namespace=%s, err=%s\n", c, namespace, err)
			tc.Ye = myerror.NewYceError(myerror.EKUBE_GET_PODS_BY_SERVICE, "")
			return false
		}

		for _, pod := range podList {
			uid = string(pod.UID)
			if _, ok := tc.topology.Items[uid]; ok {
				tc.topology.Items[uid] = PodType{
					Kind: "Pod",
					ApiVersion: "v1beta1",
					Pod: pod,
				}
			}

			relation := RelationsType {
				Source: string(svc.UID),
				Target: string(pod.UID),
			}

			tc.topology.Relations = append(tc.topology.Relations, relation)
		}
	}

	return true
}

func (tc *TopologyController) encodeTopology() string {

	data, err := json.MarshalIndent(tc.topology, "", "\t")
	if err != nil {
		mylog.Log.Errorf("encodeTopology Error: err=%s\n", err)
		tc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}
	return string(data)
}

func (tc *TopologyController) initTopology() {
	tc.topology = new(Topology)
	tc.topology.Items = make(ItemType)
	tc.topology.Relations = make([]RelationsType, 0)
}

func (tc *TopologyController) getTopologyForAllDc() {
	for index, client := range tc.k8sClients {
		tc.getTopology(client, tc.orgName)
		mylog.Log.Infof("Get Topology data for every datacenter: apiServer=%s, client=%p\n", tc.apiServers[index], client)
	}
}

// GET /api/v1/organizations/{orgId}/topology
func (tc TopologyController) Get() {
	sessionIdFromClient := tc.RequestHeader("Authorization")
	orgIdStr := tc.Param("orgId")
	orgId, err := strconv.Atoi(orgIdStr)
	if err != nil {
		tc.Ye = myerror.NewYceError(myerror.EARGS, "")
		tc.WriteBack()
		return
	}
	tc.orgId = int32(orgId)

	// Validate OrgId error
	tc.validateSession(sessionIdFromClient, orgIdStr)
	if tc.Ye != nil {
		tc.WriteBack()
		return
	}

	// Get OrgName by OrgId
	tc.getOrgNameByOrgId()
	if tc.Ye != nil {
		tc.WriteBack()
		return
	}

	// Get DcIdList by OrgId
	tc.getDcIdListByOrgId()
	if tc.Ye != nil {
		tc.WriteBack()
		return
	}

	// Get k8s ApiServer by DcIdList
	tc.getApiServerList()
	if tc.Ye != nil {
		tc.WriteBack()
		return
	}


	// Create k8s clients
	tc.createK8sClients()
	if tc.Ye != nil {
		tc.WriteBack()
		return
	}

	tc.initTopology()

	tc.getTopologyForAllDc()
	if tc.Ye != nil {
		tc.WriteBack()
		return
	}

	str := tc.encodeTopology()
	if tc.Ye != nil {
		tc.WriteBack()
		return
	}

	tc.Ye = myerror.NewYceError(myerror.EOK, str)
	tc.WriteBack()

	mylog.Log.Infoln("ToplogyController over!")
	return
}


