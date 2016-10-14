package user

import (
	yce "app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	myuser "app/backend/model/mysql/user"
)

type DeleteUserController struct {
	yce.Controller
	params *DeleteUserParams
}

type DeleteUserParams struct {
	UserName string `json:"userName"`
	Op       int32 `json:"op"`
}

func (duc *DeleteUserController) deleleUserDbItem() {
	user := new(myuser.User)
	err := user.QueryUserByUserName(duc.params.UserName)
	if err != nil {
		duc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	log.Infof("DeleteUserController: name=%s", user.Name)
	err = user.DeleteUser(duc.params.Op)
	if err != nil {
		duc.Ye = myerror.NewYceError(myerror.EMYSQL_DELETE, "")
		return
	}
}

func (duc DeleteUserController) Post() {
	duc.params = new(DeleteUserParams)
	err := duc.ReadJSON(duc.params)
	if err != nil {
		duc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if duc.CheckError() {
		return
	}

	duc.deleleUserDbItem()
	if duc.CheckError() {
		return
	}

	duc.WriteOk("")
	log.Infoln("DeleteUserController Delete Over!")

}
