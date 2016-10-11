package main

import (
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"os"
	"strings"
)

var logger *log.Logger

const (
	SERVER             string = "http://172.21.1.11:8080"
	RevisionAnnotation        = "deployment.kubernetes.io/revision"
)

func init() {
	logger = log.New(os.Stdout, "", 0)
}

func getReplicaSetsByDeployment(c *client.Client, deployment *extensions.Deployment) ([]extensions.ReplicaSet, error) {

	namespace := deployment.Namespace
	selector, err := unver.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}
	rsList, err := c.Extensions().ReplicaSets(namespace).List(options)

	return rsList.Items, nil
}

func getDeploymentByReplicaSet(namespace string, c *client.Client, rs *extensions.ReplicaSet) ([]extensions.Deployment, error) {

	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}
	dps, err := c.Extensions().Deployments(namespace).List(options)
	if err != nil {
		return nil, err
	}

	return dps.Items, nil
}

func getDeploymentByReplicaSetName(namespace string, c *client.Client, rs *extensions.ReplicaSet) (*extensions.Deployment, error) {

	name := rs.Name
	index := strings.LastIndex(name, "-")
	deploymentName := name[:index]

	dp, err := c.Extensions().Deployments(namespace).Get(deploymentName)
	if err != nil {
		return nil, err
	}

	return dp, nil
}

func main() {

	config := &restclient.Config{
		Host: SERVER,
	}

	c, err := client.New(config)
	if err != nil {
		logger.Fatalf("Could not connect to k8s api: err=%s\n", err)
	}

	list, err := c.Extensions().Deployments(api.NamespaceDefault).List(api.ListOptions{})
	if err != nil {
		logger.Fatalf("Could not list deployments: err=%s\n", err)
	}

	logger.Printf("Deployment -------> ReplicaSet: ")

	for _, deployment := range list.Items {
		rses, err := getReplicaSetsByDeployment(c, &deployment)
		if err != nil {
			logger.Fatalf("GetReplicaSetsByDeployment Error: err=%s\n", err)
		}

		for _, rs := range rses {
			logger.Printf("ReplicaSet assioated with Deployment: rs-name=%s, revision=%s, dp-name=%s\n",
				rs.Name, rs.Annotations[RevisionAnnotation], deployment.Name)
		}
	}

	logger.Printf("\n\nReplicaSet -------> Deployment: ")
	rsList, err := c.Extensions().ReplicaSets(api.NamespaceDefault).List(api.ListOptions{})
	if err != nil {
		log.Fatalf("Could not list ReplicaSet")
	}

	for _, rs := range rsList.Items {
		dp, err := getDeploymentByReplicaSetName(api.NamespaceDefault, c, &rs)
		if err != nil {
			logger.Fatalf("GetDeploymentByReplicaSet Error: err=%s\n", err)
		}

		logger.Printf("Deployment assioated with ReplicaSet: rs-name=%s, revision=%s, dp-name=%s\n",
			rs.Name, rs.Annotations[RevisionAnnotation], dp.Name)
	}
}
