package deploy

import (
	hc "app/backend/common/util/http/httpclient"
	session "app/backend/common/util/session"
	organization "app/backend/model/mysql/organization"
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
		log.Printf("Get ListDeployController error: SessionId=%s, error=%s\n", sessionId, err)
	}

	instance = &ListDeployController{cli: cli}
	return instance
}

func (lc *ListDeployController) getDcHost(orgId string) ([]string, error) {
	//TODO: get Datacenter Host from MySQL

	// example below
	server := make([]string, 1)
	server[0] = "http://172.21.1.11:8080"

	return server, nil
}

func (lc *ListDeployController) getPodList(server []string, orgId string) (api.PodList, error) {

	for _, v := range server {
		newconfig := &restclient.Config{
			Host: Server,
		}
		newCli, err := client.New(newconfig)
		if err != nil {
			log.Printf("Get new restclient error: sessionId=%s, error=%s\n", sessionId, err)
			return nil, err
		}

		podlist, err := newCli.Pods(orgId).List(api.ListOptions{})
		if err != nil {
			log.Printf("Get podlist error: server=%s, orgId=%s, error=%s\n", v, orgId, err)
			return nil, err
		}
		return podlist, nil
	}
}

func (lc ListDeployController) Get() {

	sessionIdClient := ctx.RequestHeader("sessionId")
	orgId := ctx.Param("orgId")
	userId := ctx.Param("uid")
	if ok, err := session.ValidateUserId(sessionIdClient, userId); ok {
		server, err := lc.getDcHost(orgId)
		if err != nil {
			log.Printf("Get Datacenter Host error: sessionId=%s, orgId=%s, err=%s\n", sessionIdClient, orgId, err)
		}

		podlist, err := lc.getPodList(server, orgId)
		if err != nil {
			log.Printf("Get Podlist error: sessionId=%s, orgId=%s, error=%s\n", sessionIdClient, orgId, err)
		}

		//TODO: write response json
		deployList := make()

	} else {
		log.Printf("Validate Session error: sessionId=%s, error=%s\n", sessionIdClient, err)
	}

}
