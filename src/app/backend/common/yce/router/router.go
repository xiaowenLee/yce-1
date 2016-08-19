package router

import (
	mydeploy "app/backend/controller/yce/deploy"
	mylogin "app/backend/controller/yce/login"
	mylogout "app/backend/controller/yce/logout"
	mynavList "app/backend/controller/yce/navlist"
	myregistry "app/backend/controller/yce/registry"
	"github.com/kataras/iris"
)

type Router struct {
}

func (r *Router) Registe() {
	login := new(mylogin.LoginController)
	logout := new(mylogout.LogoutController)
	nav := new(mynavList.NavListController)
	listdeploy := new(mydeploy.ListDeployController)
	registry := new(myregistry.ListRegistryController)

	iris.API("/api/v1/users/login", *login)
	iris.API("/api/v1/navlist", *nav)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments", *listdeploy)
	iris.API("/api/v1/users/logout", *logout)
	iris.API("/api/v1/registry/images", *registry)

	iris.StaticServe("../frontend", "/static")
}
