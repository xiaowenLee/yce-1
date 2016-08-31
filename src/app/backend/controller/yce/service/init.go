package service

import (
	"github.com/kataras/iris"
	myorganizaiton "app/backend/model/mysql/organization"
	mynodeport "app/backend/model/mysql/nodeport"
	myerror "app/backend/common/yce/error"
	mylog "app/backend/common/util/log"
	"app/backend/model/yce/service"
	"app/backend/common/yce/organization"
	"app/backend/common/util/session"
	"encoding/json"
)

type InitServiceController struct {
	*iris.Context
	org *myorganizaiton.Organization
	Init service.InitService
	Ye *myerror.YceError
}

func (isc *InitServiceController) WriteBack() {
	isc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateServiceController Response YceError: controller=%p, code=%d, note=%s", isc, isc.Ye.Code, myerror.Errors[isc.Ye.Code].LogMsg)
	isc.Write(isc.Ye.String())
}

func (isc *InitServiceController) String() string {
	data, err := json.Marshal(isc.Init)
	if err != nil {
		mylog.Log.Errorf("InitServiceController String() Marshal Error: err=%s", err)
		return ""
	}
	return string(data)
}

// Validate Session
func (isc *InitServiceController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		isc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		// relogin
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		isc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	return
}



func (isc InitServiceController) Get() {
	sessionIdFromClient := isc.RequestHeader("Authorization")
	orgId := isc.Param("orgId")


	// Validate OrgId error
	isc.validateSession(sessionIdFromClient, orgId)

	if isc.Ye != nil {
		isc.WriteBack()
		return
	}

	// Valid session
	org, err := organization.GetOrganizationById(orgId)
	isc.org = org
	isc.Init.OrgId = orgId
	isc.Init.OrgName = isc.org.Name

	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		isc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		isc.WriteBack()
		return
	}

	// Get Datacenters by a organization
	isc.Init.DataCenters, err = organization.GetDataCentersByOrganization(isc.org)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		isc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		isc.WriteBack()
		return
	}

	//TODO: nodePort mysql model realize recommandOne()
	// Get one nodePort
	// isc.Init.Quotas, err = myqouta.QueryAllQuotas()
	isc.Init.NodePort = mynodeport.Recommand(isc.Init.DataCenters)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		isc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		isc.WriteBack()
		return
	}

	isc.Ye = myerror.NewYceError(myerror.EOK, isc.String())
	isc.WriteBack()
	mylog.Log.Infoln("InitServiceController Get over!")
	return

}


