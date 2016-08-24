package logout

import (
	mysession "app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	"github.com/kataras/iris"
	mylog "app/backend/common/util/log"
)

var log =  mylog.Log

type LogoutController struct {
	*iris.Context
}

type LogoutParams struct {
	Username  string `json:"username"`
	SessionId string `json:"sessionId"`
}

// Check is logined
func (lc *LogoutController) checkLogin(sessionId string) (*mysession.Session, error) {

	ss := mysession.SessionStoreInstance()

	session, err := ss.Get(sessionId)

	if err != nil {
		log.Errorf("Get session by sessionId error: sessionId=%s, err=%s\n", sessionId, err)
		return nil, err
	}

	if err == nil && session == nil {
		log.Errorf("Not Login or Expirated: sessionId=%s\n", sessionId)
		return nil, nil
	}

	return session, nil
}

func (lc *LogoutController) logout(sessionId string) error {

	ss := mysession.SessionStoreInstance()
	err := ss.Delete(sessionId)

	if err != nil {
		log.Errorf("Delete session by sessionId error: sessionId=%s, err=%s\n", sessionId, err)
		return err
	}

	log.Infof("Delete session successfully: sessionId=%s", sessionId)
	return nil
}

// POST /api/v1/users/logout
func (lc LogoutController) Post() {

	logoutParams := new(LogoutParams)
	lc.ReadJSON(logoutParams)

	log.Infof("User Logout: username=%s, sessionId=%s\n", logoutParams.Username, logoutParams.SessionId)

	session, err := lc.checkLogin(logoutParams.SessionId)

	if err != nil {
		log.Errorf("CheckLogin error: sessionId=%s, err=%s\n")
		ye := myerror.NewYceError(1101, err.Error(), "")
		json, _ := ye.EncodeJson()
		lc.Write(json)
		return
	}

	if session != nil {
		err = lc.logout(logoutParams.SessionId)
		if err != nil {
			log.Errorf("Logout error: sessionId=%s, userName=%s, orgId=%s, err=%s\n",
				logoutParams.SessionId, session.UserName, session.OrgId, err)
			ye := myerror.NewYceError(1102, err.Error(), "")
			json, _ := ye.EncodeJson()
			lc.Write(json)
			return
		}
	}

	ye := myerror.NewYceError(0, "OK", "")
	json, _ := ye.EncodeJson()

	log.Infof("Logout successfully: sessionId=%s, userName=%s, orgId=%s\n",
		logoutParams.SessionId, session.UserName, session.OrgId)

	lc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	lc.Write(json)
	return

}
