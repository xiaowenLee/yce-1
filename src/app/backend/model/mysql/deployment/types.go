package deployment

import (
	mylog "app/backend/common/util/log"
)

var log =  mylog.Log

type Deployment struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	ActionType int32  `json:"actionType"`
	ActionVerb string `json:"actionVerb"`
	ActionUrl  string `json:"actionUrl"`
	ActionAt   string `json:"actionAt"`
	ActionOp   int32  `json:"actionOp"`
	DcList     string `json:"dcList"`
	Success    int32  `json:"success"`
	Reason     string `json:"reason",omitempty`
	Json       string `json:"json"`
	Comment    string `json:"comment,omitempty"`
	OrgId      int32  `json:"orgId"`
}

type Statistics map[int32]int32

type OperationStat map[string]Statistics // day-->op-->total
