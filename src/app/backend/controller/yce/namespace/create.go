package namespace

import (
	myerror "app/backend/common/yce/error"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
	api "k8s.io/kubernetes/pkg/api"
	resource "k8s.io/kubernetes/pkg/api/resource"
	"encoding/json"
	"strconv"
	"strings"
	yce "app/backend/controller/yce"
)

type CreateNamespaceController struct {
	yce.Controller
	Param  *CreateNamespaceParam
	k8sClients []*client.Client
	apiServers []string
}


type CreateNamespaceParam struct {
	OrgId string `json:"orgId"`
	UserId int32 `json:"userId"`
	Name string `json:"name"`
	CpuQuota int32 `json:"cpuQuota"`
	MemQuota int32 `json:"memQuota"`
	Budget string `json:"budget"`
	Balance string `json:"balance"`
	DcIdList []int32 `json:"dcIdList"`
}

// Parse Namespace struct, insert into MySQL
func (cnc *CreateNamespaceController) createNamespaceDbItem() {

	dcIdList, err := json.Marshal(cnc.Param.DcIdList)
	if err != nil {
		cnc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	org := myorganization.NewOrganization(cnc.Param.Name, cnc.Param.Budget, cnc.Param.Balance, "", string(dcIdList),
		cnc.Param.CpuQuota, cnc.Param.MemQuota, cnc.Param.UserId)

	err = org.InsertOrganization()
	if err != nil {
		cnc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}

	mylog.Log.Infof("CreateNamespaceController createNamespaceDbItem success")

}

// Get ApiServer by dcId
func (cnc *CreateNamespaceController) getApiServerByDcId(dcId int32) string {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("getApiServerById QueryDataCenterById Error: err=%s", err)
		cnc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	apiServer := host + ":" + port

	mylog.Log.Infof("CreateDeployController getApiServerByDcId: apiServer=%s, dcId=%d", apiServer, dcId)
	return apiServer
}

// Get ApiServer List for dcIdList
func (cnc *CreateNamespaceController) getApiServerList(dcIdList []int32) {
	// Foreach dcIdList
	for _, dcId := range cnc.Param.DcIdList {
		// Get ApiServer
		apiServer := cnc.getApiServerByDcId(dcId)
		if strings.EqualFold(apiServer, "") {
			mylog.Log.Errorf("CreateDeployController getApiServerList Error")
			return
		}

		cnc.apiServers = append(cnc.apiServers, apiServer)
	}

	mylog.Log.Infof("CreateNamespaceController getApiServerList: len(apiServer)=%d", len(cnc.apiServers))

	return
}

// Create k8sClients for every ApiServer
func (cnc *CreateNamespaceController) createK8sClients() {

	// Foreach every ApiServer to create it's k8sClient
	cnc.k8sClients = make([]*client.Client, 0)

	for _, server := range cnc.apiServers {
		config := &restclient.Config{
			Host: server,
		}

		c, err := client.New(config)
		if err != nil {
			mylog.Log.Errorf("createK8sClient Error: err=%s", err)
			cnc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return
		}

		cnc.k8sClients = append(cnc.k8sClients, c)
		cnc.apiServers = append(cnc.apiServers, server)
		mylog.Log.Infof("Append a new client to cnc.k8sClients array: c=%p, apiServer=%s", c, server)
	}

	mylog.Log.Infof("CreateNamespaceController createK8sClients: len(k8sClient)=%d", len(cnc.k8sClients))
	return
}

// Create Namespace for every ApiServer
func (cnc *CreateNamespaceController) createNamespace() {
	namespace := new(api.Namespace)
	namespace.ObjectMeta.Name = cnc.Param.Name

	// Foreach every k8sClient to create namespace resource
	for index, cli := range cnc.k8sClients {
		_, err := cli.Namespaces().Create(namespace)
		if err != nil {
			mylog.Log.Errorf("createNamespace Error: apiServer=%s, namespace=%s, err=%s",
				cnc.apiServers[index], namespace, err)
			cnc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_NAMESPACE, "")
			return
		}
	}

	mylog.Log.Infof("CreateNamespaceController createNamespace success")

}

// TODO: 由于数据中心配额表,它定义了每个数据中心有不同的配额,第一版本默认每个数据中心都是一样的配额,第二版本在实现资源增减的逻辑
// Create ResourceQuota for every ApiServer
func (cnc *CreateNamespaceController) createResourceQuota() {
	resourceQuota := new(api.ResourceQuota)
	resourceQuota.ObjectMeta.Name = cnc.Param.Name + "-quota"

	// translate into "resource.Quantity"
	cpuQuota := resource.NewQuantity(int64(cnc.Param.CpuQuota) * CPU_MULTIPLIER, resource.DecimalSI)
	memQuota := resource.NewQuantity(int64(cnc.Param.MemQuota) * MEM_MULTIPLIER, resource.BinarySI)

	resourceQuota.Spec.Hard[api.ResourceCPU] = *cpuQuota
	resourceQuota.Spec.Hard[api.ResourceMemory] = *memQuota

	// Foreach every k8sClient to create resourceQuota
	for index, cli := range cnc.k8sClients {
		_, err := cli.ResourceQuotas(cnc.Param.Name).Create(resourceQuota)
		if err != nil {
			mylog.Log.Errorf("createResoruceQuota Error: apiServer=%s, namespace=%s, err=%s",
				cnc.apiServers[index], cnc.Param.Name, err)
			cnc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_NAMESPACE, "")
		}
	}

	mylog.Log.Infof("CreateNamespaceController createResourceQuota: create Resource Quota success")

}

// Post /api/v1/organizations
func (cnc *CreateNamespaceController) Post() {
	// Parse create organization params
	cnc.Param = new(CreateNamespaceParam)
	err := cnc.ReadJSON(cnc.Param)
	if err != nil {
		mylog.Log.Errorf("CreateNamespaceController ReadJSON Error: error=%s", err)
		cnc.Ye =  myerror.NewYceError(myerror.EJSON, "")
		cnc.WriteBack()
		return
	}

	// Validate Session
	sessionIdFromClient := cnc.RequestHeader("Authorization")
	mylog.Log.Debugf("CreateNamespaceController Params: sessionId=%s", sessionIdFromClient)

	cnc.ValidateSession(sessionIdFromClient, cnc.Param.OrgId)
	if cnc.CheckError() {
		return
	}

	// Create Organization struct and insert it into MySQL
	cnc.createNamespaceDbItem()
	if cnc.CheckError() {
		return
	}

	// Get DcIdList
	cnc.getApiServerList(cnc.Param.DcIdList)
	if cnc.CheckError() {
		return
	}

	// Create k8s clients
	cnc.createK8sClients()
	if cnc.CheckError() {
		return
	}

	// Create Namespace
	cnc.createNamespace()
	if cnc.CheckError() {
		return
	}

	// Create ResourceQuota
	cnc.createResourceQuota()
	if cnc.CheckError() {
		return
	}

	cnc.WriteOk("")
	mylog.Log.Infoln("CreateNamespaceController over!")
	return
}