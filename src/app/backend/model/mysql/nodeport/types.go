package nodeport


import (
	mylog "app/backend/common/util/log"
)


var log = mylog.Log

type NodePort struct {
	Port       int32  `json:"port"`
	DcId       int32  `json:"dcId"`
	SvcName    string `json:"svcName"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `jsno:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}
