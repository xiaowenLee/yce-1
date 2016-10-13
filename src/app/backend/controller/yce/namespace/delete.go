package namespace

import (
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	myerror "app/backend/common/yce/error"
	api "k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	myorganization "app/backend/model/mysql/organization"
)

type DeleteNamespaceController struct {
	yce.Controller
	apiServers []string
	k8sClients []*client.Client

	params *DeleteNamespaceParams
}

type DeleteNamespaceParams struct {
	OrgName string `json:"orgName"` // the name of namespace will be deleted
	OrgId   string `json:"orgId"`   // the orgId of admin
	Op      int32 `json:"op"`       // the modifiedOp, default is admin
	DcIdList []int32 `json:"dcIdList"`
}

func (dnc *DeleteNamespaceController) getDcIdList() {
	dcIdList, ye := yceutils.GetDcIdListByOrgName(dnc.params.OrgName)
	if ye != nil {
		dnc.Ye = ye
		return
	}

	dnc.params.DcIdList = dcIdList
}

func (dnc *DeleteNamespaceController) deleteNamespace() {
	namespace := new(api.Namespace)
	namespace.ObjectMeta.Name = dnc.params.OrgName

	// Foreach every k8sClient to delete namespace resource
	for index, cli := range dnc.k8sClients {
		err := cli.Namespaces().Delete(namespace.Name)
		if err != nil {
			log.Errorf("deleteNamespace Error: apiServer=%s, namespace=%s, err=%s",
				dnc.apiServers[index], namespace, err)
			dnc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_NAMESPACE, "")
			return
		}
	}

	log.Infof("DeleteNamespaceController deleteNamespace success")
}

func (dnc *DeleteNamespaceController) deleteQuota() {
	resourceQuota := new(api.ResourceQuota)
	resourceQuota.ObjectMeta.Name = dnc.params.OrgName + "-quota"

	// Foreach every k8sClient to create resourceQuota
	for index, cli := range dnc.k8sClients {
		//_, err := cli.ResourceQuotas(dnc.params.Name).Create(resourceQuota)
		err := cli.ResourceQuotas(dnc.params.OrgName).Delete(resourceQuota.Name)
		if err != nil {
			log.Errorf("deleteResoruceQuota Error: apiServer=%s, namespace=%s, err=%s",
				dnc.apiServers[index], dnc.params.OrgName, err)
			dnc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_NAMESPACE, "")
		}
		log.Infof("DeleteNamespaceController DeleteResourceQuota: resourceQuotaName=%s", resourceQuota.Name)
	}

	log.Infof("DeleteNamespaceController DeleteResourceQuota: delete Resource Quota success")
}

func (dnc *DeleteNamespaceController) deleteK8s() {
	dnc.deleteNamespace()
	if dnc.Ye != nil {
		return
	}
	dnc.deleteQuota()
	if dnc.Ye != nil {
		return
	}
}

func (dnc *DeleteNamespaceController) deleteNamespaceDbItem() {
	org := new(myorganization.Organization)
	err := org.QueryOrganizationByName(dnc.params.OrgName)
	if err != nil {
		dnc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	err = org.DeleteOrganization(dnc.params.Op)
	if err != nil {
		dnc.Ye = myerror.NewYceError(myerror.EMYSQL_DELETE, "")
		return
	}
}

func (dnc DeleteNamespaceController) Post() {
	dnc.params = new(DeleteNamespaceParams)

	err := dnc.ReadJSON(dnc.params)
	if err != nil {
		dnc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if dnc.CheckError() {
		return
	}

	dnc.getDcIdList()
	if dnc.CheckError() {
		return
	}

	// Get ApiServer List
	dnc.apiServers, dnc.Ye = yceutils.GetApiServerList(dnc.params.DcIdList)
	if dnc.CheckError() {
		return
	}


	// Create K8sClient List
	dnc.k8sClients, dnc.Ye = yceutils.CreateK8sClientList(dnc.apiServers)
	if dnc.CheckError() {
		return
	}
	dnc.deleteK8s()
	if dnc.CheckError() {
		return
	}
	dnc.deleteNamespaceDbItem()
	if dnc.CheckError() {
		return
	}

	dnc.WriteOk("")
	log.Infoln("DeleteNamespaceController Delete Over!")

}
