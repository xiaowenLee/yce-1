package namespace

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	myorganization "app/backend/model/mysql/organization"
	api "k8s.io/kubernetes/pkg/api"
	resource "k8s.io/kubernetes/pkg/api/resource"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

type CreateNamespaceController struct {
	yce.Controller
	Param      *CreateNamespaceParam
	k8sClients []*client.Client
	apiServers []string
}

type CreateNamespaceParam struct {
	UserId   int32   `json:"userId"`         // ModifiedOp, default is admin
	Name     string  `json:"name"`
	DcIdList []int32 `json:"dcIdList"`
	OrgId    string  `json:"orgId"`        // auto increase in MySQL
	CpuQuota int32   `json:"cpuQuota,omitempty"`     // blow 4 items will be modified in update.go
	MemQuota int32   `json:"memQuota,omitempty"`
	Budget   string  `json:"budget,omitempty"`
	Balance  string  `json:"balance,omitempty"`
}

// Parse Namespace struct, insert into MySQL
func (cnc *CreateNamespaceController) createNamespaceDbItem() {

	if cnc.Param.DcIdList == nil {
		cnc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}

	//dcIdList, err := json.Marshal(cnc.Param.DcIdList)
	dcIdList, ye := yceutils.EncodeDcIdList(cnc.Param.DcIdList)
	if ye != nil {
		cnc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	org := myorganization.NewOrganization(cnc.Param.Name, cnc.Param.Budget, cnc.Param.Balance, "", dcIdList,
		cnc.Param.CpuQuota, cnc.Param.MemQuota, cnc.Param.UserId)

	err := org.InsertOrganization()
	if err != nil {
		cnc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}

	log.Infof("CreateNamespaceController createNamespaceDbItem success")

}



// Create Namespace for every ApiServer
func (cnc *CreateNamespaceController) createNamespace() {
	namespace := new(api.Namespace)
	namespace.ObjectMeta.Name = cnc.Param.Name

	// Foreach every k8sClient to create namespace resource
	for index, cli := range cnc.k8sClients {
		_, err := cli.Namespaces().Create(namespace)
		if err != nil {
			log.Errorf("createNamespace Error: apiServer=%s, namespace=%s, err=%s",
				cnc.apiServers[index], namespace, err)
			cnc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_NAMESPACE, "")
			return
		}
	}

	log.Infof("CreateNamespaceController createNamespace success")

}

// TODO: 由于数据中心配额表,它定义了每个数据中心有不同的配额,第一版本默认每个数据中心都是一样的配额,第二版本在实现资源增减的逻辑
// Create ResourceQuota for every ApiServer
func (cnc *CreateNamespaceController) createResourceQuota() {
	resourceQuota := new(api.ResourceQuota)
	resourceQuota.ObjectMeta.Name = cnc.Param.Name + "-quota"
	resourceQuota.Spec.Hard = make(api.ResourceList, 0)

	// translate into "resource.Quantity"
	cpuQuota := resource.NewQuantity(int64(cnc.Param.CpuQuota)*CPU_MULTIPLIER, resource.DecimalSI)
	memQuota := resource.NewQuantity(int64(cnc.Param.MemQuota)*MEM_MULTIPLIER, resource.BinarySI)

	//TODO: didn't create quota or limits
	resourceQuota.Spec.Hard[api.ResourceCPU] = *cpuQuota
	resourceQuota.Spec.Hard[api.ResourceMemory] = *memQuota

	// Foreach every k8sClient to create resourceQuota
	for index, cli := range cnc.k8sClients {
		_, err := cli.ResourceQuotas(cnc.Param.Name).Create(resourceQuota)
		if err != nil {
			log.Errorf("createResoruceQuota Error: apiServer=%s, namespace=%s, err=%s",
				cnc.apiServers[index], cnc.Param.Name, err)
			cnc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_NAMESPACE, "")
		}
	}

	log.Infof("CreateNamespaceController createResourceQuota: create Resource Quota success")

}

// Post /api/v1/organizations
func (cnc CreateNamespaceController) Post() {
	// Parse create organization params
	cnc.Param = new(CreateNamespaceParam)
	err := cnc.ReadJSON(cnc.Param)
	if err != nil {
		log.Errorf("CreateNamespaceController ReadJSON Error: error=%s", err)
		cnc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cnc.CheckError() {
		return
	}

	// Validate Session
	sessionIdFromClient := cnc.RequestHeader("Authorization")
	log.Debugf("CreateNamespaceController Params: sessionId=%s", sessionIdFromClient)

	cnc.ValidateSession(sessionIdFromClient, cnc.Param.OrgId)
	if cnc.CheckError() {
		return
	}

	// Create Organization struct and insert it into MySQL
	cnc.createNamespaceDbItem()
	if cnc.CheckError() {
		return
	}

	// Get ApiServer List
	cnc.apiServers, cnc.Ye = yceutils.GetApiServerList(cnc.Param.DcIdList)
	if cnc.CheckError() {
		return
	}

	// Create K8sClient List
	cnc.k8sClients, cnc.Ye = yceutils.CreateK8sClientList(cnc.apiServers)
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
	log.Infoln("CreateNamespaceController over!")
	return
}
