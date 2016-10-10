package operationstat



type OperationStatController struct {
	yce.Controller
}

// GET /api/v1/organizations/{:orgId}/operationstat
func (osc OperationStatController) Get() {
	sessionIdFromClient := osc.RequestHeader("Authorization")
	orgId = osc.Param("orgId")

	// Validate OrgId error
	osc.ValidateSession(sessionIdFromClient, orgId)
	if osc.CheckError() {
		return
	}
}

