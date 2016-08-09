package user

import (
	"testing"
	mysql "app/backend/common/util/mysql"
	encrypt "app/backend/common/util/encrypt"
	// "fmt"
)

/*
func Test_NewUser(*testing.T) {
	user := NewUser("dawei.li", "123456", "add dawei.li", 1, VALID, 2)
	fmt.Printf("User Name: %s\n", user.Name)
	fmt.Printf("User Comment: %s\n", user.Comment)
	fmt.Printf("User CreateAt: %s\n", user.CreatedAt)
	fmt.Printf("User ModifiedAt: %s\n", user.ModifiedAt)
}

func Test_EncodeJson_DecodeJson(*testing.T) {

	user := NewUser("dawei.li", "123456", "add dawei.li", 1, VALID, 2)
	fmt.Println(user.EncodeJson())

	u := new(User)
	str, _ := user.EncodeJson()
	u.DecodeJson(str)

	fmt.Println(u.Name)
}

func Test_GetUserByNameAndPassword(t *testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	// exists
	user := new(User)
	user.QueryUserByNameAndPassword("jingru.zhang", "234567")
	fmt.Printf("%v\n", user)

	// not exists
	u := new(User)
	err := user.QueryUserByNameAndPassword("jingru.zhang", "123456")

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%v\n", u)

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

	user := NewUser("dawei.li.richard", "123456", "add dawei.li", 1, VALID, 2)
	user.InsertUser(2)
}

*/

func Test_DeleteUser(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	user := new(User)
	user.QueryUserById(6)
	user.DeleteUser(3)
}

func Test_UpdateUser(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	user := new(User)
	user.QueryUserById(6)

	user.Password = "234567"
	user.UpdateUser(2)

	u := new(User)
	u.QueryUserById(7)
	u.Password = encrypt.NewEncryption("hello").String()
	u.OrgId = 1
	u.UpdateUser(2)
}
