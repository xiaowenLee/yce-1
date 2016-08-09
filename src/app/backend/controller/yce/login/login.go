package login

import (
	"log"
	"github.com/kataras/iris"
	"app/backend/common/util/encrypt"
	myuser "app/backend/model/mysql/user"
	mysession "app/backend/common/util/session"
	"strconv"
	"fmt"
)

type LoginController struct {
	*iris.Context
}


// POST /api/v1/users/{email}/login

func (lc LoginController) Post() {

	email := lc.Param("email")
	password := string(lc.FormValue("password"))

	encryptPass := encrypt.NewEncryption(password).String()

	// New a User
	user := new(myuser.User)

	err := user.QueryUserByNameAndPassword(email, encryptPass)

	if err != nil {
		log.Printf("Can not find the user: username=%s, err=%s\n", email, err)
		return
	}

	// Store (id,orgId) in SessionStore
	id := strconv.Itoa(int(user.Id))
	orgId := strconv.Itoa(int(user.OrgId))

	session := mysession.NewSession(id, user.Name, orgId)

	ss := mysession.SessionStoreInstance()

	err = ss.Set(session)

	fmt.Printf("Session: sessionId=%s, userId=%s, userName=%s, orgId=%s\n", session.SessionId, session.UserId, session.UserName, session.OrgId)

	if err != nil {
		log.Fatal("Set session error: sessionId=%s, err=%s\n", session.SessionId, err)
		return
	}

	// auth pass
	lc.Write("Hello world")
}