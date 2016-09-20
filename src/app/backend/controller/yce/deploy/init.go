package deploy

import (
	myorganization "app/backend/model/mysql/organization"
	"app/backend/model/yce/deploy"
	"encoding/json"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type InitDeploymentController struct {
	yce.Controller
	org  *myorganization.Organization
	Init deploy.InitDeployment
}

func (idc *InitDeploymentController) String() string {
	data, err := json.Marshal(idc.Init)
	if err != nil {
		log.Errorf("InitDeploymentController String() Marshal Error: err=%s", err)
		return ""
	}
	return string(data)
}

// GET /api/v1/organizations/{orgId}/users/{uid}/deployments/init
func (idc InitDeploymentController) Get() {
	sessionIdFromClient := idc.RequestHeader("Authorization")
	orgId := idc.Param("orgId")

	log.Debugf("InitDeploymentController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// Valid session
	idc.ValidateSession(sessionIdFromClient, orgId)
	// Validate OrgId error
	if idc.CheckError() {
		return
	}

	idc.Init.OrgId = orgId
	idc.Init.OrgName, idc.Ye = yceutils.GetOrgNameByOrgId(orgId)
	if idc.CheckError() {
		return
	}

	// Get Datacenters
	idc.Init.DataCenters, idc.Ye = yceutils.GetDatacentersByOrgId(orgId)
	if idc.CheckError() {
		return
	}

	// idc.getAllQuotas()
	// Get All Quotas
	idc.Init.Quotas, idc.Ye = yceutils.GetAllQuotasOrderByCpu()
	if idc.CheckError() {
		return
	}

	idc.WriteOk(idc.String())
	log.Infoln("InitDeploymentController Get over!")

	return

}
