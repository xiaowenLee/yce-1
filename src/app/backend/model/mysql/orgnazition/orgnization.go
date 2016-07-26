package mysql

import (
	"time"
)

type Orgnazition struct {
	Id             int32     `json:"id"`
	Name           string    `json:"name"`
	CpuQuota       int32     `json:"cpu_quota"`
	MemQuota       int32     `json:"mem_quota"`
	Budget         int32     `json:"buget"`
	Balance        int32     `json:"balance"`
	CreatedTs      time.Time `json:"created_ts"`
	LastModifiedTs time.Time `json:"last_modified_ts"`
	LastModifiedOp int       `json:"last_modified_op"`
	Comment        string    `json:"comment"`
}
