package main

import (
	"app/backend/common/util/mysql"
	mylogin "app/backend/controller/yce/login"
	"github.com/kataras/iris"
	mysession "app/backend/common/util/session"
    // mydeploy "app/backend/controller/yce/deploy"
	mylogout "app/backend/controller/yce/logout"
	mynavList "app/backend/controller/yce/navlist"
	myregistry "app/backend/controller/yce/registry"
)

func main() {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	mysession.NewSessionStore()

	login := new(mylogin.LoginController)
	logout := new(mylogout.LogoutController)
	nav := new(mynavList.NavListController)
	// listdeploy := new(mydeploy.ListDeployController)
	registry := new(myregistry.ListRegistryController)

	iris.API("/api/v1/users/login", *login)
	iris.API("/api/v1/navlist", *nav)
	// iris.API("/api/v1/organization/:id/deployments", *listdeploy)
	iris.API("/api/v1/users/logout", *logout)
	iris.API("/api/v1/registry/images", *registry)


	iris.StaticServe("../frontend", "/static")
	iris.Listen(":8080")

}
