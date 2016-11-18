package datacenter

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	mydatacenter "app/backend/model/mysql/datacenter"
	yceutils "app/backend/controller/yce/utils"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/api"
	"strconv"
)

type DeleteDatacenterController struct {
	yce.Controller
	k8sClient  *client.Client
	apiServer  string

	dc *mydatacenter.DataCenter
	params *DeleteDatacenterParams
}

type DeleteDatacenterParams struct {
	Name string `json:"name"`
	Op   int32  `json:"op"`

	OrgId string `json:"orgId"`
}

// Publish k8s.Service to every datacenter which in dcIdList
func (ddc *DeleteDatacenterController) deleteService(namespace, svcName string) {

	svcList, err := ddc.k8sClient.Services(api.NamespaceAll).List(api.ListOptions{})
	if err != nil {
		log.Errorf("ListService Error: apiServer=%s, namespace=%s, error=%s", ddc.apiServer, api.NamespaceAll, err)
		ddc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_SERVICE, "")
		return
	}

	for _, svc := range svcList.Items {
 		//TODO: same name conflicts ? testsvc in ops and testsvc in dev, how to display them in namespace all ?
		err := ddc.k8sClient.Services(api.NamespaceAll).Delete(svc.Name)
		if err != nil {
			log.Errorf("deleteService Error: apiServer=%s, namespace=%s, error=%s", ddc.apiServer, api.NamespaceAll, err)
			ddc.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_SERVICE, "")
			return
		}
	}

	log.Infof("Delete Service successfully: namespace=%s, apiServer=%s", namespace, ddc.apiServer)

	return
}

func (ddc *DeleteDatacenterController) deleteNodePortDbItem() {

	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterByName(ddc.params.Name)
	if err != nil {
		ddc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	log.Infof("DeleteDatacenterController deleteNodePortDbItem: id=%d, nodePort=%s", dc.Id, dc.NodePort)
	nodePortList, ye := yceutils.DecodeNodePort(dc.NodePort)
	if ye != nil {
		ddc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	nodePortLowerLimit, err := strconv.Atoi(nodePortList[0])
	if err != nil {
		ddc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}
	nodePortUpperLimit, err := strconv.Atoi(nodePortList[1])
	if err != nil {
		ddc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}

	ddc.Ye = yceutils.ValidateNodePort(int32(nodePortLowerLimit), int32(nodePortUpperLimit))
	if ddc.Ye != nil {
		return
	}


	//ddc.Ye = yceutils.InitNodePortTableOfDatacenter(cdc.params.NodePort, cdc.params.DcId, cdc.params.Op)
	ddc.Ye = yceutils.DeleteNodePortTableOfDatacenter(nodePortList, dc.Id, ddc.params.Op)
	if ddc.Ye != nil {
		return
	}
}

func (ddc *DeleteDatacenterController) deleteDcDbItem() {
	ddc.dc = new(mydatacenter.DataCenter)
	err := ddc.dc.QueryDataCenterByName(ddc.params.Name)
	if err != nil {
		ddc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	ddc.dc.DeleteDataCenter(ddc.params.Op)
}



func (ddc DeleteDatacenterController) Post() {
	ddc.params = new(DeleteDatacenterParams)
	ddc.dc = new(mydatacenter.DataCenter)

	err := ddc.ReadJSON(ddc.params)
	if err != nil {
		ddc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if ddc.CheckError() {
		return
	}

	dcId, ye := yceutils.GetDcIdByDcName(ddc.params.Name)
	if ye != nil {
		ddc.Ye = ye
	}
	if ddc.CheckError() {
		return
	}

	ddc.apiServer, ddc.Ye = yceutils.GetApiServerByDcId(dcId)
	if ddc.CheckError() {
		return
	}

	ddc.k8sClient, ddc.Ye = yceutils.CreateK8sClient(ddc.apiServer)
	if ddc.CheckError() {
		return
	}

	//TODO: delete all services in this datacenter. Need to think.
	//TODO: deleting means cann't be access temperorily or clean up all the stuff ?
	//TODO: modify user, organization existed in this datacenter, change to INVALID
	//ddc.deleteService()


	ddc.deleteNodePortDbItem()
	if ddc.CheckError() {
		return
	}

	ddc.deleteDcDbItem()
	if ddc.CheckError() {
		return
	}


	ddc.WriteOk("")
	log.Infoln("DeleteDatacenterController Delete Over!")
}
