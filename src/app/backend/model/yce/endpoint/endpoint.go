package endpoint

import (
	"k8s.io/kubernetes/pkg/api"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
)

type Endpoints struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
	EndpointsList api.EndpointsList `json:"endpointsList`
}

type ListEndpoints struct {
	Organization *myorganization.Organization
	DcIdList []int32
	DcName []string
}
