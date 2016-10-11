package service

import (
	mydatacenter "app/backend/model/mysql/datacenter"
	mynodeport "app/backend/model/mysql/nodeport"
	myorganization "app/backend/model/mysql/organization"
	"k8s.io/kubernetes/pkg/api"
)

type Service struct {
	DcId        int32           `json:"dcId"`
	DcName      string          `json:"dcName"`
	ServiceList api.ServiceList `json:"serviceList"`
}

type ListService struct {
	Organization *myorganization.Organization
	DcIdList     []int32
	DcName       []string
}

type InitService struct {
	OrgId       string                    `json:"orgId"`
	OrgName     string                    `json:"orgName"`
	DataCenters []mydatacenter.DataCenter `json:"dataCenters"`
	NodePort    *mynodeport.NodePort      `json:"nodePort"`
}

type CreateService struct {
	ServiceName string      `json:"serviceName"`
	OrgName     string      `json:"orgName"`
	DcIdList    []int32     `json:"dcIdList"`
	Service     api.Service `json:"service"`
}
