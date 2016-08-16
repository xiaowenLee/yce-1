package main

import (
	"app/backend/common/util/mysql"
	mylogin "app/backend/controller/yce/login"
	"github.com/kataras/iris"
	// mylogout "app/backend/controller/yce/logout"
	mysession "app/backend/common/util/session"
	// mydeploy "app/backend/controller/yce/deploy"
	mylogout "app/backend/controller/yce/logout"
	mynavList "app/backend/controller/yce/navlist"
)

func main() {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	mysession.NewSessionStore()

	login := new(mylogin.LoginController)
	logout := new(mylogout.LogoutController)
	nav := new(mynavList.NavListController)
	// deploy := new(mydeploy.CreateDeployController)

	iris.StaticServe("../frontend", "/static")

	iris.API("/api/v1/users/login", *login)
	iris.API("/api/v1/navlist", *nav)

	// iris.Get("/api/v1/organization/:oid/deployments/:id", deploy.Describe)
	// iris.Get("/api/v1/organization/:oid/deployments", deploy.List)
	// iris.Post("/api/v1/deployments", deploy.Create)

	// iris.API("/api/v1/users/:email/logout", *logout)
	iris.API("/api/v1/users/logout", *logout)

	iris.Listen(":8080")

}
