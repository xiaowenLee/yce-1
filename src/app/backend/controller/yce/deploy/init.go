package deploy

import (
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	myorganization "app/backend/model/mysql/organization"
	myqouta "app/backend/model/mysql/quota"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
	yce "app/backend/controller/yce"
)

type InitDeployController struct {
	yce.Controller
	org  *myorganization.Organization
	Init deploy.InitDeployment
}

func (idc *InitDeployController) String() string {
	data, err := json.Marshal(idc.Init)
	if err != nil {
		log.Errorf("InitDeployController String() Marshal Error: err=%s", err)
		return ""
	}
	return string(data)
}

func (idc *InitDeployController) getOrgName(orgId string) {
	org, err := organization.GetOrganizationById(orgId)

	if err != nil {
		log.Errorf("Get Organization By orgId error: orgId=%s, error=%s", orgId, err)
		idc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	idc.org = org
	idc.Init.OrgId = orgId
	idc.Init.OrgName = idc.org.Name
	log.Debugf("InitDeployController Params: org=%p, orgId=%s", idc.org, idc.Init.OrgId, idc.Init.OrgName)

}

func (idc *InitDeployController) getDatacenters() {
	// Get Datacenters by a organization
	idc.Init.DataCenters, err = organization.GetDataCentersByOrganization(idc.org)
	if err != nil {
		log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		idc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	return
}

func (idc *InitDeployController) getAllQuotas() {
	// Get all quotas
	idc.Init.Quotas, err = myqouta.QueryAllQuotas()
	if err != nil {
		log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		idc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	return
}

// GET /api/v1/organizations/{orgId}/users/{uid}/deployments/init
func (idc InitDeployController) Get() {
	sessionIdFromClient := idc.RequestHeader("Authorization")
	orgId := idc.Param("orgId")

	log.Debugf("InitDeployController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// Valid session
	idc.ValidateSession(sessionIdFromClient, orgId)
	// Validate OrgId error
	if idc.CheckError() {
		return
	}

	idc.getOrgName(orgId)
	if idc.CheckError() {
		return
	}

	idc.getDatacenters()
	if idc.CheckError() {
		return
	}

	idc.getAllQuotas()
	if idc.CheckError() {
		return
	}

	idc.WriteOk(idc.String())
	log.Infoln("InitDeployController Get over!")

	return

}
