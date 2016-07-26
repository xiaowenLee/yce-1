package mysql

import (
	"fmt"
	"testing"
)

func Test_MysqlClient_New(*testing.T) {
	client := NewMysqlClient(DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, MAX_POOL_SIZE)

	fmt.Printf("%v\n", client)
}

func Test_MysqlClient_Open(*testing.T) {

	client := NewMysqlClient(DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, MAX_POOL_SIZE)

	client.Open()

}

func Test_MysqlClient_Close(*testing.T) {

	client := NewMysqlClient(DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, MAX_POOL_SIZE)

	client.Open()

	client.Close()
}

func Test_MysqlClient_Ping(*testing.T) {

	client := NewMysqlClient(DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, MAX_POOL_SIZE)

	client.Open()

	// go client.Ping()

}

func TestMysqlClient_Query(*testing.T) {

	client := NewMysqlClient(DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, MAX_POOL_SIZE)

	client.Open()

	var str string
	q := "SELECT name FROM yce.users"
	err := client.Conn().QueryRow(q).Scan(&str)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(str)

	client.Close()
}
