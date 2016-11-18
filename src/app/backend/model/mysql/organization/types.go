package organization

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

type QuotaType struct {
	CpuQuota int32 `json:"cpuQuota"`
	MemQuota int32 `json:"memQuota"`
}

type Organization struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	CpuQuota   int32  `json:"cpu_quota"`
	MemQuota   int32  `json:"mem_quota"`
	Budget     string `json:"buget"`
	Balance    string `json:"balance"`
	Status     int32  `json:"status"`
	DcIdList   string `json:"dcIdList"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment,omitempty"`
}
