package datacenter

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	mydatacenter "app/backend/model/mysql/datacenter"
)

type DeleteDatacenterController struct {
	yce.Controller

	params *DeleteDatacenterParams
}

type DeleteDatacenterParams struct {
	Name string `json:"name"`
	Op   int32  `json:"op"`

	OrgId string `json:"orgId"`
}

func (ddc *DeleteDatacenterController) deleteDcDbItem() {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterByName(ddc.params.Name)
	if err != nil {
		ddc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	dc.DeleteDataCenter(ddc.params.Op)
}

func (ddc DeleteDatacenterController) Post() {
	ddc.params = new(DeleteDatacenterParams)

	err := ddc.ReadJSON(ddc.params)
	if err != nil {
		ddc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
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
