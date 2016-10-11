package operationstat

import (
	yce "app/backend/controller/yce"
	"strconv"
	myerror "app/backend/common/yce/error"
)


type OperationStatController struct {
	yce.Controller
}

// GET /api/v1/organizations/{:orgId}/operationstat
func (osc OperationStatController) Get() {
	sessionIdFromClient := osc.RequestHeader("Authorization")
	orgId := osc.Param("orgId")

	// Validate OrgId error
	osc.ValidateSession(sessionIdFromClient, orgId)
	if osc.CheckError() {
		return
	}

	id, err := strconv.Atoi(orgId)
	if err != nil {
		log.Errorf("Invalide OrgId value: err=%s", err)
		osc.Ye = myerror.NewYceError(myerror.EARGS, "")
	}

	if osc.CheckError() {
		return
	}

	ops := NewOperationStatistics()

	str, err := ops.Transform(int32(id))

	if err != nil {
		log.Errorf("OperationStatController Transform from mysql to json Error: err=%s", err)
		osc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
	}

	if osc.CheckError() {
		return
	}

	osc.WriteOk(str)
	log.Infoln("OperationStatController over!")

	return
}

