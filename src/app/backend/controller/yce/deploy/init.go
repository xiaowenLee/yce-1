package deploy

import (
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	myorganization "app/backend/model/mysql/organization"
	myqouta "app/backend/model/mysql/quota"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
	mylog "app/backend/common/util/log"
)

var log =  mylog.Log

type InitDeployController struct {
	*iris.Context
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

// Validate Session
func (idc *InitDeployController) validateSession(sessionId, orgId string) (*myerror.YceError, error) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		idc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		idc.Write(errJson)
		return ye, err
	}

	// Session invalide
	if !ok {
		// relogin
		log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		idc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		idc.Write(errJson)
		return ye, err
	}

	return nil, nil
}

// GET /api/v1/organizations/{orgId}/users/{uid}/deployments/init
func (idc InitDeployController) Get() {
	sessionIdFromClient := idc.RequestHeader("Authorization")
	orgId := idc.Param("orgId")

	// Validate OrgId error
	ye, err := idc.validateSession(sessionIdFromClient, orgId)

	if ye != nil || err != nil {
		log.Errorf("ListDeployController validateSession: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		errJson, _ := ye.EncodeJson()
		idc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		idc.Write(errJson)
		return
	}

	// Valid session
	idc.org, err = organization.GetOrganizationById(orgId)
	idc.Init.OrgId = orgId
	idc.Init.OrgName = idc.org.Name

	if err != nil {
		log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		idc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		idc.Write(errJson)
		return
	}

	// Get Datacenters by a organization
	idc.Init.DataCenters, err = organization.GetDataCentersByOrganization(idc.org)
	if err != nil {
		log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		idc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		idc.Write(errJson)
		return
	}

	// Get all quotas
	idc.Init.Quotas, err = myqouta.QueryAllQuotas()
	if err != nil {
		log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		idc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		idc.Write(errJson)
		return
	}

	ye = myerror.NewYceError(0, "OK", idc.String())
	errJson, _ := ye.EncodeJson()
	idc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	idc.Write(errJson)

	log.Infoln("InitDeployController Get over!")
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
