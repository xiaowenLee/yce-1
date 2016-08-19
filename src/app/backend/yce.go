package main

import (
	"log"
	"github.com/kataras/iris"
	"app/backend/common/util/mysql"
	myrouter "app/backend/common/yce/router"
	mysession "app/backend/common/util/session"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()
	mysession.NewSessionStore()
}

func main() {

	r := myrouter.NewRouter()
	r.Registe()

	iris.Listen(":8080")
}
