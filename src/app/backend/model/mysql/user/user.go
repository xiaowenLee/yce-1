package user

import (
	"time"
)

const (
	USER_SELECT = "SELECT "
	USER_INSERT = ""
	USER_UPDATE = ""
	USER_DELETE = ""
)

type User struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	OrgId      string `json:"orgId"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int    `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

func (u *User) AddUser() {

}
