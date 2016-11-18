package user

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log



type User struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	OrgId      int32  `json:"orgId"`
	Password   string `json:"password"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
	NavList    string `json:"navList"`
}
