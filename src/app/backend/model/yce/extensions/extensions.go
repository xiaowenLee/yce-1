package extensions

import (
	"k8s.io/kubernetes/pkg/api"
	myorganization "app/backend/model/mysql/organization"
)

type Extensions struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
	ServiceList api.ServiceList `json:"serviceList"`
	EndpointList api.EndpointsList `json:"endpointsList"`
}

type ListExtensions struct {
	Organization *myorganization.Organization
	DcIdList []int32
	DcName []string
}
