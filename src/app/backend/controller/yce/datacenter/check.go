package datacenter

import (
	yce "app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	yceutils "app/backend/controller/yce/utils"
)

type CheckDatacenterController struct{
	yce.Controller

	params *CheckDatacenterParams
}

type CheckDatacenterParams struct {
	Name string `json:"name"`
	OrgId string `json:"orgId"`
}

func (cdc *CheckDatacenterController) checkDuplicatedDc() {
	_, ye := yceutils.QueryDuplicatedDcName(cdc.params.Name)
	if ye != nil {
		return
	}

	cdc.Ye = myerror.NewYceError(myerror.EYCE_EXISTED_NAME, "")
	return
}

func (cdc CheckDatacenterController) Post() {
	cdc.params = new(CheckDatacenterParams)

	err := cdc.ReadJSON(cdc.params)
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cdc.CheckError() {
		return
	}

	cdc.checkDuplicatedDc()
	if cdc.CheckError() {
		return
	}

	cdc.WriteOk("")
	log.Infoln("CheckDatacenterController Check Over!")
}
