package deploy

import (
	hc "app/backend/common/util/http/httpclient"
	dp "app/backend/model/yce/deploy"
	"fmt"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"strings"
)

const (
	SERVER string = "http://172.21.1.11:8080"
)

type CreateDeployController struct {
	cli *client.Client
}

func NewCreateDeployController(server string) *CreateDeployController {
	config := &restclient.Config{
		Host: server,
	}
	cli, err := client.New(config)
	if err != nil {
		log.Printf("Get CreateDeployController error. SessionID=%s, error=%s\n", sessionID, err)
	}

	instance = &CreateDeployController{cli: cli}
	return instance
}

func (cc *CreateDeployController) Create(ctx *iris.Context) {
	//TODO: ValidateSession
	//TODO: unmarshal resquest json
	//e.g.
	myAppDeploy := new(dp.AppDeployment)
	err := ctx.ReadJSON(myAppDeploy)
	if err != nil {
		log.Printf("Read JSON error. SessionID=%s error=%s\n", sessionID, err)
	}
	//TODO: get OrgID, and DcHost refer to UserID
	//e.g.
	oid := api.NamespaceDefault
	dclen := len(myAppDeploy.Datacenters)
	dc := make([]deploy.AppDc, dclen)
	dc = myAppDeploy.Datacenters

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
			log.Printf("Create restclient error. SessionID=%s, error=%s\n", sessionID, err)
		}
		deployment := new(extensions.Deployment)
		deployment = &myAppDeploy.Deployment

		_, err = newCli.Extensions().Deployments(oid).Create(deployment)
		if err != nil {
			log.Printf("Create deployment error. Datacenter=%s, Organization=%s, error=%s\n", v.DcID, OrgID, err)
		}
		//TODO: decode create response status
		//TODO: make response json
	}
	//TODO: according to create status write to MySQL deploy log
	//TODO: write back response json
}
