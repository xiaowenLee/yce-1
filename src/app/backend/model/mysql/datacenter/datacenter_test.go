package datacenter

import (
	"fmt"
	"testing"
	testclient "app/backend/common/yce/testclient"
)

func Test_NewDataCenter(*testing.T){
	dc := NewDataCenter("testDatacenter", "172.21.1.11", "", "", "add testDatacenter", 8080, 7)
	fmt.Printf("%v\n", dc)
}

func Test_QueryDataCenterById(t *testing.T) {
	testclient.Instance().ConnectDB()

	dc := new(DataCenter)
	err := dc.QueryDataCenterById(1)
	if err != nil {
		t.Fatalf("Error: %s", err)
		t.Fail()
	}

	fmt.Printf("DataCenter: %d, %s, %s", dc.Id, dc.Name, dc.Host)

}

func Test_QueryDataCenterByName(t *testing.T) {
	testclient.Instance().ConnectDB()

	dc := new(DataCenter)
	err := dc.QueryDataCenterByName("办公网")
	if err != nil {
		t.Fatalf("Error: %s", err)
		t.Fail()
	}

	fmt.Printf("DataCenter: %d, %s, %s", dc.Id, dc.Name, dc.Host)
}

func Test_QueryDuplicatedName(t *testing.T) {
	testclient.Instance().ConnectDB()

	dc := new(DataCenter)
	err := dc.QueryDuplicatedName("办公网")
	if err != nil {
		fmt.Printf("DataCenter doesn't have duplicated name: %d, %s, %s", dc.Id, dc.Name, dc.Host)
		t.Fatalf("Error: %s", err)
		t.Fail()
	}

	fmt.Printf("DataCenter has duplicated name: %d, %s, %s", dc.Id, dc.Name, dc.Host)
}

func Test_InsertDataCenter(t *testing.T) {
	testclient.Instance().ConnectDB()

	dc := &DataCenter{
		Name: "testDatacenter",
		Host: "172.21.1.11",
		Port: 8080,
	}

	err := dc.InsertDataCenter(7)
	if err != nil {
		t.Fatalf("Error: %s", err)
		t.Fail()
	}

	fmt.Printf("DataCenter: %d, %s, %s", dc.Id, dc.Name, dc.Host)
}

func Test_UpdateDataCenter(t *testing.T) {
	testclient.Instance().ConnectDB()

	dc := new(DataCenter)
	dc.QueryDuplicatedName("testDatacenter")

	dc.Name = "testDatacenter-abc"

	err := dc.UpdateDataCenter(7)
	if err != nil {
		t.Fatalf("Error: %s", err)
		t.Fail()
	}

	fmt.Printf("DataCenter: %d, %s, %s", dc.Id, dc.Name, dc.Host)
}

func Test_DeleteDataCenter(t *testing.T) {
	testclient.Instance().ConnectDB()

	dc := new(DataCenter)
	dc.QueryDataCenterByName("testDatacenter-abc")

	err := dc.DeleteDataCenter(7)
	if err != nil {
		t.Fatalf("Error: %s", err)
		t.Fail()
	}

	fmt.Printf("DataCenter: %d, %s, %s", dc.Id, dc.Name, dc.Host)
}

func Test_QueryAllDatacenters(t *testing.T) {
	testclient.Instance().ConnectDB()

	dcList, err := QueryAllDatacenters()
	if err != nil {
		t.Fatalf("Error: %s\n", err)
		t.Fail()
	}

	for _, dc := range dcList {
		fmt.Printf("DataCenter: %d, %s, %s\n", dc.Id, dc.Name, dc.Host)
	}
}

/*
func Test_NewDataCenter(*testing.T) {
	dc := NewDataCenter("dianxin", "10.149.149.3", "", "add dianxin", 8080, 2)
	fmt.Printf("%v\n", dc)
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
	fmt.Printf("%s\n", dc.EncodJson())

	d := new(DataCenter)
	d.DecodeJson(d.EncodeJson())
	fmt.Printf("%v\n", d)
}

/*
func Test_InsertDataCenter(t *testing.T) {
	fmt.Println("Test_InsertDataCenter")
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	dc := NewDataCenter("dianxin", "10.149.149.3", "", "add dianxin", 8080, 2)
	dc.InsertDataCenter(2)
}
*/
