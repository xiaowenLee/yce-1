package namespace

import (
	yce "app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	yceutils "app/backend/controller/yce/utils"
)

type CheckNamespaceController struct {
	yce.Controller

	params *CheckNamespaceParams
}

type CheckNamespaceParams struct {
	OrgName string `json:"orgName"`
	OrgId   string `json:"orgId"`
}

func (cnc *CheckNamespaceController) checkDuplicatedName() {
	_, ye := yceutils.QueryDuplicatedOrgName(cnc.params.OrgName)
	// not found
	if ye != nil {
		return
	}

	// found
	cnc.Ye = myerror.NewYceError(myerror.EYCE_EXISTED_NAME, "")
	return
}



func (cnc CheckNamespaceController) Post() {
	cnc.params = new(CheckNamespaceParams)


	// TODO rethink of this authentication way
	sessionIdFromClient := cnc.RequestHeader("Authorization")

	err := cnc.ReadJSON(cnc.params)
	if err != nil {
		cnc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cnc.CheckError() {
		return
	}

	// validate admin's session
	cnc.ValidateSession(sessionIdFromClient, cnc.params.OrgId)
	if cnc.CheckError() {
		return
	}



	//checkDuplicatedName
	cnc.checkDuplicatedName()
	if cnc.CheckError() {
		return
	}

	cnc.WriteOk("")
	log.Infoln("CheckNamespaceController Post Over!")
}
