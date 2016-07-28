package mysql

import (
	"time"
)

type organization struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	CpuQuota   int32  `json:"cpu_quota"`
	MemQuota   int32  `json:"mem_quota"`
	Budget     int32  `json:"buget"`
	Balance    int32  `json:"balance"`
	CreatedTs  string `json:"createdAt"`
	ModifiedTs string `json:"modifiedAt"`
	ModifiedOp int32    `json:"modifiedOp"`
	Comment    string `json:"comment,omitempty"`
}


func NewOrgnization(name, comment string, cpuQuota, memQuota, modifiedOp int32, budget, balance decimal) {

}