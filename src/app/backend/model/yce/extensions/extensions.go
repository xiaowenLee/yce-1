package extensions

import (
	"k8s.io/kubernetes/pkg/api"
	myorganization "app/backend/model/mysql/organization"
	//myservice "app/backend/model/yce/service"
	//myendpoint "app/backend/model/yce/endpoint"
)

type Extensions struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
	ServiceList api.ServiceList `json:"serviceList"`
	EndpointList api.EndpointsList `json:"endpointsList"`
	//ServiceList myservice.Service `json:"serviceList`
	//EndpointsList myendpoint.Endpoints `json:"endpointsList"`

}

type ListExtensions struct {
	Organization *myorganization.Organization
	DcIdList []int32
	DcName []string
}
