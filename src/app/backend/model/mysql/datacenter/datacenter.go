package datacenter

import (
	"time"
)

type DataCenter struct {
	Id             int32     `json:"id"`
	Name           string    `json:"name"`
	Host           string    `json:"host"`
	Port           int32     `json:"port"`
	Secret         string    `json:"secret"` // maybe error
	CreatedTs      time.Time `json:"created_ts"`
	LastModifiedTs time.Time `json:"last_modified_ts"`
	LastModifiedOp int       `json:"last_modified_op"`
	Comment        string    `json:"comment"`
}
