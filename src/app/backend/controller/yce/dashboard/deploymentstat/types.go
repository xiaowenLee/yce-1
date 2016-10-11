package deploymentstat

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

type DeploymentsType struct {
	DcId 	       int32 `json:"dcId"`
	DcName 	       string `json:"dcName"`
	DeploymentName string   `json:"deploymentName"`
	RsName         string   `json:"rsName"`
	PodName        []string `json:"podName"`
}

type DeploymentStatType struct {
	DeploymentStat []DeploymentsType `json:"deploymentStat"`
}
