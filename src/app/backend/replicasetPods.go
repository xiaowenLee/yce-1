package main

import (
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	deploymentutil "k8s.io/kubernetes/pkg/controller/deployment/util"
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

/*
// FindNewReplicaSet returns the new RS this given deployment targets (the one with the same pod template).
func FindNewReplicaSet(deployment *extensions.Deployment, rsList []extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	newRSTemplate := GetNewReplicaSetTemplate(deployment)
	for i := range rsList {
		equal, err := equalIgnoreHash(rsList[i].Spec.Template, newRSTemplate)
		if err != nil {
			return nil, err
		}
		if equal {
			// This is the new ReplicaSet.
			return &rsList[i], nil
		}
	}
	// new ReplicaSet does not exist.
	return nil, nil
}

// GetNewReplicaSetTemplate returns the desired PodTemplateSpec for the new ReplicaSet corresponding to the given ReplicaSet.
func GetNewReplicaSetTemplate(deployment *extensions.Deployment) api.PodTemplateSpec {
	// newRS will have the same template as in deployment spec, plus a unique label in some cases.
	newRSTemplate := api.PodTemplateSpec{
		ObjectMeta: deployment.Spec.Template.ObjectMeta,
		Spec:       deployment.Spec.Template.Spec,
	}
	newRSTemplate.ObjectMeta.Labels = labelsutil.CloneAndAddLabel(
		deployment.Spec.Template.ObjectMeta.Labels,
		extensions.DefaultDeploymentUniqueLabelKey,
		podutil.GetPodTemplateSpecHash(newRSTemplate))
	return newRSTemplate
}
*/

func main() {

	config := &restclient.Config{
		Host: SERVER,
	}

	c, err := client.New(config)
	if err != nil {
		logger.Fatalf("Could not connect to k8s api: err=%s\n", err)
	}

	// Deployment
	dp, err := c.Extensions().Deployments(api.NamespaceDefault).Get(DEPLOYMENT)
	if err != nil {
		logger.Fatalf("Could not list deployment: err=%s\n", err)
		return
	}

	// Get rs by the deployment
	rsList, err := getReplicaSetsByDeployment(c, dp)
	if err != nil {
		logger.Fatalf("Could not list deployment: err=%s\n", err)
		return
	}

	// Find the newest rs
	rs, err := deploymentutil.FindNewReplicaSet(dp, rsList)
	if err != nil {
		logger.Fatalf("Could not list deployment: err=%s\n", err)
		return
	}

	// Pods
	logger.Printf("ReplicaSet:\t%s\n", rs.Name)
	podList, err := getPodsByReplicaSet(c, rs)
	if err != nil {
		logger.Fatalf("GetPodsByReplicaSet Error: err=%s\n", err)
	}

	for _, pod := range podList {
		logger.Printf("\tPodName:\t%s\n", pod.Name)
	}

	/*
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
	*/
}
