package namespace

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	myorganization "app/backend/model/mysql/organization"
	"encoding/json"
)

type ListNamespaceController struct {
	yce.Controller

	params *NamespaceList
}

type NamespaceList struct {
	Organizations []myorganization.Organization `json:"organizations"`
}

func (lnc *ListNamespaceController) getNamespaceList() string {
	// get Namespace
	organizations, err := myorganization.QueryAllOrganizations()
	if err != nil {
		lnc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
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
