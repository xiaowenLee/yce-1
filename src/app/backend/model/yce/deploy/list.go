package deploy

import (
	"k8s.io/kubernetes/pkg/api"
)

/*
type DeployList struct {
	Code    float64          `json:"code"`
	Message []string         `json:"message"`
	Data    []DeployListData `json:"data"`
}

type DeployListData struct {
	DcId    string      `json:"dcId"`
	PodList api.PodList `json:"podList"`
}
*/

type Data struct {
	DataCenter string      `json:"dataCenter"`
	PodList    api.PodList `json:"podList"`
}
