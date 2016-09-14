package login

import (
	"app/backend/common/util/encrypt"
	mylog "app/backend/common/util/log"
	mysession "app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	myuser "app/backend/model/mysql/user"
	"github.com/kataras/iris"
	"strconv"
)

var log = mylog.Log

type LoginController struct {
	*iris.Context
	Ye *myerror.YceError
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lc *LoginController) WriteBack() {
	lc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("LoginController Response YceError: controller=%p, code=%d, note=%s", lc, lc.Ye.Code, myerror.Errors[lc.Ye.Code].LogMsg)
	lc.Write(lc.Ye.String())
}

// Check username && password
func (lc *LoginController) check(name, password string) *myuser.User {

	encryptPass := encrypt.NewEncryption(password).String()

	// New a User
	user := new(myuser.User)

	err := user.QueryUserByNameAndPassword(name, encryptPass)

	if err != nil {
		//log.Errorf("Can not find the user: username=%s, err=%s", name, err)
		mylog.Log.Errorf("Can not find the user: username=%s, err=%s", name, err)
		lc.Ye = myerror.NewYceError(myerror.EYCE_LOGIN, "")
		return nil
	}

	mylog.Log.Infof("LoginController check success: name=%s", user.Name)
	return user

}

// Store Session and Set cookie
func (lc *LoginController) session(user *myuser.User) *mysession.Session {

	// Store (id,orgId) in SessionStore
	id := strconv.Itoa(int(user.Id))
	orgId := strconv.Itoa(int(user.OrgId))

	session := mysession.NewSession(id, user.Name, orgId)

	ss := mysession.SessionStoreInstance()

	err := ss.Set(session)

	log.Infof("Session: sessionId=%s, userId=%s, userName=%s, orgId=%s", session.SessionId, session.UserId, session.UserName, session.OrgId)

	if err != nil {
		log.Errorf("Set session error: sessionId=%s, err=%s", session.SessionId, err)
		lc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return nil
	}

	return session

}

// POST /api/v1/users/login
func (lc LoginController) Post() {

	loginParams := new(LoginParams)

	err := lc.ReadJSON(loginParams)
	if err != nil {
		mylog.Log.Errorf("LoginController ReadJSON Error=%s", err)
		lc.Ye = myerror.NewYceError(myerror.EYCE_LOGIN, "")
		lc.WriteBack()
		return
	}

	user := lc.check(loginParams.Username, loginParams.Password)
	if lc.Ye != nil {
		lc.WriteBack()
		return
	}

	session := lc.session(user)
	if lc.Ye != nil {
		lc.WriteBack()
		return
	}

	// Auth pass
	sessionStr, _ := session.EncodeJson()
	log.Errorf("User Login: sessionId=%s, userId=%s, userName=%s, orgId=%s",
		session.SessionId, session.UserId, session.UserName, session.OrgId)

	lc.Ye = myerror.NewYceError(myerror.EOK, sessionStr)
	lc.WriteBack()
	return
}
