package endpoint

import (
	mydatacenter "app/backend/model/mysql/datacenter"
	mynodeport "app/backend/model/mysql/nodeport"
	myorganization "app/backend/model/mysql/organization"
	"k8s.io/kubernetes/pkg/api"
)

type Endpoints struct {
	DcId          int32             `json:"dcId"`
	DcName        string            `json:"dcName"`
	EndpointsList api.EndpointsList `json:"endpointsList`
}

type ListEndpoints struct {
	Organization *myorganization.Organization
	DcIdList     []int32
	DcName       []string
}

type InitEndpoints struct {
	OrgId       string                    `json:"orgId"`
	OrgName     string                    `json:"orgName"`
	DataCenters []mydatacenter.DataCenter `json:"dataCenters"`
	NodePort    *mynodeport.NodePort      `json:"nodePort"`
}

type CreateEndpoints struct {
	EndpointsName string        `json:"endpointsName`
	OrgName       string        `json:"orgName"`
	DcIdList      []int32       `json:"dcIdList"`
	Endpoints     api.Endpoints `json:"endpoints"`
}
