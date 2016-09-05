package main

import (
	//"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"log"
	"os"
	// unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

var logger *log.Logger

var annotations = map[string]string{
	"Image":                      "nginx:1.7.9",
	"UserId":                     "2",
	"kubernetes.io/change-cause": "版本不匹配1",
}

const (
	DEPLOYMENT         string = "nginx-test"
	REVERSION          int64  = 4
	SERVER             string = "http://172.21.1.11:8080"
	RevisionAnnotation        = "deployment.kubernetes.io/revision"
)

func init() {
	logger = log.New(os.Stdout, "", 0)
}

func rollBack(c *client.Client, dp *extensions.Deployment, revision int64) error {

	dr := new(extensions.DeploymentRollback)
	dr.Name = dp.Name
	dr.UpdatedAnnotations = annotations
	dr.RollbackTo = extensions.RollbackConfig{Revision: revision}

	// Rollback
	err := c.Extensions().Deployments("ops").Rollback(dr)
	if err != nil {
		logger.Printf("Deployment Rollback Error: err=%s\n", err)
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

	dp, err := c.Extensions().Deployments("ops").Get(DEPLOYMENT)
	if err != nil {
		logger.Fatalf("Could not list deployments: err=%s\n", err)
	}

	rollBack(c, dp, REVERSION)

}
