package deploy

import (
	"k8s.io/kubernetes/pkg/api"
)

type AppList struct {
	code    float64                `json:"code"`
	message []string               `json:"message"`
	data    map[string]api.PodList `json:"data"`
}
