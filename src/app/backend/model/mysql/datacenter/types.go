package datacenter


import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

type DataCenter struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int32  `json:"port"`
	Secret     string `json:"secret"` // maybe error
	Status     int32  `json:"status"`
	NodePort   string `json:"nodePort"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}
