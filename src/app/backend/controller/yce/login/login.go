package login

import (
	"app/backend/common/util/encrypt"
	mysession "app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	myuser "app/backend/model/mysql/user"
	"github.com/kataras/iris"
	mylog "app/backend/common/util/log"
	"strconv"
)

var log =  mylog.Log

type LoginController struct {
	*iris.Context
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Check username && password
func (lc *LoginController) check(name, password string) (*myuser.User, *myerror.YceError) {

	encryptPass := encrypt.NewEncryption(password).String()

	// New a User
	user := new(myuser.User)

	err := user.QueryUserByNameAndPassword(name, encryptPass)

	if err != nil {
		log.Errorf("Can not find the user: username=%s, err=%s", name, err)
		ye := myerror.NewYceError(1001, err.Error(), "")
		return nil, ye
	}

	return user, nil

}

// Store Session and Set cookie
func (lc *LoginController) session(user *myuser.User) (*mysession.Session, *myerror.YceError) {

	// Store (id,orgId) in SessionStore
	id := strconv.Itoa(int(user.Id))
	orgId := strconv.Itoa(int(user.OrgId))

	session := mysession.NewSession(id, user.Name, orgId)

	ss := mysession.SessionStoreInstance()

	err := ss.Set(session)

	log.Infof("Session: sessionId=%s, userId=%s, userName=%s, orgId=%s", session.SessionId, session.UserId, session.UserName, session.OrgId)

	if err != nil {
		log.Fatalf("Set session error: sessionId=%s, err=%s", session.SessionId, err)
		ye := myerror.NewYceError(2001, err.Error(), "")
		return nil, ye
	}

	return session, nil

}

// POST /api/v1/users/login
func (lc LoginController) Post() {

	loginParams := new(LoginParams)

	lc.ReadJSON(loginParams)

	user, ye := lc.check(loginParams.Username, loginParams.Password)
	if ye != nil {
		json, _ := ye.EncodeJson()
		lc.Write(json)
		return
	}

	session, ye := lc.session(user)
	if ye != nil {
		json, _ := ye.EncodeJson()
		lc.Write(json)
		return
	}

	// Set cookie
	// lc.SetCookieKV("sessionId", session.SessionId)

	// Auth pass
	sessionStr, _ := session.EncodeJson()

	ye = myerror.NewYceError(0, "OK", sessionStr)

	json, _ := ye.EncodeJson()

	log.Errorf("User Login: sessionId=%s, userId=%s, userName=%s, orgId=%s",
		session.SessionId, session.UserId, session.UserName, session.OrgId)

	lc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	lc.Write(json)

	return
}
