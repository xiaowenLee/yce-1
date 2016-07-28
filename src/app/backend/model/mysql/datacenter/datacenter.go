package datacenter

import (
	"time"
)

type DataCenter struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int32  `json:"port"`
	Secret     string `json:"secret"` // maybe error
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int    `json:"modifiedOp"`
	Comment    string `json:"comment"`
}
