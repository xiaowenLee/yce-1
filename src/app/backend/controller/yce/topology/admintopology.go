package topology

/*
import (
	"encoding/json"
	"strconv"
)


import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mydatacenter "app/backend/model/mysql/datacenter"
	"encoding/json"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	deployutil "k8s.io/kubernetes/pkg/controller/deployment/util"
	"strconv"
)

type DcList struct {
	Items []string `json:"dcList"`
}

type AdminTopologyController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
	topology   *Topology
	dcIdList   []int32
}

func (atc *AdminTopologyController) getAllDcIdList() error {
	dcs, err := mydatacenter.QueryAllDatacenters()

	if err != nil {
		atc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return err
	}

	atc.dcIdList = make([]int32, 0)

	for _, dc := range dcs {
		atc.dcIdList = append(atc.dcIdList, dc.Id)
	}
	return nil
}

func (atc *AdminTopologyController) getTopology(c *client.Client, namespace string) bool {
	// Get Deployments.List
	dpList, ye := yceutils.GetDeploymentByNamespace(c, namespace)

	if ye != nil {
		atc.Ye = ye
	}
	if atc.CheckError() {
		return false
	}

	// Foreach Deployments.List
	for _, dp := range dpList {
		//rsList, err := getReplicaSetsByDeployment(c, &dp)
		rsList, ye := yceutils.GetReplicaSetsByDeployment(c, &dp)
		if ye != nil {
			atc.Ye = ye
		}
		if atc.CheckError() {
			return false
		}

		// For all replicasets
		rs, err := deployutil.FindNewReplicaSet(&dp, rsList)
		if err != nil {
			log.Errorf("FindNewReplicaSet Error: err=%s", err)
			atc.Ye = myerror.NewYceError(myerror.EKUBE_FIND_NEW_REPLICASET, "")
			return false
		}

		// Topology.Items
		uid := string(rs.UID)
		atc.topology.Items[uid] = ReplicaSetType{
			Kind:       "ReplicaSet",
			ApiVersion: "v1beta2",
			ReplicaSet: *rs,
		}

		podList, ye := yceutils.GetPodListByReplicaSet(c, rs)

		if ye != nil {
			atc.Ye = ye
		}
		if atc.CheckError() {
			return false
		}

		for _, pod := range podList.Items {
			uid = string(pod.UID)
			atc.topology.Items[uid] = PodType{
				Kind:       "Pod",
				ApiVersion: "v1beta2",
				Pod:        pod,
			}

			relation := RelationsType{
				Source: string(rs.UID),
				Target: string(pod.UID),
			}

			atc.topology.Relations = append(atc.topology.Relations, relation)

			node, ye := yceutils.GetNodeByPod(c, &pod)
			if ye != nil {
				atc.Ye = ye
			}
			if atc.CheckError() {
				return false
			}

			uid = string(node.UID)
			atc.topology.Items[uid] = NodeType{
				Kind:       "Node",
				ApiVersion: "v1beata2",
				Node:       *node,
			}

			relation = RelationsType{
				Source: string(node.UID),
				Target: string(pod.UID),
			}
			atc.topology.Relations = append(atc.topology.Relations, relation)
		}
		// }
		// For all replicasets
	}

	// Get Services.List
	svcList, ye := yceutils.GetServicesByNamespace(c, namespace)
	if ye != nil {
		atc.Ye = ye
	}
	if atc.CheckError() {
		return false
	}

	for _, svc := range svcList {
		uid := string(svc.UID)
		atc.topology.Items[uid] = ServiceType{
			Kind:       "Service",
			ApiVersion: "v1beta1",
			Service:    svc,
		}

		podList, ye := yceutils.GetPodsByService(c, &svc)
		if ye != nil {
			atc.Ye = ye
		}
		if atc.CheckError() {
			return false
		}
		for _, pod := range podList {
			uid = string(pod.UID)
			if _, ok := atc.topology.Items[uid]; ok {
				atc.topology.Items[uid] = PodType{
					Kind:       "Pod",
					ApiVersion: "v1beta1",
					Pod:        pod,
				}
			}

			relation := RelationsType{
				Source: string(svc.UID),
				Target: string(pod.UID),
			}

			atc.topology.Relations = append(atc.topology.Relations, relation)
		}
	}

	return true
}

func (atc *AdminTopologyController) encodeTopology() string {

	data, err := json.Marshal(atc.topology)
	if err != nil {
		log.Errorf("encodeTopology Error: err=%s\n", err)
		atc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}
	return string(data)
}

func (atc *AdminTopologyController) initTopology() {
	atc.topology = new(Topology)
	atc.topology.Items = make(ItemType)
	atc.topology.Relations = make([]RelationsType, 0)
}

func (atc *AdminTopologyController) getTopologyForAllDc() {
	for index, client := range atc.k8sClients {
		namespaces, ye := yceutils.GetAllNamespaces(client)
		if ye != nil {
			atc.Ye = ye
			return
		}

		for _, name := range namespaces {
			atc.getTopology(client, name)
			log.Infof("Get Topology data for every datacenter: apiServer=%s, client=%p, namespace=%s\n", atc.apiServers[index], client, name)
		}
	}
}

// GET /api/v1/organizations/{orgId}/topology
func (atc AdminTopologyController) Get() {
	sessionIdFromClient := atc.RequestHeader("Authorization")
	orgIdStr := atc.Param("orgId")

	log.Debugf("TopologyController Params: sessionId=%s, orgIdStr=%s", sessionIdFromClient, orgIdStr)

	orgId, err := strconv.Atoi(orgIdStr)
	if err != nil {
		atc.Ye = myerror.NewYceError(myerror.EARGS, "")
	}
	if atc.CheckError() {
		return
	}

	atc.orgId = int32(orgId)

	// Validate OrgId error
	atc.ValidateSession(sessionIdFromClient, orgIdStr)
	if atc.CheckError() {
		return
	}

	// Get all dcId for all datacenters
	atc.getAllDcIdList()
	if atc.CheckError() {

	}

	// Get k8s ApiServer by DcIdList
	atc.apiServers, atc.Ye = yceutils.GetApiServerList(atc.dcIdList)
	if atc.CheckError() {
		return
	}

	// Create k8s clients
	atc.k8sClients, atc.Ye = yceutils.CreateK8sClientList(atc.apiServers)
	if atc.CheckError() {
		return
	}

	atc.initTopology()

	atc.getTopologyForAllDc()
	if atc.CheckError() {
		return
	}

	str := atc.encodeTopology()
	if atc.CheckError() {
		return
	}

	atc.WriteOk(str)
	log.Infoln("ToplogyController over!")

	return
}
*/
