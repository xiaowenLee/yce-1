package main

import (
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"os"
	"time"
	"strings"
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

func encodeMapToString(labels map[string]string) string {
	var ss []string
	for key, value := range labels {
		str := key + ":" + value
		ss = append(ss, str)
	}

	return strings.Join(ss, ",")
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

func displayReplicaSet(c *client.Client, dp *extensions.Deployment) {
	rses, err := getReplicaSetsByDeployment(c, dp)
	if err != nil {
		logger.Fatalf("Could not getReplicaSetsByDeployment: err=%s\n", err)
	}

	for _, rs := range rses {
		name := rs.Name
		namespace := rs.Namespace
		desireReplicas := rs.Spec.Replicas
		selector := encodeMapToString(rs.Spec.Selector.MatchLabels)
		image := rs.Spec.Template.Spec.Containers[0].Image
		actualReplicas := rs.Status.Replicas
		revision := rs.Annotations[RevisionAnnotation]

		logger.Printf("Deployment Revision: %s\n", revision)
		logger.Printf("\tName: %s\n", name)
		logger.Printf("\tNamespace: %s\n", namespace)
		logger.Printf("\tImage: %s\n", image)
		logger.Printf("\tSelector: %s\n", selector)
		logger.Printf("\tReplicas: %d current / %d desired\n\n", actualReplicas, desireReplicas)
	}
}

func displayDeployment(dp *extensions.Deployment) {

	// Display Deployment
	logger.Println("Deployment:\t")
	name := dp.Name
	logger.Printf("Name:\t\t%s\n", name)

	namespace := dp.Namespace
	logger.Printf("Namespace:\t\t%s\n", namespace)

	creationTimestamp := dp.CreationTimestamp.Time.Format(time.RFC1123Z)
	logger.Printf("CreationTimestamp:\t\t%s\n", creationTimestamp)

	selector := encodeMapToString(dp.Spec.Selector.MatchLabels)
	logger.Printf("Selector:\t\t%s\n", selector)

	labels := encodeMapToString(dp.Labels)
	logger.Printf("Labels:\t\t%s\n", labels)


	// Annotations
	annotations := encodeMapToString(dp.Annotations)
	logger.Printf("Annotations:\t\t%s\n", annotations)

	// Replicas
	updatedReplicas := dp.Status.UpdatedReplicas
	totalReplicas := dp.Spec.Replicas
	availableReplicas := dp.Status.AvailableReplicas
	unavailableReplicas := dp.Status.UnavailableReplicas
	logger.Printf("Replicas:\t\t%d updated | %d total | %d available | %d unavailable\n", updatedReplicas, totalReplicas, availableReplicas, unavailableReplicas)

	strategyType := dp.Spec.Strategy.Type
	logger.Printf("StrategyType:\t\t%s\n", strategyType)

	if dp.Spec.Strategy.RollingUpdate != nil {
		ru := dp.Spec.Strategy.RollingUpdate
		logger.Printf("RollingUpdateStrategy:\t\t%s max unavailable, %s max surge\n", ru.MaxUnavailable.String(), ru.MaxSurge.String())
	}

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

	dp, err := c.Extensions().Deployments(api.NamespaceDefault).Get(DEPLOYMENT)
	if err != nil {
		logger.Fatalf("Could not list deployment: err=%s\n", err)
		return
	}

	displayDeployment(dp)

	logger.Printf("\n\n")

	displayReplicaSet(c, dp)

}
