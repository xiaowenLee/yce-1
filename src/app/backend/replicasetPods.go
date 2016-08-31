package main


import (
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"os"
)

var logger *log.Logger

const (
	DEPLOYMENT         string = "nginx-deployment"
	SERVER             string = "http://172.21.1.11:8080"
	RevisionAnnotation string = "deployment.kubernetes.io/revision"
)

func init() {
	logger = log.New(os.Stdout, "", 0)
}

func getPodsByReplicaSet(c *client.Client, rs *extensions.ReplicaSet) ([]api.Pod, error) {
	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}

	podList, err := c.Pods(api.NamespaceDefault).List(options)
	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}

func main() {

	config := &restclient.Config {
		Host: SERVER,
	}

	c, err := client.New(config)
	if err != nil {
		logger.Fatalf("Could not connect to k8s api: err=%s\n", err)
	}

	rsList, err := c.Extensions().ReplicaSets(api.NamespaceDefault).List(api.ListOptions{})
	if err != nil {
		logger.Fatalf("Could not list deployments: err=%s\n", err)
	}

	// Foreach replicaSet
	for _, rs := range rsList.Items {
		logger.Printf("ReplicaSet:\t%s\n", rs.Name)
		podList, err := getPodsByReplicaSet(c, &rs)
		if err != nil {
			logger.Fatalf("GetPodsByReplicaSet Error: err=%s\n", err)
		}

		for _, pod := range podList {
			logger.Printf("\tPodName:\t%s\n", pod.Name)
		}
	}
}