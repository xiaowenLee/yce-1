package endpoint


import (
	"github.com/kataras/iris"
	myorganizaiton "app/backend/model/mysql/organization"
	myerror "app/backend/common/yce/error"
	mylog "app/backend/common/util/log"
	"app/backend/model/yce/endpoint"
	"app/backend/common/yce/organization"
	"app/backend/common/util/session"
	"encoding/json"
)

type InitEndpointsController struct {
	*iris.Context
	org *myorganizaiton.Organization
	Init endpoint.InitEndpoints
	Ye *myerror.YceError
}

func (iec *InitEndpointsController) WriteBack() {
	iec.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateEndpointsController Response YceError: controller=%p, code=%d, note=%s", iec, iec.Ye.Code, myerror.Errors[iec.Ye.Code].LogMsg)
	iec.Write(iec.Ye.String())
}

func (iec *InitEndpointsController) String() string {
	data, err := json.Marshal(iec.Init)
	if err != nil {
		mylog.Log.Errorf("InitEndpointsController String() Marshal Error: err=%s", err)
		return ""
	}
	return string(data)
}

// Validate Session
func (iec *InitEndpointsController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		iec.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		// relogin
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		iec.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	return
}



func (iec InitEndpointsController) Get() {
	sessionIdFromClient := iec.RequestHeader("Authorization")
	orgId := iec.Param("orgId")


	// Validate OrgId error
	iec.validateSession(sessionIdFromClient, orgId)

	if iec.Ye != nil {
		iec.WriteBack()
		return
	}

	// Valid session
	org, err := organization.GetOrganizationById(orgId)
	iec.org = org
	iec.Init.OrgId = orgId
	iec.Init.OrgName = iec.org.Name

	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		iec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		iec.WriteBack()
		return
	}

	// Get Datacenters by a organization
	iec.Init.DataCenters, err = organization.GetDataCentersByOrganization(iec.org)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		iec.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		iec.WriteBack()
		return
	}

	iec.Ye = myerror.NewYceError(myerror.EOK, iec.String())
	iec.WriteBack()
	mylog.Log.Infoln("InitEndpointsController Get over!")
	return

}
