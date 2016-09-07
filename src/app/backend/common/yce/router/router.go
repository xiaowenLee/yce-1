package router

import (
	mydeploy "app/backend/controller/yce/deploy"
	mylogin "app/backend/controller/yce/login"
	mylogout "app/backend/controller/yce/logout"
	mynavList "app/backend/controller/yce/navlist"
	myregistry "app/backend/controller/yce/registry"
	myservice "app/backend/controller/yce/service"
	myendpoint "app/backend/controller/yce/endpoint"
	myextensions "app/backend/controller/yce/extensions"
	mynamespace "app/backend/controller/yce/namespace"
	"github.com/kataras/iris"
)

type Router struct {
	Login        *mylogin.LoginController
	Logout       *mylogout.LogoutController
	Nav          *mynavList.NavListController
	ListDeploy   *mydeploy.ListDeployController
	Registry     *myregistry.ListRegistryController
	ListService      *myservice.ListServiceController
	InitService *myservice.InitServiceController
	CreateService *myservice.CreateServiceController
	ListEndpoints     *myendpoint.ListEndpointsController
	InitEndpoints *myendpoint.InitEndpointsController
	CreateEndpoints *myendpoint.CreateEndpointsController
	ListExtensions *myextensions.ListExtensionsController
	InitDeploy   *mydeploy.InitDeployController
	CreateDeploy *mydeploy.CreateDeployController
	RollingDeploy *mydeploy.RollingDeployController
	RollbackDeploy *mydeploy.RollbackDeployController
	ScaleDeploy *mydeploy.ScaleDeploymentController
	ListOperationLog *mydeploy.ListOperationLogController
	InitNamespace *mynamespace.InitNamespaceController
	DeleteService *myservice.DeleteServiceController
	DeleteEndpoint *myendpoint.DeleteEndpointsController
	HistoryDeploy *mydeploy.HistoryDeployController
}

func NewRouter() *Router {
	r := new(Router)
	r.Login = new(mylogin.LoginController)
	r.Logout = new(mylogout.LogoutController)
	r.Nav = new(mynavList.NavListController)
	r.ListDeploy = new(mydeploy.ListDeployController)
	r.Registry = new(myregistry.ListRegistryController)
	r.ListService = new(myservice.ListServiceController)
	r.InitService = new(myservice.InitServiceController)
	r.CreateService = new(myservice.CreateServiceController)
	r.ListEndpoints = new(myendpoint.ListEndpointsController)
	r.InitEndpoints = new(myendpoint.InitEndpointsController)
	r.CreateEndpoints = new(myendpoint.CreateEndpointsController)
	r.ListExtensions = new(myextensions.ListExtensionsController)
	r.InitDeploy = new(mydeploy.InitDeployController)
	r.CreateDeploy = new(mydeploy.CreateDeployController)
	r.RollingDeploy = new(mydeploy.RollingDeployController)
	r.RollbackDeploy = new(mydeploy.RollbackDeployController)
	r.ScaleDeploy = new(mydeploy.ScaleDeploymentController)
	r.ListOperationLog = new(mydeploy.ListOperationLogController)
	r.InitNamespace = new(mynamespace.InitNamespaceController)
	r.DeleteService = new(myservice.DeleteServiceController)
	r.DeleteEndpoint = new(myendpoint.DeleteEndpointsController)
	r.HistoryDeploy = new(mydeploy.HistoryDeployController)

	return r
}

func (r *Router) Registe() {

	iris.API("/api/v1/users/login", *r.Login)
	iris.API("/api/v1/users/logout", *r.Logout)
	iris.API("/api/v1/navlist", *r.Nav)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments", *r.ListDeploy)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments/init", *r.InitDeploy)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments/new", *r.CreateDeploy)
	iris.API("/api/v1/organizations/:orgId/deployments/:deploymentName/rolling", *r.RollingDeploy)
	iris.API("/api/v1/organizations/:orgId/deployments/:deploymentName/rollback", *r.RollbackDeploy)
	iris.API("/api/v1/organizations/:orgId/deployments/:deploymentName/scale", *r.ScaleDeploy)
	iris.API("/api/v1/organizations/:orgId/operationlog", *r.ListOperationLog)
	iris.API("/api/v1/registry/images", *r.Registry)
	iris.API("/api/v1/organizations/:orgId/users/:userId/services", *r.ListService)
	iris.API("/api/v1/organizations/:orgId/users/:userId/services/init", *r.InitService)
	iris.API("/api/v1/organizations/:orgId/users/:userId/services/new", *r.CreateService)
	iris.API("/api/v1/organizations/:orgId/datacenters/:dcId/users/:userId/services/:svcName", *r.DeleteService)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints", *r.ListEndpoints)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints/init", *r.InitEndpoints)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints/new", *r.CreateEndpoints)
	iris.API("/api/v1/organizations/:orgId/datacenters/:dcId/endpoints/:epName", *r.DeleteEndpoint)
	iris.API("/api/v1/organizations/:orgId/users/:userId/extensions", *r.ListExtensions)
	iris.API("/api/v1/organizations/init", *r.InitNamespace)
	iris.API("/api/v1/organizations/:orgId/datacenters/:dcId/deployments/:name/history", *r.HistoryDeploy)

	iris.StaticServe("../frontend", "/static")
}
