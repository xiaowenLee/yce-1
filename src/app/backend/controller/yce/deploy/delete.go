package deploy

import (
	hc "app/backend/common/util/http/httpclient"
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

const (
	SERVER string = "http://172.21.1.11:8080"
)

var instance *DeleteDeployController

type DeleteDeployController struct {
	cli *client.Client
}

func NewDeleteDeployController(server string) *DeleteDeployController {
	config := &restclient.Config{
		Host: server,
	}
	cli, err := client.New(config)
	if err != nil {
		log.Printf("Get DeleteDeployController error. SessionID=%s, error=%s\n", sessionID, err)
	}

	instance = &DeleteDeployController{cli: cli}
	return instance
}

func (dc *DeleteDeployController) Delete(c *client.Client) error {
	err := c.Extensions().Deployments(api.NamespaceDefault).Delete("nginx-deploy-app", &api.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
