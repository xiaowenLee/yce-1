package mysql

import (
	"time"
)

type User struct {
	Id             int32     `json:"id"`
	Name           string    `json:"name"`
	Password       string    `json:"password"`
	OrgId          string    `json:"org_id"`
	CreatedTs      time.Time `json:"created_ts"`
	LastModifiedTs time.Time `json:"last_modified_ts"`
	LastModifiedOp int       `json:"last_modified_op"`
	Comment        string    `json:"comment"`
}
