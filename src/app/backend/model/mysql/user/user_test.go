package mysql

import (
	"testing"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


// INSERT INTO USERS(name, password, org_id, created_ts, last_modified_ts, last_modifed_op) VALUES('litanhua', 'root', 0, now(), now(), 0)

func Test_Qurey_UserName(*testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(172.21.1.11:32306)/yce?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var str string
	q := "SELECT name from yce.users"
	err = db.QueryRow(q).Scan(&str)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(str)
}

func Test_Query_User(*testing.T) {

	db, err := sql.Open("mysql", "root:root@tcp(172.21.1.11:32306)/yce?parseTime=true")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	user := new(User)
	q := "SELECT id, name, password, org_id, created_ts, last_modified_ts, last_modifed_op from yce.users"

	err = db.QueryRow(q).Scan(
		&user.Id,
		&user.Name,
		&user.Password,
		&user.OrgId,
		&user.CreatedTs,
		&user.LastModifiedTs,
		&user.LastModifiedOp)
	// err = db.QueryRow(q).Scan(&user)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v\n", user)
	// fmt.Println(user.LastModifiedTs)
}
