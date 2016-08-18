package deploy

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

// AppDeployment for Frontend fulfillment
type AppDeployment struct {
	Datacenters []AppDc               `json:"dataCenters"`
	Deployment  extensions.Deployment `json:"deployment"`
}

type AppDc struct {
	DcID float64 `json:"dcID,omitempty"`
}


type Data struct {
	DataCenter string      `json:"dataCenter"`
	PodList    api.PodList `json:"podList"`
}