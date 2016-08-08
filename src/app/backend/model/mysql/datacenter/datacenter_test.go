package datacenter

import (
	mysql "app/backend/common/util/mysql"
	"fmt"
	"testing"
)

func Test_NewDataCenter(*testing.T) {
	dc := NewDataCenter("dianxin", "10.149.149.3", "", "add dianxin", 8080, 2)
	fmt.Printf("%v\n", dc)
}

func Test_InsertDataCenter(t *testing.T) {
	fmt.Println("Test_InsertDataCenter")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	dc := NewDataCenter("dianxin", "10.149.149.3", "", "add dianxin", 8080, 2)
	dc.InsertDataCenter(2)
}

func Test_QueryDataCenterById(*testing.T) {
	fmt.Println("Test_QueryDataCenter")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	dc := new(DataCenter)
	dc.QueryDataCenterById(1)
	fmt.Printf("%v\n", dc)

}

func Test_UpdateDataCenter(t *testing.T) {
	fmt.Println("Test_UpdateDataCenter")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	dc := new(DataCenter)
	dc.QueryDataCenterById(2)

	dc.Host = "172.21.1.11"
	dc.Name = "bangongwang"
	dc.UpdateDataCenter(2)

	dc.QueryDataCenterById(2)
	fmt.Printf("%v\n", dc)
}

func Test_DeleteDataCenter(t *testing.T) {
	fmt.Println("Test_DeleteDataCenter")

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	dc := new(DataCenter)
	dc.QueryDataCenterById(2)
	dc.DeleteDataCenter(2)
}

func Test_EncodeJson_DecodeJson(*testing.T) {

	dc := NewDataCenter("dianxin", "10.149.149.3", "", "add dianxin", 8080, 2)
	fmt.Printf("%s\n", dc.EncodeJson())

	d := new(DataCenter)
	d.DecodeJson(d.EncodeJson())
	fmt.Printf("%v\n", d)
}
