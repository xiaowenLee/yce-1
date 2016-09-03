package deploy

import (
	myqouta "app/backend/model/mysql/quota"
	mydatacenter "app/backend/model/mysql/datacenter"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	myorganization "app/backend/model/mysql/organization"
	mydeployment "app/backend/model/mysql/deployment"
)

// AppDeployment for Frontend fulfillment
type AppCreateDeployment struct {
	Datacenters []AppDc               `json:"dataCenters"`
	Deployment  extensions.Deployment `json:"deployment"`
}

type AppDc struct {
	DcID float64 `json:"dcID,omitempty"`
}

type AppDisplayDeployment struct {
	DcId    int32       `json:"dcId"`
	DcName  string      `json:"dcName"`
	PodList api.PodList `json:"podList"`
}

type DeployAndPodList struct {
	Deploy *extensions.Deployment `json:"deploy"`
	PodList *api.PodList `json:"podList"`
}


// Response List Deployments
type Deployment struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
	Deployments []DeployAndPodList `json:"deployments"`
}

// GET
type ListDeployment struct {
	Organization *myorganization.Organization
	DcIdList []int32
	DcName []string
}

// Get .../init
type InitDeployment struct {
	OrgId string `json:"orgId"`
	OrgName string `json:"orgName"`
	DataCenters []mydatacenter.DataCenter `json:"dataCenters"`
	Quotas []myqouta.Quota `json:"quotas"`
}

// DcIdList
type DcIdListType struct {
	DcIdList []int32 `json:"dcIdList"`
}

//TODO: Change DcIdList 2 []int32
// Post .../new
type CreateDeployment struct {
	AppName  string `json:"appName"`
	OrgName  string `json:"orgName"`
	DcIdList []int32 `json:"dcIdList"`
	//DcIdList DcIdListType `json:"dcIdList"`
	Deployment extensions.Deployment `json:"deployment"`
}

// RollingUpdate Strategy
type RollingStrategy struct {
	MaxUnavailable int32 `json:"maxUnavailable"`
	Image string `json:"image"`
	UpdateInterval int32 `json:"updateInterval"`
}

//TODO: Change DcIdList 2 []int32
// RollingUpdate Deployment
type RollingDeployment struct {
	AppName string `json:"appName"`
	OrgName string `json:"orgName"`
	//DcIdList DcIdListType `json:"dcIdList"`
	DcIdList []int32 `json:"dcIdList"`
	UserId int32 `json:"userId"`
	Strategy RollingStrategy `json:"strategy"`
	Comments string `json:"comments"`
}

// Response List OperationLog

/*
type DcIdListType struct {
	DcIdList []int32 `json:"dcIdList"`
}
*/

type OperationLog struct {
	DcName []string `json:"dcName"`
	UserName string `json:"userName"`
	Record *mydeployment.Deployment `json:"records"`
}


// Get .../history
type ListOperationLog struct {
	DcName []string
	DcIdList []int32
	Organization *myorganization.Organization
}