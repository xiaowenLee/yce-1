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
	if cdc.CheckError() {
		return
	}

	id, _ := strconv.Atoi(userId)

	user := new(myuser.User)
	user.QueryUserById(id)

	navList := user.NavList

	nlc.WriteOk(navList)
	log.Infoln("NavListController over!")

}
