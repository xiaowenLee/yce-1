package topology

import (
	"k8s.io/kubernetes/pkg/api"
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

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
