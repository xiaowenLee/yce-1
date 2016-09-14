package logout

import (
	mylog "app/backend/common/util/log"
	mysession "app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	"github.com/kataras/iris"
)

var log = mylog.Log

type LogoutController struct {
	*iris.Context
	Ye *myerror.YceError
}

type LogoutParams struct {
	Username  string `json:"username"`
	SessionId string `json:"sessionId"`
}

func (lc *LogoutController) WriteBack() {
	lc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("LoginController Response YceError: controller=%p, code=%d, note=%s", lc, lc.Ye.Code, myerror.Errors[lc.Ye.Code].LogMsg)
	lc.Write(lc.Ye.String())
}

// Check is logined
func (lc *LogoutController) checkLogin(sessionId string) *mysession.Session {

	ss := mysession.SessionStoreInstance()

	session, err := ss.Get(sessionId)
	if err != nil {
		log.Errorf("Get session by sessionId error: sessionId=%s, err=%s", sessionId, err)
		lc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return nil
	}

	if err == nil && session == nil {
		log.Errorf("Not Login or Expirated: sessionId=%s", sessionId)
		lc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return nil
	}

	return session
}

func (lc *LogoutController) logout(sessionId string) {

	ss := mysession.SessionStoreInstance()
	err := ss.Delete(sessionId)

	if err != nil {
		log.Errorf("Delete session by sessionId error: sessionId=%s, err=%s", sessionId, err)
		lc.Ye = myerror.NewYceError(myerror.EYCE_SESSION_DEL, "")
		return
	}

	log.Infof("Delete session successfully: lc=%p, sessionId=%s", lc, sessionId)
	return
}

// POST /api/v1/users/logout
func (lc LogoutController) Post() {

	logoutParams := new(LogoutParams)
	err := lc.ReadJSON(logoutParams)
	if err != nil {
		mylog.Log.Errorf("LogoutController ReadJSON Error: error=%s", err)
		lc.Ye = myerror.NewYceError(myerror.EYCE_LOGOUT, "")
		lc.WriteBack()
		return
	}

	log.Infof("User Logout: username=%s, sessionId=%s", logoutParams.Username, logoutParams.SessionId)

	session := lc.checkLogin(logoutParams.SessionId)
	if lc.Ye != nil {
		lc.WriteBack()
		return
	}

	if session != nil {
		lc.logout(logoutParams.SessionId)
		if lc.Ye != nil {
			log.Errorf("Logout error: sessionId=%s, userName=%s, orgId=%s",
				logoutParams.SessionId, session.UserName, session.OrgId)
			lc.WriteBack()
			return
		}
	}

	lc.Ye = myerror.NewYceError(myerror.EOK, "")
	log.Infof("Logout successfully: sessionId=%s, userName=%s, orgId=%s",
		logoutParams.SessionId, session.UserName, session.OrgId)

	lc.WriteBack()
	return

}
