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
	"encoding/json"
)


type CreateNamespaceController struct {
	*iris.Context
	Ye *myerror.YceError
	Param  *CreateNamespaceParam
}


type CreateNamespaceParam struct {
	OrgId string `json:"orgId"`
	UserId int32 `json:"userId"`
	Name string `json:"name"`
	CpuQuota int32 `json:"cpuQuota"`
	MemQuota int32 `json:"memQuota"`
	Budget string `json:"budget"`
	Balance string `json:"balance"`
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
func (cnc *CreateNamespaceController) CreateNamespaceDbItem() {

	dcIdList, err := json.Marshal(cnc.Param.DcIdList)
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	// CreateNamespaceDbItem
	org := myorganization.NewOrganization(cnc.Param.Name, cnc.Param.Budget, "", string(dcIdList),
		cnc.Param.CpuQuota, cnc.Param.MemQuota, cnc.Param.UserId)

	err = org.InsertOrganization()
	if err != nil {

	}
}


// Post /api/v1/organizations
func (cnc *CreateNamespaceController) Post() {
	// Parse create organization params
	cnc.Param = new(CreateNamespaceParam)
	cnc.ReadJSON(cnc.Param)

	// Validate Session
	sessionIdFromClient := cdc.RequestHeader("Authorization")
	cdc.validateSession(sessionIdFromClient, cnc.Param.OrgId)

	if cdc.Ye != nil {
		cdc.WriteBack()
		return
	}


}