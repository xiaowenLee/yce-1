package user

import (
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	myerror "app/backend/common/yce/error"
	myuser "app/backend/model/mysql/user"
	"encoding/json"
)

type ListUserController struct {
	yce.Controller

	params *UserList
}

type UserList struct {
	Users []myuser.User `json:"users"`
}

func (luc *ListUserController) getUsers() string {
	users, ye := yceutils.GetUsers()
	log.Infof("ListUserController: len(users)=%d", len(users))
	if ye != nil {
		luc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}
	luc.params.Users = users
	usersJSON, err := json.Marshal(luc.params)
	if err != nil {
		luc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	usersString := string(usersJSON)
	log.Infof("ListUserController: users=%s", usersString)
	return usersString

}

func (luc ListUserController) Get() {
	luc.params = new(UserList)
	users := luc.getUsers()
	if luc.CheckError() {
		return
	}

	luc.WriteOk(users)
	log.Infoln("ListUserController Get Over!")
}


