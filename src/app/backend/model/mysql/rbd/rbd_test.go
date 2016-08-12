package rbd

import (
	mysql "app/backend/common/util/mysql"
	"fmt"
	"testing"
)

func Test_NewRbd(*testing.T) {
	r := NewRbd("arbd", "rbd", "add rbd", 100, 2, 2, 2)
	fmt.Printf("%v\n", r)
}

func Test_InsertRbd(*testing.T) {
	fmt.Println("Test_InsertDcQuota")

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	r := NewRbd("arbd", "rbd", "add rbd", 100, 2, 2, 2)
	r.InsertRbd(2)
}

func Test_QueryRbdById(*testing.T) {
	fmt.Println("Test_QueryDcQuotaById")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	r := new(Rbd)
	r.QueryRbdById(1)
	fmt.Printf("%v\n", r)
}

func Test_UpdateRbd(*testing.T) {
	fmt.Println("Test_UpdateRbd")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	r := new(Rbd)
	r.QueryRbdById(2)

	r.FileSystem = "xfs"
	r.Name = "newRbd"

	r.UpdateRbd(2)

	r.QueryRbdById(2)

	fmt.Printf("%v\n", r)
}

func Test_DeleteRbd(*testing.T) {
	fmt.Println("Test_DeleteRbd")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	r := new(Rbd)
	r.QueryRbdById(3)
	r.DeleteRbd(2)

}

func Test_EncodeJson_DecodeJson(*testing.T) {
	r := NewRbd("brbd", "rbd", "add rbd", 100, 2, 2, 2)

	fmt.Printf("%s\n", r.EncodeJson())

	rbd := new(Rbd)
	rbd.DecodeJson(r.EncodeJson())
	fmt.Printf("%v\n", rbd)
}
