package deploy

import (
	mylog "app/backend/common/util/log"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	myorganization "app/backend/model/mysql/organization"
	myqouta "app/backend/model/mysql/quota"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
)

type InitDeployController struct {
	*iris.Context
	org  *myorganization.Organization
	Init deploy.InitDeployment
	Ye *myerror.YceError
}

func (idc *InitDeployController) WriteBack() {
	idc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateDeployController Response YceError: controller=%p, code=%d, note=%s", idc, idc.Ye.Code, myerror.Errors[idc.Ye.Code].LogMsg)
	idc.Write(idc.Ye.String())
}

func (idc *InitDeployController) String() string {
	data, err := json.Marshal(idc.Init)
	if err != nil {
		mylog.Log.Errorf("InitDeployController String() Marshal Error: err=%s", err)
		return ""
	}
	return string(data)
}

// Validate Session
func (idc *InitDeployController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		idc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		// relogin
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		idc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	mylog.Log.Infof("InitDeployController ValidateSession success")

	return
}

// GET /api/v1/organizations/{orgId}/users/{uid}/deployments/init
func (idc InitDeployController) Get() {
	sessionIdFromClient := idc.RequestHeader("Authorization")
	orgId := idc.Param("orgId")

	mylog.Log.Debugf("InitDeployController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	idc.validateSession(sessionIdFromClient, orgId)
	// Validate OrgId error
	if idc.Ye != nil {
		idc.WriteBack()
		return
	}

	// Valid session
	org, err := organization.GetOrganizationById(orgId)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		idc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		idc.WriteBack()
		return
	}

	idc.org = org
	idc.Init.OrgId = orgId
	idc.Init.OrgName = idc.org.Name
	mylog.Log.Debugf("InitDeployController Params: org=%p, orgId=%s", idc.org, idc.Init.OrgId, idc.Init.OrgName)

	// Get Datacenters by a organization
	idc.Init.DataCenters, err = organization.GetDataCentersByOrganization(idc.org)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		idc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		idc.WriteBack()
		return
	}

	// Get all quotas
	idc.Init.Quotas, err = myqouta.QueryAllQuotas()
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		idc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		idc.WriteBack()
		return
	}

	idc.Ye = myerror.NewYceError(myerror.EOK, idc.String())
	idc.WriteBack()
	mylog.Log.Infoln("InitDeployController Get over!")
	return

}

/*
response example:
{
	"code": 0,
	"message": "....",
	"data": {
		"orgId":  "1",
		"orgName": "Ops",
		"dataCenters": [
			{
				"id": "1",
				"name": "世纪互联",
				"budget": 10000000,
				"balance": 10000000
			},
			{
				"id": "2",
				"name": "电信机房",
				"budget": 10000000,
				"balance": 10000000
			},
			{
				"id": "3",
				"name": "电子城机房",
				"budget": 10000000,
				"balance": 10000000
			}
		],
		"quotas": [
			{
				"id": "1",
				"name": "2C4G50G",
				"cpu": 2,
				"mem": 4,
				"rbd": 50,
				"price": 1000
			},
			{
				"id": "2",
				"name": "4C8G100G",
				"cpu": 4,
				"mem": 8,
				"rbd": 100,
				"price": 18000
			},
			{
				"id": "3",
				"name": "4C16G200G",
				"cpu": 4,
				"mem" 16,
				"rbd": 200,
				"price": 2860

			}

		]
	}
}
*/
