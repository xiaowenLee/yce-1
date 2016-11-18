package template

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

type Template struct {
	Id int32 `json:"id"`
	Name string `json:"name"`
	OrgId int32 `json:"orgId"`
	Deployment string `json:"deployment"`
	Service string `json:"service"`
	Endpoints string `json:"endpoints"`
	Status int32 `json:"status"`
	Comment string `json:"comment"`
	CreatedAt string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32 `json:"modifiedOp"`
}
