package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"os"
)

/*==========================================================================
 Constant
 logger
 topology
==========================================================================*/
const (
	SERVER    string = "http://172.21.1.11:8080"
	NAMESPACE string = "ops"
)

var logger *log.Logger

var topology *Topology

/*==========================================================================
 Definations
==========================================================================*/
type PodType struct {
	api.Pod
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

type ServiceType struct {
	api.Service
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

type ReplicaSetType struct {
	extensions.ReplicaSet
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

type NodeType struct {
	api.Node
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
}

type ItemType map[string]interface{}

type RelationsType struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Topology struct {
	Items     ItemType        `json:"items"`
	Relations []RelationsType `json:"relations"`
}

/*==========================================================================
 Helper funcs
==========================================================================*/

func getDeploymentsByNamespace(c *client.Client, namespace string) ([]extensions.Deployment, error) {
	dps, err := c.Extensions().Deployments(namespace).List(api.ListOptions{})
	if err != nil {
		logger.Fatalf("getDeploymentsByNamespace Error: err=%s\n", err)
		return nil, err
	}

	return dps.Items, nil
}

func getReplicaSetsByDeployment(c *client.Client, deployment *extensions.Deployment) ([]extensions.ReplicaSet, error) {

	namespace := deployment.Namespace
	selector, err := unver.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		logger.Fatalf("getReplicaSetsByDeployment Error: err=%s\n", err)
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}
	rsList, err := c.Extensions().ReplicaSets(namespace).List(options)

	logger.Printf("getReplicaSetsByDeployment: dp.Name=%s, len(rs.Items)=%d\n", deployment.Name, len(rsList.Items))

	return rsList.Items, nil
}

func getPodsByReplicaSet(c *client.Client, namespace string, rs *extensions.ReplicaSet) ([]api.Pod, error) {
	selector, err := unver.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		logger.Fatalf("getPodsByReplicaSet Error: err=%s\n", err)
		return nil, err
	}
	options := api.ListOptions{LabelSelector: selector}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		logger.Fatalf("getPodsByReplicaSet Error: err=%s\n", err)
		return nil, err
	}

	logger.Printf("getPodsByReplicaSet: rs.Name=%s, len(rs.Items)=%d\n", rs.Name, len(podList.Items))

	return podList.Items, nil
}

func getNodeByPod(c *client.Client, pod *api.Pod) (*api.Node, error) {
	name := pod.Spec.NodeName
	node, err := c.Nodes().Get(name)
	if err != nil {
		logger.Fatalf("getNodeByPod Error: err=%s\n", err)
		return nil, err
	}
	return node, nil
}

func getServicesByNamespace(c *client.Client, namespace string) ([]api.Service, error) {
	svcs, err := c.Services(namespace).List(api.ListOptions{})
	if err != nil {
		logger.Fatalf("getServicesByNamespace Error: err=%s\n", err)
		return nil, err
	}

	return svcs.Items, nil
}

func getPodByService(c *client.Client, namespace string, svc *api.Service) ([]api.Pod, error) {
	selector := new(unver.LabelSelector)
	selector.MatchLabels = svc.Spec.Selector

	s, err := unver.LabelSelectorAsSelector(selector)
	if err != nil {
		logger.Fatalf("getPodByService Error: err=%s\n", err)
		return nil, err
	}

	options := api.ListOptions{LabelSelector: s}

	podList, err := c.Pods(namespace).List(options)
	if err != nil {
		logger.Fatalf("getPodByService Error: err=%s\n", err)
		return nil, err
	}

	return podList.Items, nil
}

/*==========================================================================
 Topology
==========================================================================*/
/*

begin:
	ops --> Deployments.List
	Foreach deployment in Deployments.List
		rs := findNewReplicaSet()
		rs --> Select Pods.List
		Foreach pod in Pods.List(){
			Pod.Name --> Node: pod <---> node
			rs <---> pod
		}
	ops --> Services.List
	Foreach service in Services.List
		service --> Select Pods.List
		service <---> pod
:end
*/
func getTopology(c *client.Client, namespace string) bool {
	// Get Deployments.List
	dpList, err := getDeploymentsByNamespace(c, namespace)
	if err != nil {
		return false
	}

	// Foreach Deployments.List
	for _, dp := range dpList {
		rsList, err := getReplicaSetsByDeployment(c, &dp)
		if err != nil {
			logger.Fatalf("getTopology Error: err=%s\n", err)
			return false
		}

		for _, rs := range rsList {
			// Topology.Items
			uid := string(rs.UID)
			topology.Items[uid] = ReplicaSetType{
				Kind:       "ReplicaSet",
				ApiVersion: "v1beta2",
				ReplicaSet: rs,
			}

			podList, err := getPodsByReplicaSet(c, namespace, &rs)
			if err != nil {
				return false
			}
			for _, pod := range podList {
				uid = string(pod.UID)
				topology.Items[uid] = PodType{
					Kind:       "Pod",
					ApiVersion: "v1beta2",
					Pod:        pod,
				}

				relation := RelationsType{
					Source: string(rs.UID),
					Target: string(pod.UID),
				}

				topology.Relations = append(topology.Relations, relation)

				node, err := getNodeByPod(c, &pod)
				if err != nil {
					return false
				}

				uid = string(node.UID)
				topology.Items[uid] = NodeType{
					Kind:       "Node",
					ApiVersion: "v1beata2",
					Node:       *node,
				}

				relation = RelationsType{
					Source: string(node.UID),
					Target: string(pod.UID),
				}
				topology.Relations = append(topology.Relations, relation)
			}
		}
	}

	// Get Services.List
	svcList, err := getServicesByNamespace(c, namespace)
	if err != nil {
		logger.Fatalf("getTopology Error: err=%s\n", err)
		return false
	}

	for _, svc := range svcList {
		uid := string(svc.UID)
		topology.Items[uid] = ServiceType{
			Kind:       "Service",
			ApiVersion: "v1beta1",
			Service:    svc,
		}

		podList, err := getPodByService(c, namespace, &svc)
		if err != nil {
			logger.Fatalf("getTopology Error: err=%s\n", err)
			return false
		}

		for _, pod := range podList {
			uid = string(pod.UID)
			if _, ok := topology.Items[uid]; ok {
				topology.Items[uid] = PodType{
					Kind:       "Pod",
					ApiVersion: "v1beta1",
					Pod:        pod,
				}
			}

			relation := RelationsType{
				Source: string(svc.UID),
				Target: string(pod.UID),
			}

			topology.Relations = append(topology.Relations, relation)
		}
	}

	return true
}

func encodeTopology() (string, error) {

	data, err := json.MarshalIndent(topology, "", "\t")
	if err != nil {
		logger.Fatalf("encodeTopology Error: err=%s\n", err)
		return "", err
	}
	return string(data), nil
}

/*==========================================================================
==========================================================================*/
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
	}

	topology = new(Topology)
	topology.Items = make(ItemType)
	topology.Relations = make([]RelationsType, 0)

	if ok := getTopology(c, NAMESPACE); ok {
		str, err := encodeTopology()
		if err != nil {
			logger.Fatalf("Error: err=%s\n", err)
		}
		fmt.Println(str)
	}
}
