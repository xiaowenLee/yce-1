package logout

import (
	"log"
	// "fmt"
	// "strconv"
	"github.com/kataras/iris"
	mysession "app/backend/common/util/session"
	// myerror "app/backend/common/error"
)

type LogoutController struct {
	*iris.Context
}

// Check is logined

func (lc *LogoutController) checkLogin(sessionId string) (*mysession.Session, error) {

	ss := mysession.SessionStoreInstance()

	session, err := ss.Get(sessionId)

	if err != nil {
		log.Printf("Get session by sessionId error: sessionId=%s, err=%s\n", sessionId, err)
		return nil, err
	}

	if err == nil && session == nil {
		log.Printf("Not Login or Expirated: sessionId=%s\n", sessionId)
		return nil, nil
	}

	return session, true
}

func (lc *LogoutController) logout(sessionId string) error {

	ss := mysession.SessionStoreInstance()
	err := ss.Delete(sessionId)

	if err != nil {
		log.Printf("Delete session by sessionId error: sessionId=%s, err=%s\n", sessionId, err)
		return err
	}

	log.Printf("Delete session successfully: sessionId=%s", sessionId)
	return nil
}

// POST /api/v1/users/{email}/login
func (lc *LogoutController) Post() {
	email := lc.Param("email")
	sessionId := string(lc.FormValue("sessionId"))

	log.Printf("Logout: username=%s, sessionId=%s\n", email, sessionId)

	session, err := lc.checkLogin(sessionId)

	if err != nil {
		log.Printf("CheckLogin error: sessionId=%s, err=%s\n")
		// lc.Write("")
		return
	}

	if session != nil {
		err = lc.logout(sessionId)
		if err != nil {
			log.Println("Logout error: sessionId=%s, userName=%s, orgId=%s, err=%s\n",
				sessionId, session.UserName, session.OrgId, err)
			// lc.Write("")
			return
		}
	}

	log.Printf("Logout successfully: sessionId=%s, userName=%s, orgId=%s\n",
		sessionId, session.UserName, session.OrgId)
	return
	// lc.Write("UserName: " + email + ", SessionID: " + sessionId)

}