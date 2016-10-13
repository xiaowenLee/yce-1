package namespace

import (
	myorganization "app/backend/model/mysql/organization"
	yceutils "app/backend/controller/yce/utils"
	yce "app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	resource "k8s.io/kubernetes/pkg/api/resource"
	api "k8s.io/kubernetes/pkg/api"
)

type UpdateNamespaceController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client

	params *UpdateNamespaceParam
}
type UpdateNamespaceParam struct {
	OrgId    string  `json:"orgId"`                  // auto increase in MySQL
	UserId   int32   `json:"userId"`                 // ModifiedOp, default is admin
	Name     string  `json:"name"`                   // forbid to be modified
	DcIdList []int32 `json:"dcIdList,omitempty"`     // modificaiton will be handle later
	CpuQuota int32   `json:"cpuQuota,omitempty"`     // blow 4 items will be modified in update.go
	MemQuota int32   `json:"memQuota,omitempty"`
	Budget   string  `json:"budget,omitempty"`
	Balance  string  `json:"balance,omitempty"`
}

func (unc *UpdateNamespaceController) updateNamespaceDbItem() {
	dcIdList, ye := yceutils.EncodeDcIdList(unc.params.DcIdList)
	if ye != nil {
		unc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	org := myorganization.NewOrganization(unc.params.Name, unc.params.Budget, unc.params.Balance, "", string(dcIdList),
		unc.params.CpuQuota, unc.params.MemQuota, unc.params.UserId)

	err  := org.UpdateOrganization(unc.params.UserId)

	if err != nil {
		unc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}

	log.Infof("UpdateNamespaceController updateNamespaceDbItem success")
}


func (unc *UpdateNamespaceController) updateResourceQuota() {
	resourceQuota := new(api.ResourceQuota)
	resourceQuota.ObjectMeta.Name = unc.params.Name + "-quota"
	resourceQuota.Spec.Hard = make(api.ResourceList, 0)

	// translate into "resource.Quantity"
	cpuQuota := resource.NewQuantity(int64(unc.params.CpuQuota)*CPU_MULTIPLIER, resource.DecimalSI)
	memQuota := resource.NewQuantity(int64(unc.params.MemQuota)*MEM_MULTIPLIER, resource.BinarySI)

	//TODO: didn't create quota or limits
	resourceQuota.Spec.Hard[api.ResourceCPU] = *cpuQuota
	resourceQuota.Spec.Hard[api.ResourceMemory] = *memQuota

	// Foreach every k8sClient to create resourceQuota
	for index, cli := range unc.k8sClients {
		_, err := cli.ResourceQuotas(unc.params.Name).Create(resourceQuota) //TODO: change to Update()
		if err != nil {
			log.Errorf("updateResoruceQuota Error: apiServer=%s, namespace=%s, err=%s",
				unc.apiServers[index], unc.params.Name, err)
			unc.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_NAMESPACE, "")
		}
	}

	log.Infof("UpdateNamespaceController UpdateResourceQuota: create Resource Quota success")

}

func (unc UpdateNamespaceController) Post() {
	unc.params = new(UpdateNamespaceParam)
	err := unc.ReadJSON(unc.params)
	if err != nil {
		log.Errorf("UpdateNamespaceController ReadJSON Error: error=%s", err)
		unc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}

	if unc.CheckError() {
		return
	}

	// Validate Session
	sessionIdFromClient := unc.RequestHeader("Authorization")
	log.Debugf("UpdateNamespaceController paramss: sessionId=%s", sessionIdFromClient)

	unc.ValidateSession(sessionIdFromClient, unc.params.OrgId)
	if unc.CheckError() {
		return
	}

	// Get ApiServer List
	unc.apiServers, unc.Ye = yceutils.GetApiServerList(unc.params.DcIdList)
	if unc.CheckError() {
		return
	}

	// Create K8sClient List
	unc.k8sClients, unc.Ye = yceutils.CreateK8sClientList(unc.apiServers)
	if unc.CheckError() {
		return
	}

	unc.updateResourceQuota()
	if unc.CheckError() {
		return
	}

	unc.updateNamespaceDbItem()
	if unc.CheckError() {
		return
	}

	unc.WriteOk("")
	log.Infoln("UpdateNamespaceController Post Over!")
	return
}