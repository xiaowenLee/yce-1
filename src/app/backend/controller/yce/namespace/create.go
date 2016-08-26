package namespace

import (
	"github.com/kataras/iris"
	"app/backend/common/util/session"
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)


type CreateNamespaceController struct {
	*iris.Context
	Ye *myerror.YceError
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