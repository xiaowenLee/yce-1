package rbd

import (
	"time"
)

type Rbd struct {
	Id         int32     `json:"id"`
	Name       string    `json:"name"`
	Pool       string    `json:"pool"`
	Size       int32     `json:"size"`
	FileSystem string    `jsonN:"filesystem"`
	OrgId      int32     `json:"orgId"`
	DcID       int32     `json:"dcId"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedTs time.Time `json:"modifiedAt"`
	ModifiedOp int       `json:"modifiedOp"`
	Comment    string    `json:"comment"`
}
