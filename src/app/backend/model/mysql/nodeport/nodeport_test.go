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
/*
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
		fmt.Printf("Insert Failed: %s\n", err)
	} else {
		fmt.Printf("Insert Succeed")
	}
}
*/
func TestNodePort_UpdateNodePortByPortAndDcId(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	np := &NodePort{
		Port: 30061,
		DcId: 1,
		SvcId: 10,
		Comment: "lb service",
	}

	err := np.UpdateNodePortByPortAndDcId(1)
	if err != nil {
		fmt.Printf("Update Failed: %s\n", err)
	} else {
		fmt.Printf("Update Succeed")
	}
}
/*
func TestNodePort_UpdateNodePortByPortAndDcId(t *testing.T) {

}
}

func TestNodePort_QueryNodePortByPortAndDcId(t *testing.T) {

}

func TestNodePort_DeleteNodePortByPortAndDcId(t *testing.T) {

}
*/
/*
func TestNodePort_QueryNodePortByPort(t *testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	d := new(NodePort)

	if err := d.QueryNodePortByPortAndDcId(29999, 1); err != nil {
		fmt.Println("No such Port")
	} else {
		fmt.Println(err)
	}

	if err := d.QueryNodePortByPortAndDcId(30000, 1); err != nil {
		fmt.Println("No such Port")
	} else {
		fmt.Println(err)
	}

	if err := d.QueryNodePortByPortAndDcId(32767, 1); err != nil {
		fmt.Println("No such Port")
	} else {
		fmt.Println(err)
	}

	if err := d.QueryNodePortByPortAndDcId(32768, 1); err != nil {
		fmt.Println("No such Port")
	} else {
		fmt.Println(err)
	}

	if err := d.QueryNodePortByPortAndDcId(30061, 1); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Find Port")
	}

}
func TestNodePort_All(t *testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	d := new(NodePort)

	if err := d.InsertNodePort(30000, 1, "none app"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Insert OK")
	}

	if err := d.InsertNodePort(30061, 1, "1st app"); err != nil {
		fmt.Println("Insert Error")
	} else {
		fmt.Println(err)
	}
	if err := d.InsertNodePort(30061, 2, "2nd app"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Insert OK")
	}


	if err := d.UpdateNodePortByPortAndDcId(30000, 1, "update app"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Update OK")
	}

	if err := d.UpdateNodePortByPortAndDcId(30000, 2, "none app"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Update Error")
	}

	if err := d.UpdateNodePortByPortAndDcId(30001, 1, "none app"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Update Error")
	}

	if err := d.DeleteNodePortByPortAndDcId(30000, 1); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Delete Success")
	}

	if err := d.DeleteNodePortByPortAndDcId(30000, 2); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Delete Failed")
	}

	if err := d.DeleteNodePortByPortAndDcId(30001, 1); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Delete Failed")
	}

}
*/
/*
func TestNodePort_InsertNodePort(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()



	d := new(NodePort)
	err := d.InsertNodePort(30061, 2, "lb service")
	fmt.Println(err)
}

func TestNodePort_UpdateNodePortByPort(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	d := new(NodePort)
	err := d.UpdateNodePortByPort(30061, 1, "lb service occupied")
	fmt.Println(err)
}

func TestNodePort_DeleteNodePortByPort(t *testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	d := new(NodePort)
	err := d.DeleteNodePortByPort(2, 30061)
	fmt.Println(err)
}
*/
