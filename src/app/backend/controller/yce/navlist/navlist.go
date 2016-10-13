package navlist

import (
	yce "app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	myuser "app/backend/model/mysql/user"
	"strconv"
)

type NavListController struct {
	yce.Controller
}

// GET /api/v1/organizations/orgId/users/userId/navList

func (nlc NavListController) Get() {

	sessionIdFromClient := nlc.RequestHeader("Authorization")
	orgId := nlc.Param("orgId")
	userId := nlc.Param("userId")

	log.Debugf("CreateDeploymentController get Params:  sessionIdFromClient=%s, orgId=%s, userId=%s", sessionIdFromClient, orgId, userId)

	// Validate OrgId error
	nlc.ValidateSession(sessionIdFromClient, orgId)
	if nlc.CheckError() {
		return
	}

	id, _ := strconv.Atoi(userId)

	user := new(myuser.User)
	err := user.QueryUserById(int32(id))
	if err != nil {
		log.Errorf("NavListController QueryUserById Error: err=%s", err)
		nlc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
	}

	if nlc.CheckError() {
		return
	}

	navList := user.NavList

	nlc.WriteOk(navList)
	log.Infoln("NavListController over!")

}
