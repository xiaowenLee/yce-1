package rbd

import (
	"time"
)

type Rbd struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Pool       string `json:"pool"`
	Size       int32  `json:"size"`
	FileSystem string `jsonN:"filesystem"`
	OrgId      int32  `json:"orgId"`
	DcID       int32  `json:"dcId"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedTs string `json:"modifiedAt"`
	ModifiedOp int    `json:"modifiedOp"`
	Comment    string `json:"comment"`
}
