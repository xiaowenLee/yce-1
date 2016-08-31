package main

import (
	"k8s.io/kubernetes/pkg/api"
	// unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/util/intstr"
	"log"
	"os"
)

var logger *log.Logger

const (
	DEPLOYMENT           string = "nginx-deployment"
	RevisionHistoryLimit int32  = 5
	SERVER               string = "http://172.21.1.11:8080"
	RevisionAnnotation   string = "deployment.kubernetes.io/revision"
)

func init() {
	logger = log.New(os.Stdout, "", 0)
}

func rollingUpdate(c *client.Client, dp *extensions.Deployment) error {

	// New a DeploymentStrategy
	ds := new(extensions.DeploymentStrategy)
	ds.Type = extensions.RollingUpdateDeploymentStrategyType
	ds.RollingUpdate = new(extensions.RollingUpdateDeployment)
	ds.RollingUpdate.MaxUnavailable = intstr.FromInt(int(dp.Spec.Replicas))

	dp.Spec.Strategy = *ds

	// Image
	dp.Spec.Template.Spec.Containers[0].Image = "nginx:1.9.7"

	// Update
	_, err := c.Extensions().Deployments(api.NamespaceDefault).Update(dp)
	if err != nil {
		logger.Printf("Update Deployment Error: err=%s\n", err)
		return err
	}
	return nil
}

func main() {
	config := &restclient.Config{
		Host: SERVER,
	}

	c, err := client.New(config)
	if err != nil {
		logger.Fatalf("Could not connect to k8s api: err=%s\n", err)
	}

	dp, err := c.Extensions().Deployments(api.NamespaceDefault).Get(DEPLOYMENT)
	if err != nil {
		logger.Fatalf("Could not list deployments: err=%s\n", err)
	}

	rollingUpdate(c, dp)

}
