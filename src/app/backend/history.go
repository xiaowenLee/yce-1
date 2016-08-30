package main

import (
	"k8s.io/kubernetes/pkg/api"
	// "k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"os"
)

var logger *log.Logger

const (
	SERVER string = "http://204.11.99.12:8080"
)


func init() {
	logger = log.New(os.Stdout, "", 0)
}

func main() {

	config := &restclient.Config{
		Host: SERVER,
	}

	c, err := client.New(config)

	if err != nil {
		logger.Fatalf("Could not connect to k8s api: err=%s\n", err)
		return
	}

	list, err := c.Extensions().Deployments(api.NamespaceDefault).List(api.ListOptions{})
	if err != nil {
		logger.Fatalf("Could not list deployment: err=%s\n", err)
		return
	}

	// Foreach item in the List
	for index, item := range list.Items {
		name := item.ObjectMeta.Name

		revisionHistoryLimit := item.Spec.RevisionHistoryLimit

		selector := item.Spec.Selector

		rollbackTo := item.Spec.RollbackTo

		logger.Printf("Deploymenet %s: index=%d, revisionHistoryLimit=%v, revision=%d, selector=%s\n",
			name, index, *revisionHistoryLimit, rollbackTo, selector)

	}
}
