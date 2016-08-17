package quota

import (
	mysql "app/backend/common/util/mysql"
	"testing"
	// encrypt "app/backend/common/util/encrypt"
	"fmt"
)

func Test_QueryQuotaById(*testing.T) {
	fmt.Println("Test_QueryQuotaById")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	q := new(Quota)
	q.QueryQuotaById(2)
	fmt.Printf("%v\n", q)
}

/*
func Test_InsertQuota(*testing.T) {
	fmt.Println("Test_InsertQuota")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	q := NewQuota("quota", "100000", "add quota", 200, 400, 500, 2)
	q.InsertQuota(2)
	fmt.Printf("%v\n", q)
}

func Test_DeleteQuota(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	q := new(Quota)
	q.QueryQuotaById(1)
	q.DeleteQuota(3)
}

func Test_EncodeJSON_DecodeJson(*testing.T) {
	q := NewQuota("quota", "100000", "add quota", 200, 400, 500, 2)
	fmt.Println(q.EncodeJson())

	quota := new(Quota)
	quota.DecodeJson(q.EncodeJson())

	fmt.Printf("%v\n", quota)
}

func Test_NewQuota(*testing.T) {
	q := NewQuota("quota", "100000", "add quota", 200, 400, 500, 2)
	fmt.Printf("%v\n", q)
}

func Test_UpdateQuota(*testing.T) {
	fmt.Println("Test_UpdateQuota")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	q := new(Quota)
	q.QueryQuotaById(1)

	q.Name = "LimitRange"
	q.UpdateQuota(2)

	quota := new(Quota)
	quota.QueryQuotaById(2)

	fmt.Printf("%v\n", quota)

}
*/
