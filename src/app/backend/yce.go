package main


import (
	"github.com/kataras/iris"
	mylogin "app/backend/controller/yce/login"
	"app/backend/common/util/mysql"
	mysession "app/backend/common/util/session"
)

func main() {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	mysession.NewSessionStore()

	lc := new(mylogin.LoginController)

	iris.API("/api/v1/users/:email/login", *lc)

	iris.Listen(":8080")

}
