package topology

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mydatacenter "app/backend/model/mysql/datacenter"
	"encoding/json"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	deployutil "k8s.io/kubernetes/pkg/controller/deployment/util"
	"strconv"
	"strings"
)

type DcList struct {
	Items []string `json:"dcList"`
}

type TopologyController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
	orgName    string
	orgId      int32
	topology   *Topology
	dcIdList   []int32
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

func (tc *TopologyController) getAllDcIdList() error {
	dcs, err := mydatacenter.QueryAllDatacenters()

	if err != nil {
		tc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return err
	}

	tc.dcIdList = make([]int32, 0)

	for _, dc := range dcs {
		tc.dcIdList = append(tc.dcIdList, dc.Id)
	}
	return nil
}

func (tc *TopologyController) getTopology(c *client.Client, namespace string) bool {
	// Get Deployments.List
	dpList, ye := yceutils.GetDeploymentByNamespace(c, namespace)

	if ye != nil {
		tc.Ye = ye
	}
	if tc.CheckError() {
		return false
	}

	// Foreach Deployments.List
	for _, dp := range dpList {
		//rsList, err := getReplicaSetsByDeployment(c, &dp)
		rsList, ye := yceutils.GetReplicaSetsByDeployment(c, &dp)
		if ye != nil {
			tc.Ye = ye
		}
		if tc.CheckError() {
			return false
		}

		// For all replicasets
		rs, err := deployutil.FindNewReplicaSet(&dp, rsList)
		if err != nil {
			log.Errorf("FindNewReplicaSet Error: err=%s", err)
			tc.Ye = myerror.NewYceError(myerror.EKUBE_FIND_NEW_REPLICASET, "")
			return false
		}

		// Topology.Items
		uid := string(rs.UID)
		tc.topology.Items[uid] = ReplicaSetType{
			Kind:       "ReplicaSet",
			ApiVersion: "v1beta2",
			ReplicaSet: *rs,
		}

		podList, ye := yceutils.GetPodListByReplicaSet(c, rs)

		if ye != nil {
			tc.Ye = ye
		}
		if tc.CheckError() {
			return false
		}

		for _, pod := range podList.Items {
			uid = string(pod.UID)
			tc.topology.Items[uid] = PodType{
				Kind:       "Pod",
				ApiVersion: "v1beta2",
				Pod:        pod,
			}

			relation := RelationsType{
				Source: string(rs.UID),
				Target: string(pod.UID),
			}

			tc.topology.Relations = append(tc.topology.Relations, relation)

			node, ye := yceutils.GetNodeByPod(c, &pod)
			if ye != nil {
				tc.Ye = ye
			}
			if tc.CheckError() {
				return false
			}

			uid = string(node.UID)
			tc.topology.Items[uid] = NodeType{
				Kind:       "Node",
				ApiVersion: "v1beata2",
				Node:       *node,
			}

			relation = RelationsType{
				Source: string(node.UID),
				Target: string(pod.UID),
			}
			tc.topology.Relations = append(tc.topology.Relations, relation)
		}
		// }
		// For all replicasets
	}

	// Get Services.List
	svcList, ye := yceutils.GetServicesByNamespace(c, namespace)
	if ye != nil {
		tc.Ye = ye
	}
	if tc.CheckError() {
		return false
	}

	for _, svc := range svcList {
		uid := string(svc.UID)
		tc.topology.Items[uid] = ServiceType{
			Kind:       "Service",
			ApiVersion: "v1beta1",
			Service:    svc,
		}

		podList, ye := yceutils.GetPodsByService(c, &svc)
		if ye != nil {
			tc.Ye = ye
		}
		if tc.CheckError() {
			return false
		}
		for _, pod := range podList {
			uid = string(pod.UID)
			if _, ok := tc.topology.Items[uid]; ok {
				tc.topology.Items[uid] = PodType{
					Kind:       "Pod",
					ApiVersion: "v1beta1",
					Pod:        pod,
				}
			}

			relation := RelationsType{
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

func (tc *TopologyController) getTopologyForAllDcAndAllNamespaces() {
	for index, client := range tc.k8sClients {
		namespaces, ye := yceutils.GetAllNamespaces(client)
		if ye != nil {
			tc.Ye = ye
			return
		}

		for _, name := range namespaces {
			tc.getTopology(client, name)
			log.Infof("Get Topology data for every datacenter: apiServer=%s, client=%p, namespace=%s\n", tc.apiServers[index], client, name)
		}
	}
}

// GET /api/v1/organizations/{orgId}/users/{userId}/topology
func (tc TopologyController) Get() {
	sessionIdFromClient := tc.RequestHeader("Authorization")
	orgIdStr := tc.Param("orgId")
	userIdStr := tc.Param("userId")

	log.Debugf("TopologyController Params: sessionId=%s, orgIdStr=%s", sessionIdFromClient, orgIdStr)

	orgId, err := strconv.Atoi(orgIdStr)
	if err != nil {
		tc.Ye = myerror.NewYceError(myerror.EARGS, "")
	}
	if tc.CheckError() {
		return
	}

	// Validate OrgId error
	tc.ValidateSession(sessionIdFromClient, orgIdStr)
	if tc.CheckError() {
		return
	}

	if strings.EqualFold(userIdStr, ADMIN_ID) {
		// Get all dcId for all datacenters
		tc.getAllDcIdList()
		if tc.CheckError() {
			return
		}

	} else {

		tc.orgId = int32(orgId)
		// Get OrgName by OrgId
		tc.orgName, tc.Ye = yceutils.GetOrgNameByOrgId(orgIdStr)
		if tc.CheckError() {
			return
		}

		// Get DcIdList by OrgId
		tc.dcIdList, tc.Ye = yceutils.GetDcIdListByOrgId(orgIdStr)
		if tc.CheckError() {
			return
		}
	}

	// Get k8s ApiServer by DcIdList
	tc.apiServers, tc.Ye = yceutils.GetApiServerList(tc.dcIdList)
	if tc.CheckError() {
		return
	}

	// Create k8s clients
	tc.k8sClients, tc.Ye = yceutils.CreateK8sClientList(tc.apiServers)
	if tc.CheckError() {
		return
	}

	tc.initTopology()

	if strings.EqualFold(userIdStr, ADMIN_ID) {
		tc.getTopologyForAllDcAndAllNamespaces()
	} else {
		tc.getTopologyForAllDc()
	}

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
