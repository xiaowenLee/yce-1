package deploy

import (
	myqouta "app/backend/model/mysql/quota"
	mydatacenter "app/backend/model/mysql/datacenter"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	myorganization "app/backend/model/mysql/organization"
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

// Post .../new
type CreateDeployment struct {
	AppName  string `json: "appName"`
	OrgName  string `json: "orgName"`
	DcIdList []int32 `json:"dcIdList"`
	Deployment extensions.Deployment `json:"deployment"`
}
