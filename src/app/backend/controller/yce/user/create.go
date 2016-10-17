package user

import (
	"app/backend/controller/yce"
	myuser "app/backend/model/mysql/user"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
	"strconv"
	"app/backend/common/util/encrypt"
)

type CreateUserController struct {
	yce.Controller

	params *CreateUserParams
	orgId string
}

type CreateUserParams struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	OrgName  string `json:"orgName"`
	OrgId    string `json:"orgId"` // admin's organization
	Op 	 string `json:"op"`    // admin's userId
}

// create user.
func (cuc *CreateUserController) createUser() {
	//NOTE: not safe, for the plain password text transferred on the net.
	encryptPassword := encrypt.NewEncryption(cuc.params.Password).String()
	org := new(myorganization.Organization)
	err := org.QueryOrganizationByName(cuc.params.OrgName)
	if err != nil {
		cuc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	op, _ := strconv.Atoi(cuc.params.Op)
	user := myuser.NewUser(cuc.params.UserName, encryptPassword, "", org.Id, int32(op))

	err = user.InsertUser(int32(op))
	if err != nil {
		cuc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}
}

func (cuc CreateUserController) Post() {
	SessionIdFromClient := cuc.RequestHeader("Authorization")
	cuc.params = new(CreateUserParams)

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
	cuc.createUser()
	if cuc.CheckError() {
		return
	}

	cuc.WriteOk("")
	log.Infoln("CreateUserController Post Over!")

}

