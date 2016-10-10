package deploymentstat


import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

type DeploymentsType struct {
	DeploymentName string `json:"deploymentName"`
	RsName string `json:"rsName"`
	PodName []string `json:"podName"`
}


type DatacentersType struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
	Deployments []DeploymentsType `json:"deployments"`
}


type DeploymentStatType struct {
	DeploymentStat []DatacentersType `json:"deploymentStat"`
}

