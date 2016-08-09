package login

import (
	"log"
	"github.com/kataras/iris"
	"app/backend/common/util/encrypt"
	"app/backend/model/mysql/user"
	"app/backend/common/util/session"
)

type LoginController struct {
	*iris.Context
}


// POST /api/v1/users/{email}/login

func (lc *LoginController) Login() {

	// Parse username(email) and password
	email := lc.FormValue("username")
	password := lc.FormValue("password")

	encryptPass := encrypt.NewEncryption(password)

	// New a User
	user := new(User)

	err := user.QueryUserByNameAndPassword(email, encryptPass)

	if err != nil {
		log.Printf("Can not find the user: username=%s, err=%s\n", email, err)
		return
	}

	// Store (id,orgId) in SessionStore
	session := NewSession(user.Id, user.Name, user.OrgId)

	ss := session.SessionStoreInstance()

	err = ss.Set(session)

	if err != nil {
		log.Fatal("Set session error: sessionId=%s, err=%s\n", session.SessionId, err)
		return
	}

	// auth pass
}