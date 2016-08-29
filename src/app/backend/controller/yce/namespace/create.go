package namespace

import (
	"github.com/kataras/iris"
	"app/backend/common/util/session"
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
)


type CreateNamespaceController struct {
	*iris.Context
	Ye *myerror.YceError
}


type CreateNamespaceParams struct {
	OrgId string `json:"orgId"`
	Name string `json:"name"`
	CpuQuota int32 `json:"cpuQuota"`
	MemQuota int32 `json:"memQuota"`
	Budget int32 `json:"budget"`
	Balance int32 `json:"balance"`
	DcIdList []int32 `json:"dcIdList"`
}

func (cnc *CreateNamespaceController) WriteBack() {
	cnc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("CreateDeployController Response YceError: controller=%p, code=%d, note=%s", cnc, cnc.Ye.Code, myerror.Errors[cnc.Ye.Code].LogMsg)
	cnc.Write(cnc.Ye.String())
}

func (cnc *CreateNamespaceController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		cnc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		cnc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}
	return
}

// Parse Namespace struct, insert into MySQL
func (cnc *CreateNamespaceController) CreateNamespaceDbItem(org *myorganization.Organization) {

}


// Post /api/v1/organizations
func (cnc *CreateNamespaceController) Post() {
	sessionIdFromClient := cdc.RequestHeader("Authorization")
	// Validate OrgId error
	cdc.validateSession(sessionIdFromClient, orgId)

	if cdc.Ye != nil {
		cdc.WriteBack()
		return
	}


}