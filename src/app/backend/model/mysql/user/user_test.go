package user

import (
	"fmt"
	"testing"
	mysql "app/backend/common/util/mysql"
)

// INSERT INTO USERS(name, password, org_id, created_ts, last_modified_ts, last_modifed_op) VALUES('litanhua', 'root', 0, now(), now(), 0)

func Test_NewUser(*testing.T) {
 	user := NewUser("dawei.li", "123456", "add dawei.li", 3, 1, VALID, 2)
	fmt.Printf("User Name: %s\n", user.Name)
	fmt.Printf("User Comment: %s\n", user.Comment)
	fmt.Printf("User CreateAt: %s\n", user.CreatedAt)
	fmt.Printf("User ModifiedAt: %s\n", user.ModifiedAt)
}

func Test_GetUserByID(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	user := new(User)
	user.QueryUserById(2)
	fmt.Printf("%s\n", user.Name)
}

func Test_InsertUser(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	user := NewUser("dawei.li", "123456", "add dawei.li", 3, 1, VALID, 2)
	user.InsertUser()
}

func Test_DeleteUser(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	user := new(User)
	user.QueryUserById(6)
	user.DeleteUser()
}

/*
func Test_Qurey_UserName(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	db := mysql.MysqlInstance().Conn()

	var str string
	q := "SELECT name from yce.user where id = 2"
	err := db.QueryRow(q).Scan(&str)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(str)
}

/*
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
*/
