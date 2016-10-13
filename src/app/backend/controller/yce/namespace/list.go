package namespace

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	myorganization "app/backend/model/mysql/organization"
	yceutils "app/backend/controller/yce/utils"
	"encoding/json"
)

type ListNamespaceController struct {
	yce.Controller

	params *NamespaceList
}

type NamespaceList struct {
	Organizations []myorganization.Organization `json:"organizations"`
	DcList        []yceutils.DcIdAndNameType `json:"dcList"`
}


func (lnc *ListNamespaceController) getDcList() {
	datacenters, ye := yceutils.QueryAllDatacenters()
	if ye != nil {
		lnc.Ye = ye
		return
	}

	for _, dc := range datacenters {
		d := new(yceutils.DcIdAndNameType)
		d.DcId = dc.Id
		d.DcName = dc.Name

		lnc.params.DcList = append(lnc.params.DcList, *d)
	}

}

func (lnc *ListNamespaceController) getNamespaceList() string {
	// get Namespace
	organizations, err := myorganization.QueryAllOrganizations()
	if err != nil {
		lnc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}

	lnc.getDcList()
	if lnc.Ye != nil {
		return ""
	}

	lnc.params.Organizations = organizations
	orgListJSON, err := json.Marshal(lnc.params)
	if err != nil {
		lnc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	orgListString := string(orgListJSON)
	return orgListString
}


func (lnc ListNamespaceController) Get() {
	//TODO: rethink of session authroization. Here it is omitted.
	//SessionIdFromClient := iuc.RequestHeader("Authorization")

	lnc.params = new(NamespaceList)

	orgList := lnc.getNamespaceList()
	if lnc.CheckError() {
		return
	}

	lnc.WriteOk(orgList)
	log.Infoln("ListNamespaceController Get Over!")
}
