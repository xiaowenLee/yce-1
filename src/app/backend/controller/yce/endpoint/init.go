package endpoint

import (
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	yce "app/backend/controller/yce"
	myorganizaiton "app/backend/model/mysql/organization"
	"app/backend/model/yce/endpoint"
	"encoding/json"
)

type InitEndpointsController struct {
	yce.Controller
	org  *myorganizaiton.Organization
	Init endpoint.InitEndpoints
}

func (iec *InitEndpointsController) String() string {
	data, err := json.Marshal(iec.Init)
	if err != nil {
		mylog.Log.Errorf("InitEndpointsController String() Marshal Error: err=%s", err)
		return ""
	}
	return string(data)
}

func (iec InitEndpointsController) Get() {
	sessionIdFromClient := iec.RequestHeader("Authorization")
	orgId := iec.Param("orgId")
	mylog.Log.Debugf("InitEndpointsController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// Validate OrgId error
	iec.ValidateSession(sessionIdFromClient, orgId)
	if iec.CheckError() {
		return
	}

	org, err := organization.GetOrganizationById(orgId)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		iec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	if iec.CheckError() {
		return
	}

	iec.org = org
	iec.Init.OrgId = orgId
	iec.Init.OrgName = iec.org.Name
	mylog.Log.Debugf("InitEndpointsController Params: orgName=%s", iec.Init.OrgName)

	// Get Datacenters by a organization
	iec.Init.DataCenters, err = organization.GetDataCentersByOrganization(iec.org)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		iec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	if iec.CheckError() {
		return
	}

	iec.WriteOk(iec.String())
	mylog.Log.Infoln("InitEndpointsController Get over!")
	return

}
