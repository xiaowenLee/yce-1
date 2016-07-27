package deployment

import (
	"time"
)

type Deployment struct {
	Id         int32     `json:"id"`
	Name       string    `json:"name"`
	ActionType int32     `json:"actionType"`
	ActionVerb string    `json:"actionVerb"`
	ActionUrl  string    `json:"actionUrl"`
	ActionTs   string `json:"actionTs"`
	ActionOp   int32     `json:"actionOp"`
	DcList     string    `json:"dcList"`
	Success    int32     `json:"success"`
	Reason     string    `json:"reason"`
	Json       string    `json:"json"`
	Comment    string    `json:"comment"`
}
