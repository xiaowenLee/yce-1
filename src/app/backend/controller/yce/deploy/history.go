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
)

const (
	REVISION_ANNOTATION string = "deployment.kubernetes.io/revision"
)

type HistoryDeployController struct {
	Ye *myerror.YceError
	*iris.Context
	apiServer  string
	k8sClient  *client.Client
	deployment *extensions.Deployment
	dcId       string
	orgId      string
	name       string          // deployment-name
	list       *ReplicaSetList // ReplicaSets
}

func (hdc *HistoryDeployController) WriteBack() {
	hdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("Create ListDeployController Response Error: controller=%p, code=%d, note=%s", hdc, hdc.Ye.Code, myerror.Errors[hdc.Ye.Code].LogMsg)
	hdc.Write(hdc.Ye.String())
}

type ReplicaType struct {
	Current int32 `json: "current"`
	Desire  int32 `json: "desire"`
}

type HistoryReturn struct {
	Revision  string      `json: "revision"`
	Name      string      `json: "name"`
	Namespace string      `json: "name"`
	Image     string      `json: "image"`
	Selector  string      `json: "image"`
	Replicas  ReplicaType `json: "replicas"`
}

type ReplicaSetList []HistoryReturn

func (hdc *HistoryDeployController) encodeMapToString(labels map[string]string) string {
	var ss []string
	for key, value := range labels {
		str := key + ":" + value
		ss = append(ss, str)
	}

	return strings.Join(ss, ",")
}

func (hdc *HistoryDeployController) validateSessionId(sessionId, orgId string) {
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	// validate error
	if err != nil {
		mylog.Log.Errorf("Create ListDeployController Error: sessionId=%s, orgId=%s, error=%s", sessionId, orgId, err)
		hdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// invalid sessionId
	if !ok {
		mylog.Log.Errorf("Create ListDeployController Failed: sessionId=%s, orgId=%s", sessionId, orgId)
		hdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	return
}

// Get ReplicaSet List by deployment via LabelSelectorAsSelector func
func (hdc *HistoryDeployController) getReplicaSetsByDeployment() []extensions.ReplicaSet {

	namespace := hdc.deployment.Namespace
	selector, err := unver.LabelSelectorAsSelector(hdc.deployment.Spec.Selector)
	if err != nil {
		mylog.Log.Errorf("LabelSelectorAsSelector Error: apiServer=%s, namespace=%s, deployment=%s, err=%s",
			hdc.apiServer, namespace, hdc.deployment.Name, err)
		hdc.Ye = myerror.NewYceError(myerror.EKUBE_LABEL_SELECTOR, "")
		return nil
	}
	options := api.ListOptions{LabelSelector: selector}
	rsList, err := hdc.k8sClient.Extensions().ReplicaSets(namespace).List(options)

	mylog.Log.Infof("HistoryDeployController GetReplicaSetByDeployment over!")
	return rsList.Items
}

// Get ApiServer by DcId
func (hdc *HistoryDeployController) getApiServerAndK8sClientByDcId() {

	// ApiServer
	dc := new(mydatacenter.DataCenter)
	dcId, _ := strconv.Atoi(hdc.dcId)
	err := dc.QueryDataCenterById(int32(dcId))
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: dcId=%d, err=%s", dcId, err)
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
		mylog.Log.Errorf("createK8sClient Error: err=%s", err)
		hdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	hdc.k8sClient = c
	mylog.Log.Infof("GetApiServerAndK8sClientByDcId over: apiServer=%s, k8sClient=%p",
		hdc.apiServer, hdc.k8sClient)
}

// Get Deployment by deployment-name
func (hdc *HistoryDeployController) getDeploymentByName() {

	// Get namespace(org.Name) by orgId
	org, err := organization.GetOrganizationById(hdc.orgId)
	if err != nil {
		mylog.Log.Errorf("getDatacentersByOrgId Error: orgId=%s, error=%s", hdc.orgId, err)
		hdc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}

	namespace := org.Name
	dp, err := hdc.k8sClient.Extensions().Deployments(namespace).Get(hdc.name)
	if err != nil {
		mylog.Log.Errorf("getDeployByName Error: apiServer=%s, namespace=%s, deployment-name=%s, err=%s\n",
			hdc.apiServer, namespace, hdc.name, err)
		hdc.Ye = myerror.NewYceError(myerror.EKUBE_GET_DEPLOYMENT, "")
		return
	}

	hdc.deployment = dp

	mylog.Log.Infof("GetDeploymentByName over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		hdc.apiServer, namespace, hdc.name, dp)
}

// Foreach ReplicaSets to return
func (hdc *HistoryDeployController) getReplicaSetList() {

	hdc.list = new(ReplicaSetList)

	// Get ReplicaSets By Deployment
	rsList := hdc.getReplicaSetsByDeployment()
	if rsList == nil {
		mylog.Log.Errorf("GetReplicaSetList Error: hdc=%p, apiServer=%s, deployment-name=%s",
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

		mylog.Log.Debugf("GetReplicaSetList replicaset: name=%s, namespace=%s, image=%s, revision=%s, current=%d, desired=%d",
			hr.Name, hr.Namespace, hr.Image, hr.Revision, hr.Replicas.Current, hr.Replicas.Desire)

		*hdc.list = append(*hdc.list, hr)
	}

	mylog.Log.Infof("GetReplicaList over!")
}

// Encode ReplicaSetList to string
func (hdc *HistoryDeployController) encodeReplicaSetList() string {
	data, err := json.Marshal(hdc.list)
	if err != nil {
		mylog.Log.Errorf("EncodeReplicaSetList Error: err=%s", err)
		hdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	return string(data)
}

// GET /api/v1/organizations/{orgId}/datacenters/{dcId}/deployments/{name}/history
func (hdc HistoryDeployController) Get() {
	hdc.orgId = hdc.Param("orgId")
	hdc.dcId = hdc.Param("dcId")
	hdc.name = hdc.Param("name")

	// validateSessionId
	sessionIdFromClient := hdc.RequestHeader("Authorization")
	hdc.validateSessionId(sessionIdFromClient, hdc.orgId)
	if hdc.Ye != nil {
		hdc.WriteBack()
		return
	}

	// Get ApiServer and K8sClient
	hdc.getApiServerAndK8sClientByDcId()
	if hdc.Ye != nil {
		hdc.WriteBack()
		return
	}

	// Get Deployment by name
	hdc.getDeploymentByName()
	if hdc.Ye != nil {
		hdc.WriteBack()
		return
	}

	// Get ReplicaSets by deployment
	hdc.getReplicaSetList()
	if hdc.Ye != nil {
		hdc.WriteBack()
		return
	}

	// Return to browser
	ret := hdc.encodeReplicaSetList()
	if hdc.Ye != nil {
		hdc.WriteBack()
		return
	}

	hdc.Ye = myerror.NewYceError(myerror.EOK, ret)
	hdc.WriteBack()
	mylog.Log.Infoln("HistoryDeployController over!")
	return
}
