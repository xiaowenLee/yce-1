package resourcestat

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

type UsageType struct {
	Total int32 `json:"total"`
	Used  int32 `json:"used"`
}

type DatacentersType struct {
	DcId   int32     `json:"dcId"`
	DcName string    `json:"dcName"`
	Cpu    UsageType `json:"cpu"`
	Mem    UsageType `json:"mem"`
}

// treat overall as one special datacenter with propertities: dcId=-1 and dcName=总览
type ResourceStat struct {
	Datacenters []DatacentersType `json:"datacenters"`
}
