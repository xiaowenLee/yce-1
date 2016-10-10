package resourcestat

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

type ResourceType struct {
	Cpu int32 `json:"cpu"`
	Mem int32 `json:"mem"`
}

type OverallType struct {
	Total ResourceType `json:"total"`
	Used ResourceType `json:"used"`
}

type DatacentersType struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
	Total ResourceType `json:"total"`
	Used ResourceType `json:"used"`
}

type ResourceStat struct {
	Overall OverallType `json:"overall"`
	Datacenters []DatacentersType `json:"datacenters"`
}