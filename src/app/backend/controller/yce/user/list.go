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
	OrgList  map[int32]string `json:"orgList"`
}

func (luc *ListUserController) getOrgNames() {
	organizations, ye := yceutils.GetAllOrganizations()
	if ye != nil {
		luc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	for _, org := range organizations {
		luc.params.OrgList[org.Id] = org.Name
	}
	log.Infof("ListUserController getOrgNames: len(OrgNames)=%d", len(luc.params.OrgList))
}

func (luc *ListUserController) getUsers() string {
	users, ye := yceutils.GetUsers()
	if ye != nil {
		luc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}
	luc.params.Users = users
	log.Infof("ListUserController getUsers: len(users)=%d", len(luc.params.Users))

	luc.getOrgNames()
	if luc.Ye != nil {
		return ""
	}
	log.Infof("ListUserController getUsers: params=%v", luc.params)

	usersJSON, err := json.Marshal(luc.params)
	if err != nil {
		log.Infof("ListUserController getUsers: error=%s", err)
		luc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	usersString := string(usersJSON)
	log.Infof("ListUserController: users=%s", usersString)
	return usersString

}

func (luc ListUserController) Get() {
	luc.params = new(UserList)
	luc.params.OrgList = make(map[int32]string)
	users := luc.getUsers()
	if luc.CheckError() {
		return
	}

	luc.WriteOk(users)
	log.Infoln("ListUserController Get Over!")
}


