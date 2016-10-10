package user

import (
	myuser "app/backend/model/user"
	"app/backend/controller/yce"
)

type CreateUserController struct {
	yce.Controller
}

type CreateUserParams struct {

}

Id         int32  `json:"id"`
Name       string `json:"name"`
OrgId      int32  `json:"orgId"`
Password   string `json:"password"`
Status     int32  `json:"status"`
CreatedAt  string `json:"createdAt"`
ModifiedAt string `json:"modifiedAt"`
ModifiedOp int32  `json:"modifiedOp"`
Comment    string `json:"comment"`

func (cuc CreateUserController) Post() {

}
