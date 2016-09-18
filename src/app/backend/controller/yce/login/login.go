package login

import (
	"app/backend/common/util/encrypt"
	myerror "app/backend/common/yce/error"
	myuser "app/backend/model/mysql/user"
	"strconv"
	yce "app/backend/controller/yce"
)


type LoginController struct {
	yce.Controller
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	}
	if lc.CheckError() {
		return
	}

	user := lc.check(loginParams.Username, loginParams.Password)
	if lc.CheckError() {
		return
	}

	session := lc.session(user)
	if lc.CheckError() {
		return
	}

	// Auth pass
	sessionStr, _ := session.EncodeJson()
	lc.WriteOk(sessionStr)
	log.Infof("User Login: sessionId=%s, userId=%s, userName=%s, orgId=%s",
		session.SessionId, session.UserId, session.UserName, session.OrgId)

	return
}
