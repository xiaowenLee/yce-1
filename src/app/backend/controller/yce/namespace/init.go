package namespace

import (
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
	yce "app/backend/controller/yce"
)

type InitNamespaceController struct {
	yce.Controller
}

type InitNamespaceParams struct {
	Name string `json:"name"`
	OrgId  string `json:"orgId"`
}


// POST /api/v1/organizations/init
func (inc *InitNamespaceController) Post() {

	initNamespaceParams := new(InitNamespaceParams)
	err := inc.ReadJSON(initNamespaceParams)
	if err != nil {
		log.Errorf("InitNamespaceController ReadJSON Error: error=%s", err)
		inc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if inc.CheckError() {
		return
	}

	org := new(myorganization.Organization)
	err = org.QueryOrganizationByName(initNamespaceParams.Name)

	// Exists
	if err == nil {
		inc.Ye = myerror.NewYceError(myerror.EYCE_ORG_EXIST, "")
	}

	if inc.CheckError() {
		return
	}

	// Not Exists
	inc.WriteOk("")
	return
}