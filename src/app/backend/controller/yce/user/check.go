package user

import (
	myerror "app/backend/common/yce/error"
	"app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	myorganization "app/backend/model/mysql/organization"
)

type CheckUserController struct {
	yce.Controller
	params *CheckUserParams
	orgId  string
}

type CheckUserParams struct {
	UserName string `json:"userName"`
	OrgName  string `json:"orgName"`
	OrgId    string `json:"orgId"`
}

// check whether the name is existed
func (cuc *CheckUserController) checkDuplicatedName() {
	org := new(myorganization.Organization)
	err := org.QueryOrganizationByName(cuc.params.OrgName)
	if err != nil {
		cuc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	//_, ye := yceutils.QueryDuplicatedUserNameAndOrgId(cuc.params.UserName, org.Id)
	_, ye := yceutils.QueryDuplicatedUserName(cuc.params.UserName)
	// not found, can insert
	if ye != nil {
		//cuc.Ye = myerror.NewYceError(myerror.EOK, "")
		return
	} else {
		// found, cann't insert
		cuc.Ye = myerror.NewYceError(myerror.EYCE_EXISTED_NAME, "")
		return
	}

}

func (cuc CheckUserController) Post() {
	SessionIdFromClient := cuc.RequestHeader("Authorization")
	cuc.params = new(CheckUserParams)

	err := cuc.ReadJSON(cuc.params)
	if err != nil {
		cuc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cuc.CheckError() {
		return
	}

	cuc.orgId = cuc.params.OrgId

	cuc.ValidateSession(SessionIdFromClient, cuc.orgId)
	if cuc.CheckError() {
		return
	}

	// check if duplicated user name
	cuc.checkDuplicatedName()
	if cuc.CheckError() {
		return
	}

	cuc.WriteOk("")
	log.Infoln("CheckUserController Post Over")
}

