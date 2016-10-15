package datacenter

import (
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mydatacenter "app/backend/model/mysql/datacenter"
	myerror "app/backend/common/yce/error"
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

	udc.updateDcDbItem()
	if udc.CheckError() {
		return
	}

	udc.WriteOk("")
	log.Infoln("UpdateDatacenterController Update Over!")
}