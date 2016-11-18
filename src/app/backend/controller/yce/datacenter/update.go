package datacenter

import (
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mydatacenter "app/backend/model/mysql/datacenter"
	myerror "app/backend/common/yce/error"
	"strconv"
)

type UpdateDatacenterController struct {
	yce.Controller

	params *UpdateDatacenterParams
}

type UpdateDatacenterParams struct {
	Name string `json:"name"`
	NodePort []string `json:"nodePort"`
	Host string `json:"host"`
	Port int32 `json:"port"`
	//Secret string `json:"secret"`
	Op int32 `json:"op"`

	OrgId string `json:"orgId"`
}

func (udc *UpdateDatacenterController) updateNodePortDbItem() {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterByName(udc.params.Name)
	if err != nil {
		udc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	//oldNodePort := dc.NodePort
	oldNodePort, ye := yceutils.DecodeNodePort(dc.NodePort)
	if ye != nil {
		udc.Ye = ye
		return
	}

	nodePortLowerLimit, err := strconv.Atoi(oldNodePort[0])
	if err != nil {
		udc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}
	nodePortUpperLimit, err := strconv.Atoi(oldNodePort[1])
	if err != nil {
		udc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}

	udc.Ye = yceutils.ValidateNodePort(int32(nodePortLowerLimit), int32(nodePortUpperLimit))
	if udc.Ye != nil {
		return
	}

	//Delete old nodePorts, turn them into INVALID
	ye = yceutils.DeleteNodePortTableOfDatacenter(oldNodePort, dc.Id, udc.params.Op)
	if ye != nil {
		udc.Ye = ye
		return
	}

	newNodePort := udc.params.NodePort

	nodePortLowerLimit, err = strconv.Atoi(udc.params.NodePort[0])
	if err != nil {
		udc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}
	nodePortUpperLimit, err = strconv.Atoi(udc.params.NodePort[1])
	if err != nil {
		udc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}

	udc.Ye = yceutils.ValidateNodePort(int32(nodePortLowerLimit), int32(nodePortUpperLimit))
	if udc.Ye != nil {
		return
	}

	ye = yceutils.InitNodePortTableOfDatacenter(newNodePort, dc.Id, udc.params.Op)
	if ye != nil {
		udc.Ye = ye
		return
	}

}


func (udc *UpdateDatacenterController) updateDcDbItem() {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterByName(udc.params.Name)
	if err != nil {
		udc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	dc.NodePort, udc.Ye = yceutils.EncodeNodePort(udc.params.NodePort)
	dc.Host = udc.params.Host
	dc.Port = udc.params.Port
	if udc.Ye != nil {
		return
	}
	dc.UpdateDataCenter(udc.params.Op)
}

func (udc UpdateDatacenterController) Post() {
	udc.params = new(UpdateDatacenterParams)

	err := udc.ReadJSON(udc.params)
	if err != nil {
		udc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if udc.CheckError() {
		return
	}

	udc.updateNodePortDbItem()
	if udc.CheckError() {
		return
	}

	udc.updateDcDbItem()
	if udc.CheckError() {
		return
	}

	udc.WriteOk("")
	log.Infoln("UpdateDatacenterController Update Over!")
}