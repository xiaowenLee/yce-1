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
	Login *mylogin.LoginController
	Logout *mylogout.LogoutController
	Nav *mynavList.NavListController
	ListDeploy *mydeploy.ListDeployController
	Registry *myregistry.ListRegistryController
	InitDeploy *mydeploy.InitDeployController
	CreateDeploy *mydeploy.CreateDeployController
}

func NewRouter() *Router {
	r := new(Router)
	r.Login = new(mylogin.LoginController)
	r.Logout = new(mylogout.LogoutController)
	r.Nav = new(mynavList.NavListController)
	r.ListDeploy = new(mydeploy.ListDeployController)
	r.Registry = new(myregistry.ListRegistryController)
	r.InitDeploy = new(mydeploy.InitDeployController)
	r.CreateDeploy = new(mydeploy.CreateDeployController)

	return r
}

func (r *Router) Registe() {

	iris.API("/api/v1/users/login", *r.Login)
	iris.API("/api/v1/users/logout", *r.Logout)
	iris.API("/api/v1/navlist", *r.Nav)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments", *r.ListDeploy)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments/init", *r.InitDeploy)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments", *r.CreateDeploy)
	iris.API("/api/v1/registry/images", *r.Registry)

	iris.StaticServe("../frontend", "/static")
}
