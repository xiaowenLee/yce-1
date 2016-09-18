package deploy

import (
	mylog "app/backend/common/util/log"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	"strings"
	"sort"
	yce "app/backend/controller/yce"
)


type HistoryDeploymentController struct {
	yce.Controller
	apiServer  string
	k8sClient  *client.Client
	deployment *extensions.Deployment
	dcId       string
	orgId      string
	name       string          // deployment-name
	list       *ReplicaSetList // ReplicaSets
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

// Get ReplicaSet List by deployment via LabelSelectorAsSelector func
func (hdc *HistoryDeploymentController) getReplicaSetsByDeployment() []extensions.ReplicaSet {

	namespace := hdc.deployment.Namespace
	selector, err := unver.LabelSelectorAsSelector(hdc.deployment.Spec.Selector)
	if err != nil {
		log.Errorf("LabelSelectorAsSelector Error: apiServer=%s, namespace=%s, deployment=%s, err=%s",
			hdc.apiServer, namespace, hdc.deployment.Name, err)
		hdc.Ye = myerror.NewYceError(myerror.EKUBE_LABEL_SELECTOR, "")
		return nil
	}
	options := api.ListOptions{LabelSelector: selector}
	rsList, err := hdc.k8sClient.Extensions().ReplicaSets(namespace).List(options)

	log.Infof("HistoryDeploymentController GetReplicaSetByDeployment over!")
	return rsList.Items
}

// Get ApiServer by DcId
func (hdc *HistoryDeploymentController) getApiServerAndK8sClientByDcId() {

	// ApiServer
	dc := new(mydatacenter.DataCenter)
	dcId, _ := strconv.Atoi(hdc.dcId)
	err := dc.QueryDataCenterById(int32(dcId))
	if err != nil {
		log.Errorf("getApiServerById QueryDataCenterById Error: dcId=%d, err=%s", dcId, err)
		hdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	hdc.apiServer = host + ":" + port

	// K8sClient
	config := &restclient.Config{
		Host: hdc.apiServer,
	}

	c, err := client.New(config)
	if err != nil {
		log.Errorf("createK8sClient Error: err=%s", err)
		hdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	hdc.k8sClient = c
	log.Infof("GetApiServerAndK8sClientByDcId over: apiServer=%s, k8sClient=%p",
		hdc.apiServer, hdc.k8sClient)
}

// Get Deployment by deployment-name
func (hdc *HistoryDeploymentController) getDeploymentByName() {

	// Get namespace(org.Name) by orgId
	org, err := organization.GetOrganizationById(hdc.orgId)
	if err != nil {
		log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", hdc.orgId, err)
		hdc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	namespace := org.Name
	dp, err := hdc.k8sClient.Extensions().Deployments(namespace).Get(hdc.name)
	if err != nil {
		log.Errorf("getDeployByName Error: apiServer=%s, namespace=%s, deployment-name=%s, err=%s\n",
			hdc.apiServer, namespace, hdc.name, err)
		hdc.Ye = myerror.NewYceError(myerror.EKUBE_GET_DEPLOYMENT, "")
		return
	}

	hdc.deployment = dp

	log.Infof("GetDeploymentByName over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		hdc.apiServer, namespace, hdc.name, dp)
}

// Foreach ReplicaSets to return
func (hdc *HistoryDeploymentController) getReplicaSetList() {

	hdc.list = new(ReplicaSetList)

	// Get ReplicaSets By Deployment
	rsList := hdc.getReplicaSetsByDeployment()
	if rsList == nil {
		log.Errorf("GetReplicaSetList Error: hdc=%p, apiServer=%s, deployment-name=%s",
			hdc, hdc.apiServer, hdc.name)
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

		*hdc.list = append(*hdc.list, hr)
	}

	log.Infof("GetReplicaList over: len(rsList)=%d", len(rsList))
}

// Encode ReplicaSetList to string
func (hdc *HistoryDeploymentController) encodeReplicaSetList() string {
	// Sort the HistoryReturn List
	sort.Sort(hdc.list)
	data, err := json.Marshal(hdc.list)
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
	hdc.name = hdc.Param("name")
	sessionIdFromClient := hdc.RequestHeader("Authorization")

	log.Debugf("HistoryDeploymentController Params: sessionId=%s, orgId=%s, dcId=%s, name=%s", sessionIdFromClient, hdc.orgId, hdc.dcId, hdc.name)

	// ValidateSessionId
	hdc.ValidateSessionId(sessionIdFromClient, hdc.orgId)
	if hdc.CheckError() {
		return
	}

	// Get ApiServer and K8sClient
	hdc.getApiServerAndK8sClientByDcId()
	if hdc.CheckError() {
		return
	}

	// Get Deployment by name
	hdc.getDeploymentByName()
	if hdc.CheckError() {
		return
	}

	// Get ReplicaSets by deployment
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
