package deploy

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	"encoding/json"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"sort"
	"strconv"
	"strings"
)

type HistoryDeploymentController struct {
	yce.Controller
	apiServer      string
	k8sClient      *client.Client
	deployment     *extensions.Deployment
	dcId           string
	orgId          string
	orgName        string
	deploymentName string          // deployment-name
	replicaSetList *ReplicaSetList // ReplicaSets
}

type ReplicaType struct {
	Current int32 `json:"current"`
	Desire  int32 `json:"desire"`
}

type HistoryReturn struct {
	Revision  string      `json:"revision"`
	Name      string      `json:"name"`
	Namespace string      `json:"namespace"`
	Image     string      `json:"image"`
	Selector  string      `json:"selector"`
	Replicas  ReplicaType `json:"replicas"`
}

type ReplicaSetList []HistoryReturn

// Sort interface
func (slice ReplicaSetList) Len() int {
	return len(slice)
}

func (slice ReplicaSetList) Less(i, j int) bool {
	iRevision, _ := strconv.Atoi(slice[i].Revision)
	jRevision, _ := strconv.Atoi(slice[j].Revision)
	return iRevision > jRevision
}

func (slice ReplicaSetList) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (hdc *HistoryDeploymentController) encodeMapToString(labels map[string]string) string {
	var ss []string
	for key, value := range labels {
		str := key + ":" + value
		ss = append(ss, str)
	}

	return strings.Join(ss, ",")
}

// Foreach ReplicaSets to return
func (hdc *HistoryDeploymentController) getReplicaSetList() {

	hdc.replicaSetList = new(ReplicaSetList)

	// Get ReplicaSets By Deployment
	rsList, ye := yceutils.GetReplicaSetsByDeployment(hdc.k8sClient, hdc.deployment)

	if ye != nil {
		log.Errorf("GetReplicaSetList Error: hdc=%p, apiServer=%s, deployment-name=%s",
			hdc, hdc.apiServer, hdc.deploymentName)
		hdc.Ye = ye
		return
	}

	for _, rs := range rsList {
		hr := HistoryReturn{}

		hr.Name = rs.Name
		hr.Namespace = rs.Namespace
		hr.Selector = hdc.encodeMapToString(rs.Spec.Selector.MatchLabels)
		hr.Image = rs.Spec.Template.Spec.Containers[0].Image
		hr.Revision = rs.Annotations[REVISION_ANNOTATION]
		hr.Replicas.Current = rs.Status.Replicas
		hr.Replicas.Desire = rs.Spec.Replicas

		log.Debugf("GetReplicaSetList replicaset: name=%s, namespace=%s, image=%s, revision=%s, current=%d, desired=%d",
			hr.Name, hr.Namespace, hr.Image, hr.Revision, hr.Replicas.Current, hr.Replicas.Desire)

		*hdc.replicaSetList = append(*hdc.replicaSetList, hr)
	}

	log.Infof("GetReplicaList over: len(rsList)=%d", len(rsList))
}

// Encode ReplicaSetList to string
func (hdc *HistoryDeploymentController) encodeReplicaSetList() string {
	// Sort the HistoryReturn List
	sort.Sort(hdc.replicaSetList)
	data, err := json.Marshal(hdc.replicaSetList)
	if err != nil {
		log.Errorf("EncodeReplicaSetList Error: err=%s", err)
		hdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	return string(data)
}

// GET /api/v1/organizations/{orgId}/datacenters/{dcId}/deployments/{name}/history
func (hdc HistoryDeploymentController) Get() {
	hdc.orgId = hdc.Param("orgId")
	hdc.dcId = hdc.Param("dcId")
	hdc.deploymentName = hdc.Param("name")
	sessionIdFromClient := hdc.RequestHeader("Authorization")

	log.Debugf("HistoryDeploymentController Params: sessionId=%s, orgId=%s, dcId=%s, name=%s", sessionIdFromClient, hdc.orgId, hdc.dcId, hdc.deploymentName)

	// ValidateSession
	hdc.ValidateSession(sessionIdFromClient, hdc.orgId)
	if hdc.CheckError() {
		return
	}

	// Get ApiServer
	dcId, _ := strconv.Atoi(hdc.dcId)
	hdc.apiServer, hdc.Ye = yceutils.GetApiServerByDcId(int32(dcId))
	if hdc.CheckError() {
		return
	}

	// Get K8sClient
	hdc.k8sClient, hdc.Ye = yceutils.CreateK8sClient(hdc.apiServer)
	if hdc.CheckError() {
		return
	}

	// Get Deployment by name
	hdc.orgName, hdc.Ye = yceutils.GetOrgNameByOrgId(hdc.orgId)
	if hdc.CheckError() {
		return
	}

	hdc.deployment, hdc.Ye = yceutils.GetDeploymentByNameAndNamespace(hdc.k8sClient, hdc.deploymentName, hdc.orgName)
	if hdc.CheckError() {
		return
	}

	// Get ReplicaSetList(HistoryReturn)
	hdc.getReplicaSetList()
	if hdc.CheckError() {
		return
	}

	// Return to browser
	ret := hdc.encodeReplicaSetList()
	if hdc.CheckError() {
		return
	}

	hdc.WriteOk(ret)
	log.Infoln("HistoryDeploymentController over!")
	return
}
