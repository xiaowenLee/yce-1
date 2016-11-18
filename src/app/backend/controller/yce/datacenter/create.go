package datacenter

import (
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	"strconv"
)

type CreateDatacenterController struct {
	yce.Controller

	params *CreateDatacenterParams
}

type CreateDatacenterParams struct {
	Name string `json:"name"`
	NodePort []string `json:"nodePort"`
	Host string `json:"host"`
	Port int32 `json:"port"`
	OrgId string `json:"orgId"`
	Op  int32 `json:"op"`
	//Secret string `json:"secret"` //TODO: will be realized later
	DcId int32 `json:"dcId"`
}

func (cdc *CreateDatacenterController) createNodePortItems() {

	nodePortLowerLimit, err := strconv.Atoi(cdc.params.NodePort[0])
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}
	nodePortUpperLimit, err := strconv.Atoi(cdc.params.NodePort[1])
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}

	if nodePortLowerLimit >= nodePortUpperLimit || yceutils.ValidatePort(nodePortLowerLimit) != nil || yceutils.ValidatePort(nodePortUpperLimit) != nil {
		cdc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
		return
	}

	cdc.Ye = yceutils.InitNodePortTableOfDatacenter(cdc.params.NodePort, cdc.params.DcId, cdc.params.Op)
	if cdc.Ye != nil {
		return
	}

}

func (cdc *CreateDatacenterController) createDcDbItems() {
	nodePort, ye := yceutils.EncodeNodePort(cdc.params.NodePort)
	if ye != nil {
		cdc.Ye = ye
		return
	}

	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterByName(cdc.params.Name)
	if err != nil {
		dc := mydatacenter.NewDataCenter(cdc.params.Name, cdc.params.Host, "", nodePort, "", cdc.params.Port, cdc.params.Op)
		err := dc.InsertDataCenter(cdc.params.Op)
		if err != nil {
			cdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
			return
		}
	} else if dc.Status == mydatacenter.INVALID {
		dc.Status = mydatacenter.VALID
		dc.NodePort, ye = yceutils.EncodeNodePort(cdc.params.NodePort)
		if ye != nil {
			cdc.Ye = ye
			return
		}
		err := dc.UpdateDataCenter(cdc.params.Op)
		if err != nil {
			cdc.Ye = myerror.NewYceError(myerror.EMYSQL, "")
			return
		}
	} else {
		cdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT,  "")
	}

	err = dc.QueryDataCenterByName(cdc.params.Name)
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
	}
	cdc.params.DcId = dc.Id

}

func (cdc CreateDatacenterController) Post() {
	cdc.params = new(CreateDatacenterParams)

	err := cdc.ReadJSON(cdc.params)
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cdc.CheckError() {
		return
	}

	cdc.createDcDbItems()
	if cdc.CheckError() {
		return
	}

	cdc.createNodePortItems()
	if cdc.CheckError() {
		return
	}

	cdc.WriteOk("")
	log.Infoln("CreateDatacenterController Create Over!")

}
