package rbd

import (
	"time"
)

type Rbd struct {
	Id             int32     `json:"id"`
	Name           string    `json:"name"`
	Pool           string    `json:"pool"`
	Size           int32     `json:"size"`
	FileSystem     string    `jsonN:"filesystem"`
	OrgId          int32     `json:"org_id"`
	DcID           int32     `json:"dc_id"`
	CreatedTs      time.Time `json:"created_ts"`
	LastModifiedTs time.Time `json:"last_modified_ts"`
	LastModifiedOp int       `json:"last_modified_op"`
	Comment        string    `json:"comment"`
}
