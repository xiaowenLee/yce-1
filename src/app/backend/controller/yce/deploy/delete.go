package deploy

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
)

type DeleteDeployController struct {
	cli *client.Client
}

func NewDeleteDeployController(server string) *DeleteDeployController {
	config := &restclient.Config{
		Host: server,
	}
	cli, err := client.New(config)
	if err != nil {
		log.Printf("Get DeleteDeployController error: error=%s\n", err)
	}

	instance := &DeleteDeployController{cli: cli}
	return instance
}

func (dc *DeleteDeployController) Delete(c *client.Client) error {
	err := c.Extensions().Deployments(api.NamespaceDefault).Delete("nginx-deploy-app", &api.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
