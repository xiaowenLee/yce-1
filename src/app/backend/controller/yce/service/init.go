package service

import (
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	yce "app/backend/controller/yce"
	mynodeport "app/backend/model/mysql/nodeport"
	myorganizaiton "app/backend/model/mysql/organization"
	"app/backend/model/yce/service"
	"encoding/json"
)

type InitServiceController struct {
	yce.Controller
	org  *myorganizaiton.Organization
	Init service.InitService
}

func (isc *InitServiceController) String() string {
	data, err := json.Marshal(isc.Init)
	if err != nil {
		log.Errorf("InitServiceController String() Marshal Error: err=%s", err)
		return ""
	}
	return string(data)
}

//Get /api/v1/organizations/{:orgId}/users/{:userId}/services/init
func (isc InitServiceController) Get() {
	sessionIdFromClient := isc.RequestHeader("Authorization")
	orgId := isc.Param("orgId")
	log.Debugf("InitServiceController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

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
	log.Debugf("InitServiceController Params: orgId=%s, orgName=%s", isc.Init.OrgId, isc.Init.OrgName)

	if err != nil {
		log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		isc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
	}
	if isc.CheckError() {
		return
	}

	// Get Datacenters by a organization
	isc.Init.DataCenters, err = organization.GetDataCentersByOrganization(isc.org)
	if err != nil {
		log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		isc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
	}
	if isc.CheckError() {
		return
	}

	//TODO: nodePort mysql model realize Recommand()
	// Get one nodePort
	// isc.Init.Quotas, err = myqouta.QueryAllQuotas()
	isc.Init.NodePort = mynodeport.Recommand(isc.Init.DataCenters)
	if isc.CheckError() {
		return
	}

	log.Infof("InitServiceController: Recommand nodePort=%d", isc.Init.NodePort.Port)

	isc.WriteOk(isc.String())
	log.Infoln("InitServiceController Get over!")
	return
}
