package datacenter

import (
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
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
}

func (cdc *CreateDatacenterController) createDcDbItems() {
	nodePort, ye := yceutils.EncodeNodePort(cdc.params.NodePort)
	if ye != nil {
		cdc.Ye = ye
		return
	}
	dc := mydatacenter.NewDataCenter(cdc.params.Name, cdc.params.Host, "", nodePort, "", cdc.params.Port, cdc.params.Op)
	err := dc.InsertDataCenter(cdc.params.Op)
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}

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

	cdc.WriteOk("")
	log.Infoln("CreateDatacenterController Create Over!")

}
