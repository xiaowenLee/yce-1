package service

import (
	myorganizaiton "app/backend/model/mysql/organization"
	mynodeport "app/backend/model/mysql/nodeport"
	myerror "app/backend/common/yce/error"
	"app/backend/model/yce/service"
	"app/backend/common/yce/organization"
	"encoding/json"
	yce "app/backend/controller/yce"
)

type InitServiceController struct {
	yce.Controller
	org *myorganizaiton.Organization
	Init service.InitService
}


func (isc *InitServiceController) String() string {
	data, err := json.Marshal(isc.Init)
	if err != nil {
		mylog.Log.Errorf("InitServiceController String() Marshal Error: err=%s", err)
		return ""
	}
	return string(data)
}

func (isc InitServiceController) Get() {
	sessionIdFromClient := isc.RequestHeader("Authorization")
	orgId := isc.Param("orgId")
	mylog.Log.Debugf("InitServiceController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// Validate OrgId error
	isc.ValidateSession(sessionIdFromClient, orgId)

	if isc.CheckError() {
		return
	}

	// Valid session
	org, err := organization.GetOrganizationById(orgId)
	isc.org = org
	isc.Init.OrgId = orgId
	isc.Init.OrgName = isc.org.Name
	mylog.Log.Debugf("InitServiceController Params: orgId=%s, orgName=%s", isc.Init.OrgId, isc.Init.OrgName)

	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		isc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
	}
	if isc.CheckError() {
		return
	}

	// Get Datacenters by a organization
	isc.Init.DataCenters, err = organization.GetDataCentersByOrganization(isc.org)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		isc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
	}
	if isc.CheckError() {
		return
	}

	//TODO: nodePort mysql model realize recommandOne()
	// Get one nodePort
	// isc.Init.Quotas, err = myqouta.QueryAllQuotas()
	isc.Init.NodePort = mynodeport.Recommand(isc.Init.DataCenters)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		isc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
	}
	if isc.CheckError() {
		return
	}

	isc.WriteOk(isc.String())
	mylog.Log.Infoln("InitServiceController Get over!")
	return
}


