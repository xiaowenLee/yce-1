package service

import (
	"k8s.io/kubernetes/pkg/api"
)

type Service struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
	ServiceList api.ServiceList `json:"serviceList"`
}
