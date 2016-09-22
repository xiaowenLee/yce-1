package yce

import (
	"github.com/kataras/iris"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
)


type Controller struct {
	*iris.Context
	Ye *myerror.YceError
}

type IController interface {
	WriteError()
	ValidateSession(sessionId, orgId string)
	CheckError() bool // if error true, else false
	WriteOk()
}

func (c *Controller) WriteError() {
	c.Response.Header.Set("Access-Control-Allow-Origin", "*")
	log.Infof("Controller Response YceError: controller=%p, code=%d, msg=%s", c, c.Ye.Code, myerror.Errors[c.Ye.Code].LogMsg)
	c.Write(c.Ye.String())
}

func (c *Controller) ValidateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		c.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		c.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	log.Infof("Controller ValidateSession successfully")

	return
}

func (c *Controller) CheckError() bool {
	if c.Ye != nil {
		c.WriteError()
		return true
	}
	return false
}

func (c *Controller) WriteOk(msg string) {
	c.Response.Header.Set("Access-Control-Allow-Origin", "*")
	c.Ye = myerror.NewYceError(myerror.EOK, msg)
	log.Infof("Controller Response OK: controller=%p, code=%d, msg=%s", c, c.Ye.Code, myerror.Errors[c.Ye.Code].LogMsg)
	c.Write(c.Ye.String())
}
