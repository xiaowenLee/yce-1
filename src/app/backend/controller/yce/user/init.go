package user

import (
	"app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
)

type InitUserController struct {
	yce.Controller
	params *InitUserParams
	orgId string
}

type InitUserParams struct {
	UserName string `json:"userName"`
	OrgName string `json:"orgName"`
	OrgId string `json:"orgId"`
}

// check whether the name is existed
func (iuc *InitUserController) checkDuplicatedName() {
	org := new(myorganization.Organization)
	err := org.QueryOrganizationByName(iuc.params.OrgName)
	if err != nil {
		iuc.Ye = myerror.NewYceError(myerror.EMYSQL, "")
		return
	}

	_, ye := yceutils.QueryDuplicatedNameAndOrgId(iuc.params.UserName, org.Id)
	// not found, can insert
	if ye != nil {
		//iuc.Ye = myerror.NewYceError(myerror.EOK, "")
		return
	} else {
		// found, cann't insert
		iuc.Ye = myerror.NewYceError(myerror.EYCE_EXISTED_NAME, "")
		return
	}


}

func (iuc InitUserController) Post() {
	SessionIdFromClient := iuc.RequestHeader("Authorization")
	iuc.params = new(InitUserParams)

	err := iuc.ReadJSON(iuc.params)
	if err != nil {
		iuc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if iuc.CheckError() {
		return
	}

	iuc.orgId = iuc.params.OrgId

	iuc.ValidateSession(SessionIdFromClient, iuc.orgId)
	if iuc.CheckError() {
		return
	}

	// check if duplicated user name
	iuc.checkDuplicatedName()
	if iuc.CheckError() {
		return
	}

	iuc.WriteOk("")
	log.Infoln("InitUserController Post Over")
}

