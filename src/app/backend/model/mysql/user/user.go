package user

import (
	mysql "app/backend/common/util/mysql"
	"log"
	localtime "app/backend/common/util/time"
	"fmt"
)

const (
	USER_SELECT = "SELECT id, name, orgId, createdAt, modifiedAt, modifiedOp FROM user WHERE id = ? "
	USER_INSERT = "INSERT INTO user(name, password, orgId, status, createdAt, modifiedAt, modifiedOp, comment) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	USER_UPDATE = ""
	USER_DELETE = "UPDATE user SET status=0 WHERE id = ?"
	VALID = 1
	INVALID = 0
)

type User struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	OrgId      int32 `json:"orgId"`
	Password   string `json:"password"`
	Status int32 `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32    `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

func NewUser(name, password, comment string, id, orgId, status, modifiedOp int32) *User {

	return &User{
		Id: id,
		Name: name,
		Password: password,
		OrgId: orgId,
		Status: status,
		Comment: comment,
		ModifiedAt: localtime.NewLocalTime().String(),
		CreatedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
	}
}

func (u *User) QueryUserById(id int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(USER_SELECT)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer  stmt.Close()

	// Query user by id
	stmt.QueryRow(id).Scan(&u.Id, &u.Name, &u.OrgId, &u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	fmt.Printf("%v\n", u)
}

func (u *User) InsertUser() {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(USER_INSERT)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	defer stmt.Close()

	// Insert a user
	_, err = stmt.Exec(u.Name, u.Password, u.OrgId, u.Status, u.CreatedAt, u.ModifiedAt, u.ModifiedOp, u.Comment)

	if err != nil {
		log.Fatal(err)
		panic(err.Error())

	}
}

func (u *User) DeleteUser() {
	db := mysql.MysqlInstance().Conn()

	// Prepare delete-statement
	stmt, err := db.Prepare(USER_DELETE)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	defer stmt.Close()

	// Set user status  INVALED
	_, err = stmt.Exec(u.Id)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}
