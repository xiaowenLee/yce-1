package nodeport

import (
	mysql "app/backend/common/util/mysql"
	"fmt"
	"testing"
)

func TestNewNodePort(t *testing.T) {
	np := NewNodePort(30061, 1, 1, 1)
	fmt.Printf("%p\n", np)
}

func TestNodePort_InsertNodePort(t *testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := &NodePort{
		Port: 32380,
		DcId: 1,
		SvcId: 4,
	}

	err := np.InsertNodePort(1)
	if err != nil {
		fmt.Printf("Insert Failed or exist: %s\n", err)
	} else {
		fmt.Printf("Insert Succeed")
	}
}

func TestNodePort_UpdateNodePortByPortAndDcId(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := &NodePort{
		Port: 32380,
		DcId: 2,
		SvcId: 4,
		Comment: "redis slave service",
	}

	err := np.UpdateNodePortByPortAndDcId(1)
	if err != nil {
		fmt.Printf("Update Failed: %s\n", err)
	} else {
		fmt.Printf("Update Succeed or not exist:")
	}
}


func TestNodePort_QueryNodePortByPortAndDcId(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := new(NodePort)
	err := np.QueryNodePortByPortAndDcId(32380, 2)
	if err != nil {
		fmt.Printf("Query Failed or not exist! ")
	} else {
		fmt.Printf("Found")
	}
}


func TestNodePort_DeleteNodePortByPortAndDcId(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := &NodePort{
		Port: 32380,
		DcId: 2,
		SvcId: 4,
	}

	err := np.DeleteNodePortByPortAndDcId(1)
	if err != nil {
		fmt.Printf("Delete Failed: %s\n", err)
	} else {
		fmt.Printf("Delete Succeed or not exist")
	}
}

func TestNodePort_QueryServiceIdByPortAndDcId(t *testing.T) {
	np := new(NodePort)
	svcId, err := np.QueryServiceIdByPortAndDcId(32306, 2)
	if err != nil {
		fmt.Printf("Query ServiceId Failed")
	} else {
		fmt.Printf("Query Service Id Successed: svcId=%d", svcId)
	}
}

