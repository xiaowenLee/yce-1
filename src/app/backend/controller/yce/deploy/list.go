package deploy

import (
	hc "app/backend/common/util/http/httpclient"
	session "app/backend/common/util/session"
	deploy "app/backend/model/yce/deploy"
	"fmt"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"strings"
)

type ListDeployController struct {
	cli  *client.Client
	iris *iris.Context
}

func NewListDeployController(server string) *ListDeployController {
	config := &restclient.Config{
		Host: server,
	}
	cli, err := client.New(config)
	if err != nil {
		log.Printf("Get ListDeployController error. SessionID=%s, error=%s\n", sessionID, err)
	}

	instance = &ListDeployController{cli: cli}
	return instance
}

func validateSession(client, uid) (ok bool, err error) {
	//sessionIdfromClient := ctx.RequestHeader("sessionId")
	//get sessionId from Redis refer to uid
	sessionId := Redis.Get(uid)
}

func (lc *ListDeployController) getDcHost() {
	//TODO: get Datacenter Host from MySQL
	//e.g.
	dc := make([]deploy.AppDc, 1)
	dc[0].DcID = 1
}

func (lc *ListDeployController) getPodList() {

	var Server string
	for _, v := range dc {
		switch v.DcID {
		case 1:
			Server = "http://172.21.1.11:8080"
		case 2:
			Server = "http://172.21.1.11:8080"
		case 3:
			Server = "http://172.21.1.11:8080"
		}

		newconfig := &restclient.Config{
			Host: Server,
		}
		newCli, err := client.New(newconfig)
		if err != nil {
			log.Printf("Get new restclient error. SessionID=%s, error=%s\n", sessionID, err)
		}

		podlist, err := newCli.Pods(oid).List(api.ListOptions{})
		if err != nil {
			log.Printf("Get podlist error. DataCenter=%s, Organization=%s, SessionID=%s, error=%s\n", v.DcID, oid, sessionID, err)
		}

		//TODO: make response podlist struct
		//NOTE: time convertion, dc Chinese convertion
	}
}

func (lc *ListDeployController) List() {

	sessionIdfromClient := ctx.RequestHeader("sessionId")
	oid := ctx.Param("oid")
	if ok, err := validateSession(sessionIdFromClient, uid); ok {
		getDcHost()
		getPodList()
	} else {
		log.Printf("Validate Session error. sessionID=%s, error=%s\n", sessionID, err)
	}

	//TODO: write response json
}
