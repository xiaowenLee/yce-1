package nodeport

import (
	mysql "app/backend/common/util/mysql"
	"fmt"
	"testing"
)

/*
func TestNewNodePort(t *testing.T) {
	np := NewNodePort(30061, 1, "lb-service.ops", 1)
	fmt.Printf("%p\n", np)
}
*/

/*
func TestNodePort_InsertNodePort(t *testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := &NodePort{
		Port: 32380,
		DcId: 3,
		Status: INVALID,
		SvcName: "re-redis-slave.default",
	}

	err := np.InsertNodePort(1)
	if err != nil {
		fmt.Printf("Insert Failed or exist: %s\n", err)
	} else {
		fmt.Printf("Insert Succeed")
	}
}
*/
/*
func TestNodePort_UpdateNodePort(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := &NodePort{
		Port: 32380,
		DcId: 4,
		SvcName: "redis-slave.default",
		Comment: "redis slave service default",
	}

	err := np.UpdateNodePort(1)
	if err != nil {
		fmt.Printf("Update Failed: %s\n", err)
	} else {
		fmt.Printf("Update Succeed or not exist:")
	}
}
*/
/*
func TestNodePort_QueryNodePortByPortAndDcId(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := new(NodePort)
	err := np.QueryNodePortByPortAndDcId(32380, 4)
	if err != nil {
		fmt.Printf("Query Failed or not exist! ")
	} else {
		fmt.Printf("Found")
	}
}
*/
/*
func TestNodePort_DeleteNodePort(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := &NodePort{
		Port: 32380,
		DcId: 4,
	}

	err := np.DeleteNodePort(1)
	if err != nil {
		fmt.Printf("Delete Failed: %s\n", err)
	} else {
		fmt.Printf("Delete Succeed or not exist")
	}
}
*/
/*
func TestNodePort_QueryServiceIdByPortAndDcId(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()


	np := new(NodePort)
	svcName, err := np.QueryServiceNameByPortAndDcId(32380, 4)
	if err != nil {
		fmt.Printf("Query ServiceId Failed")
	} else {
		fmt.Printf("Query Service Id Successed: svcName=%s", svcName)
	}
}
*/
/*
func TestNodePort_InsertOnDuplicateKeyUpdateStatusn(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := &NodePort{
		Port: 32380,
		DcId: 3,
		SvcName: "Delete-Yeah-TEST-redis-slave.default",
	}

	err := np.InsertOnDuplicateKeyUpdateStatus(1)
	if err != nil {
		fmt.Printf("Insert On Duplicate Key Update Status Failed")
	} else {
		fmt.Printf("Insert On Duplicate Key Update Status Successed")
	}

}
*/
