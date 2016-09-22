package topology

import (
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	"strconv"
	"strings"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	deployutil "k8s.io/kubernetes/pkg/controller/deployment/util"
)

type DcList struct {
	Items []string `json:"dcList"`
}

type TopologyController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
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

// Get OrgName by orgId
func (tc *TopologyController) getOrgNameByOrgId() {
	org := new(myorganization.Organization)
	err := org.QueryOrganizationById(tc.orgId)
	if err != nil {
		tc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	tc.orgName = org.Name
	log.Infof("TopologyController getOrgNameByOrgId: orgName=%s", tc.orgName)
	return
}

// Get DcIdList by OrgId
func (tc *TopologyController) getDcIdListByOrgId() {
	org := new(myorganization.Organization)
	err := org.QueryOrganizationById(tc.orgId)
	if err != nil {
		log.Errorf("TopologyController QueryOrganizationById error: error=%s", err)
		tc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	tc.dcIdList, tc.Ye = yceutils.DecodeDcIdList(org.DcIdList)
	/*
	dcList := DcList{}
	err = json.Unmarshal([]byte(org.DcList), &dcList)
	if err != nil {
		log.Errorf("TopologyController getDcIdListByOrgId: unmarshal error: error=%s", err)
		tc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	for _, dcId := range dcList.Items {
		id, _ := strconv.Atoi(dcId)
		tc.dcIdList = append(tc.dcIdList, int32(id))
	}

	*/
	// Decode to DcIdList
	log.Infof("TopologyController getDcIdListByOrgId: len(dcIdList)=%d", len(tc.dcIdList))
	return
}

// Get ApiServer by dcId
func (tc *TopologyController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		tc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	log.Infof("TopologyController getApiServerByDcId: apiServer=%s, dcId=%d", apiServer, dcId)
	return apiServer
}

// Get ApiServer List for dcIdList
func (tc *TopologyController) getApiServerList() {
	// Foreach dcIdList
	for _, dcId := range tc.dcIdList {
		// Get ApiServer
		apiServer := tc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			log.Errorf("TopologyController getApiServerList Error")
			return
		}

		tc.apiServers = append(tc.apiServers, apiServer)
	}
	log.Infof("TopologyController getApiServerList: len(apiServer)=%d", len(tc.apiServers))
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
			log.Errorf("createK8sClient Error: err=%s", err)
			tc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		tc.k8sClients = append(tc.k8sClients, c)
		tc.apiServers = append(tc.apiServers, server)
		log.Infof("Append a new client to tc.k8sClients array: c=%p, apiServer=%s", c, server)
	}


	log.Infof("TopologyController createK8sClient: len(k8sClient)=%d", len(tc.k8sClients))
	return
}

/*==========================================================================
 Helper funcs
==========================================================================*/

func getDeploymentsByNamespace(c *client.Client, namespace string) ([]extensions.Deployment, error)  {

	dps, err := c.Extensions().Deployments(namespace).List(api.ListOptions{})
	if err != nil {
		log.Errorf("getDeploymentsByNamespace Error: err=%s\n", err)
		return nil, err
	}

	return dps.Items, nil
}

func getReplicaSetsByDeployment(c *client.Client, deployment *extensions.Deployment) ([]extensions.ReplicaSet, error) {

	namespace := deployment.Namespace
	selector, err := unver.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		log.Errorf("getReplicaSetsByDeployment Error: err=%s\n", err)
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}
	rsList, err := c.Extensions().ReplicaSets(namespace).List(options)

	log.Infof("getReplicaSetsByDeployment: dp.Name=%s, len(rs.Items)=%d\n", deployment.Name, len(rsList.Items))

	return rsList.Items, nil
}

func getPodsByReplicaSet(c *client.Client, namespace string, rs *extensions.ReplicaSet) ([]api.Pod, error) {
	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		log.Infof("getPodsByReplicaSet Error: err=%s\n", err)
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		log.Errorf("getPodsByReplicaSet Error: err=%s\n", err)
		return nil, err
	}

	log.Infof("getPodsByReplicaSet: rs.Name=%s, len(rs.Items)=%d\n", rs.Name, len(podList.Items))

	return podList.Items, nil
}

func getNodeByPod(c *client.Client, pod *api.Pod) (*api.Node, error) {
	name := pod.Spec.NodeName
	node, err := c.Nodes().Get(name)
	if err != nil {
		log.Infof("getNodeByPod Error: err=%s\n", err)
		return nil, err
	}
	return node, nil
}

func getServicesByNamespace(c *client.Client, namespace string) ([]api.Service, error) {
	svcs, err := c.Services(namespace).List(api.ListOptions{})
	if err != nil {
		log.Infof("getServicesByNamespace Error: err=%s\n", err)
		return nil, err
	}

	return svcs.Items, nil
}

func getPodByService(c *client.Client, namespace string, svc *api.Service) ([]api.Pod, error) {
	selector := new(unver.LabelSelector)
	selector.MatchLabels = svc.Spec.Selector

	s, err := unver.LabelSelectorAsSelector(selector)
	if err != nil {
		log.Infof("getPodByService Error: err=%s\n", err)
		return nil, err
	}

	options := api.ListOptions{LabelSelector: s}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		log.Fatalf("getPodByService Error: err=%s\n", err)
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
			log.Errorf("getTopology Error: err=%s\n", err)
			tc.Ye = myerror.NewYceError(myerror.EKUBE_GET_RS_BY_DEPLOYMENT, "")
			return false
		}

		// For all replicasets
		// for _, rs := range rsList {
		rs, err := deployutil.FindNewReplicaSet(&dp, rsList)
		if err != nil {
			log.Errorf("FindNewReplicaSet Error: err=%s", err)
			tc.Ye = myerror.NewYceError(myerror.EKUBE_FIND_NEW_REPLICASET, "")
			return false
		}

		// Topology.Items
		uid := string(rs.UID)
		tc.topology.Items[uid] = ReplicaSetType{
			Kind: "ReplicaSet",
			ApiVersion: "v1beta2",
			ReplicaSet: *rs,
		}

		podList, err := getPodsByReplicaSet(c, namespace, rs)
		if err != nil {
			log.Errorf("getPodsByReplicaSet Error", err)
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
		// }
		// For all replicasets
	}

	// Get Services.List
	svcList, err := getServicesByNamespace(c, namespace)
	if err != nil {
		log.Errorf("getTopology Error: client=%p, namespace=%s, err=%s\n", c, namespace, err)
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
			log.Fatalf("getTopology Error: client=%p, namespace=%s, err=%s\n", c, namespace, err)
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

	data, err := json.Marshal(tc.topology)
	if err != nil {
		log.Errorf("encodeTopology Error: err=%s\n", err)
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
		log.Infof("Get Topology data for every datacenter: apiServer=%s, client=%p\n", tc.apiServers[index], client)
	}


}

// GET /api/v1/organizations/{orgId}/topology
func (tc TopologyController) Get() {
	sessionIdFromClient := tc.RequestHeader("Authorization")
	orgIdStr := tc.Param("orgId")

	log.Debugf("TopologyController Params: sessionId=%s, orgIdStr=%s", sessionIdFromClient, orgIdStr)


	orgId, err := strconv.Atoi(orgIdStr)
	if err != nil {
		tc.Ye = myerror.NewYceError(myerror.EARGS, "")
	}
	if tc.CheckError() {
		return
	}

	tc.orgId = int32(orgId)

	// Validate OrgId error
	tc.ValidateSession(sessionIdFromClient, orgIdStr)
	if tc.CheckError() {
		return
	}

	// Get OrgName by OrgId
	//tc.getOrgNameByOrgId()
	tc.orgName, tc.Ye = yceutils.GetOrgNameByOrgId(orgIdStr)
	if tc.CheckError() {
		return
	}

	// Get DcIdList by OrgId
	//tc.getDcIdListByOrgId()
	tc.dcIdList, tc.Ye = yceutils.GetDcIdListByOrgId(orgIdStr)
	if tc.CheckError() {
		return
	}

	// Get k8s ApiServer by DcIdList
	//tc.getApiServerList()
	tc.apiServers, tc.Ye = yceutils.GetApiServerList(tc.dcIdList)
	if tc.CheckError() {
		return
	}


	// Create k8s clients
	//tc.createK8sClients()
	tc.k8sClients, tc.Ye = yceutils.CreateK8sClientList(tc.apiServers)
	if tc.CheckError() {
		return
	}

	tc.initTopology()

	tc.getTopologyForAllDc()
	if tc.CheckError() {
		return
	}

	str := tc.encodeTopology()
	if tc.CheckError() {
		return
	}

	tc.WriteOk(str)
	log.Infoln("ToplogyController over!")

	return
}


