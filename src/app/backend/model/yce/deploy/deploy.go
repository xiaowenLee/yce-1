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

type DcList struct {
	Dclist []string `json:"dcList"`
}


type Data struct {
	DcId int32 `json:"dcId"`
	DcName string      `json:"dcName"`
	PodList    api.PodList `json:"podList"`
}