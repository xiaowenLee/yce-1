package login

import (
	"log"
	"strconv"
	"github.com/kataras/iris"
	"app/backend/common/util/encrypt"
	myuser "app/backend/model/mysql/user"
	mysession "app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
)

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
		log.Printf("Can not find the user: username=%s, err=%s\n", name, err)
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

	log.Printf("Session: sessionId=%s, userId=%s, userName=%s, orgId=%s\n", session.SessionId, session.UserId, session.UserName, session.OrgId)

	if err != nil {
		log.Fatal("Set session error: sessionId=%s, err=%s\n", session.SessionId, err)
		ye := myerror.NewYceError(2001, err.Error(), "")
		return nil, ye
	}

	return session, nil

}

// POST /api/v1/users/login
func (lc LoginController) Post() {

	loginParams := new(LoginParams)

	lc.ReadJSON(loginParams)

	log.Printf("LoginParam: %v\n", loginParams)

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
	sessionStr, _:= session.EncodeJson()

	ye = myerror.NewYceError(0, "OK", sessionStr)

	json, _ := ye.EncodeJson()

	log.Printf("User Login: sessionId=%s, userId=%s, userName=%s, orgId=%s\n",
		session.SessionId, session.UserId, session.UserName, session.OrgId)

	lc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	lc.Write(json)

	return
}
