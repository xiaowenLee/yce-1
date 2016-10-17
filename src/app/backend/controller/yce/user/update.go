package user

import (
	yce "app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	myuser "app/backend/model/mysql/user"
//	myorganization "app/backend/model/mysql/organization"
	"app/backend/common/util/encrypt"
)

type UpdateUserController struct {
	yce.Controller

	params *UpdateUserParams
}

type UpdateUserParams struct {
	OrgId    string `json:"orgId"`      // admin's orgId
	Op       int32  `json:"op"`         // admin's userId
	Name     string `json:"name"`       // userName, forbidden modified

	OrgName  string `json:"orgName,omitempty"`    // user's orgName
	Password string `json:"password,omitempty"`   // user's password
}


func (uuc *UpdateUserController) updateUser() {
	u := new(myuser.User)
	err := u.QueryUserByUserName(uuc.params.Name)
	if err != nil {
		uuc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	//TODO: Only allow modified password now. orgId is foreign key.
	/*
	if uuc.params.OrgName != "" {
		o := new(myorganization.Organization)
		o.QueryOrganizationByName(uuc.params.OrgName)
		u.OrgId = o.Id
	}
	*/

	if uuc.params.Password != "" {
		encryptPassword := encrypt.NewEncryption(uuc.params.Password).String()
		u.Password = encryptPassword
	}

	err = u.UpdateUser(uuc.params.Op)
	if err != nil {
		uuc.Ye = myerror.NewYceError(myerror.EMYSQL, "")
	}
}

func (uuc UpdateUserController) Post() {
	uuc.params = new(UpdateUserParams)
	err := uuc.ReadJSON(uuc.params)
	if err != nil {
		uuc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if uuc.CheckError() {
		return
	}

	uuc.updateUser()
	if uuc.CheckError() {
		return
	}

	uuc.WriteOk("")
	log.Infoln("UpdateUserController Update Over!")
}

