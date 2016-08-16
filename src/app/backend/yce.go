package main


import (
	"github.com/kataras/iris"
	"app/backend/common/util/mysql"
	mylogin "app/backend/controller/yce/login"
	mylogout "app/backend/controller/yce/logout"
	mysession "app/backend/common/util/session"
	mynavList "app/backend/controller/yce/navlist"
)

func main() {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	mysession.NewSessionStore()

	login := new(mylogin.LoginController)
	logout := new(mylogout.LogoutController)
	nav := new(mynavList.NavListController)

	iris.StaticServe("../frontend", "/static")

	iris.API("/api/v1/users/login", *login)
	iris.API("/api/v1/navlist", *nav)
	iris.API("/api/v1/users/logout", *logout)

	iris.Listen(":8080")

}
