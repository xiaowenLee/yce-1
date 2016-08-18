package deploy

import (
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
)

//var instance *DescribeDeployController

type DescribeDeployController struct {
	cli *client.Client
}

func NewDescribeDeployController(server string) *DescribeDeployController {
	config := &restclient.Config{
		Host: server,
	}
	cli, err := client.New(config)
	if err != nil {
		log.Printf("Get DescribeDeployController error: error=%s\n", err)
	}

	instance := &DescribeDeployController{cli: cli}
	return instance
}

/*
func (dec *DescribeDeployController) Describe(ctx *iris.Context) {
	//TODO: ValidateSession
	oid := ctx.Param("oid")
	id := ctx.Param("id")

	//TODO: get Datacenter Host from MySQL ? or url ?
	//NOTE: Datacenter only takes one value when describing.
	//e.g.
	dc := make([]deploy.AppDc, 1)
	dc[0].DcID = 1

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

		poddetail, err := newCli.Pods(oid).Get(id)
		if err != nil {
			log.Printf("Get poddetails error. DataCenter=%s, Organization=%s, SessionID=%s, error=%s\n", v.DcID, oid, sessionID, err)
		}

		//TODO: make response poddetails struct
		//NOTE: time convertion, dc Chinese convertion
	}

	//TODO: write response json
}
*/
