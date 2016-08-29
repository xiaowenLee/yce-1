package service

import (
	"k8s.io/kubernetes/pkg/api"
	mydatacenter "app/backend/model/mysql/datacenter"
)

type Service struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
	ServiceList api.ServiceList `json:"serviceList"`
}

type InitService struct {
	OrgId string `json:"orgId"`
	OrgName string `json:"orgName"`
	DataCenters []mydatacenter.DataCenter `json:"dataCenters"`
}

type CreateService struct {
	ServiceName string `json:"serviceName"`
	OrgName string `json:"orgName"`
	DcIdList []int32  `json:"dcIdList"`
	Service api.Service `json:"service"`
}
