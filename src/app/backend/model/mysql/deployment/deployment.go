package deployment

import (
	"time"
)

type Deployment struct {
	Id         int32     `json:"id"`
	Name       string    `json:"name"`
	ActionType int32     `json:"action_type"`
	ActionVerb string    `json:"action_verb"`
	ActionUrl  string    `json:"action_url"`
	ActionTs   time.Time `json:"action_ts"`
	ActionOp   int32     `json:"action_op"`
	DcList     string    `json:"dc_list"`
	Success    int32     `json:"success"`
	Reason     string    `json:"reason"`
	Json       string    `json:"json"`
	Comment    string    `json:"comment"`
}
