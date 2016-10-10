package user

import (
	"app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
	myuser "app/backend/model/user/organization"
)

type InitUserCreationController struct {
	yce.Controller
	params *InitUserCreationParams
	orgId int32
}

type InitUserCreationParams struct {
	UserName string `json:"userName"`
	OrgName string `json:"orgName"`
}

func (iucc *InitUserCreationController) checkDuplicatedName() {
	org := new(myorganization.Organization)
	err := org.QueryOrganizationByName(iucc.params.OrgName)
	if err != nil {
		iucc.Ye = myerror.NewYceError(myerror.EMYSQL, "")
		return ""
	}

	user := new(myuser.User)
	ok, ye := yceutils.QueryDuplicatedNameAndOrgId(iucc.params.UserName, org.Id)
	// not found, can insert
	if ye != nil {
		iucc.Ye = myerror.NewYceError(myerror.EOK, "")
	} else {
		// found, cann't insert
		iucc.Ye = myerror.NewYceError(myerror.EYCE_DUP_NAME, "")
	}

	if iucc.CheckError() {
		return
	}

}

func (iucc InitUserCreationController) Post() {
	SessionIdFromClient := iucc.RequestHeader("Authorization")
	iucc.params = new(InitUserCreationParams)

	err = iucc.ReadJSON(iucc.params)
	if err != nil {
		iucc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if iucc.CheckError() {
		return
	}

	iucc.ValidateSession(SessionIdFromClient, iucc.orgId)
	if iucc.CheckError() {
		return
	}

	// check if duplicated user name
	iucc.checkDuplicatedName()
	if iucc.CheckError() {
		return
	}

	log.Infoln("InitUserCreationController Post Over")
}

