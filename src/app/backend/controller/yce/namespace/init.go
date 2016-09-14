package namespace

import (
	"github.com/kataras/iris"
	"app/backend/common/util/session"
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
)

type InitNamespaceController struct {
	*iris.Context
	Ye *myerror.YceError
}

type InitNamespaceParams struct {
	Name string `json:"name"`
	OrgId  string `json:"orgId"`
}

func (inc *InitNamespaceController) WriteBack() {
	inc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateDeployController Response YceError: controller=%p, code=%d, note=%s", inc, inc.Ye.Code, myerror.Errors[inc.Ye.Code].LogMsg)
	inc.Write(inc.Ye.String())
}


func (inc *InitNamespaceController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		inc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		inc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}
	return
}


// POST /api/v1/organizations/init
func (inc *InitNamespaceController) Post() {

	initNamespaceParams := new(InitNamespaceParams)
	err := inc.ReadJSON(initNamespaceParams)
	if err != nil {
		mylog.Log.Errorf("InitNamespaceController ReadJSON Error: error=%s", err)
		inc.Ye = myerror.NewYceError(myerror.EJSON, "")
		inc.WriteBack()
		return
	}

	org := new(myorganization.Organization)
	err = org.QueryOrganizationByName(initNamespaceParams.Name)

	// Exists
	if err == nil {
		inc.Ye = myerror.NewYceError(myerror.EYCE_ORG_EXIST, "")
		inc.WriteBack()
		return
	}

	// Not Exists
	inc.Ye = myerror.NewYceError(myerror.EOK, "")
	inc.WriteBack()
	return
}