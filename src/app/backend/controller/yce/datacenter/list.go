package datacenter

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	mydatacenter "app/backend/model/mysql/datacenter"
	"encoding/json"
)

type ListDatacenterController struct {
	yce.Controller

	params *DatacenterList
}

type DatacenterList struct {
	Datacenters []mydatacenter.DataCenter `json:"datacenters"`
}

func (ldc *ListDatacenterController) listDatacenters() string {
	dcList, err := mydatacenter.QueryAllDatacenters()
	if err != nil {
		ldc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	ldc.params.Datacenters = dcList

	dcListJSON, err := json.Marshal(ldc.params)
	if err != nil {
		ldc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	dcListString := string(dcListJSON)
	return dcListString
}

func (ldc ListDatacenterController) Get() {
	ldc.params = new(DatacenterList)

	err := ldc.ReadJSON(ldc.params)
	if err != nil {
		ldc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if ldc.CheckError() {
		return
	}

	result := ldc.listDatacenters()
	if ldc.CheckError() {
		return
	}

	ldc.WriteOk(result)
	log.Infoln("ListDatacenterController List Over!")
}
