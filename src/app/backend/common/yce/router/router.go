package router

import (
	mypath "app/backend/controller/yce/apis"
	mydeploymentstat "app/backend/controller/yce/dashboard/deploymentstat"
	myresourcestat "app/backend/controller/yce/dashboard/resourcestat"
	mydeploy "app/backend/controller/yce/deploy"
	myendpoint "app/backend/controller/yce/endpoint"
	myextensions "app/backend/controller/yce/extensions"
	myhealthz "app/backend/controller/yce/healthz"
	mylogin "app/backend/controller/yce/login"
	mylogout "app/backend/controller/yce/logout"
	mynamespace "app/backend/controller/yce/namespace"
	mynavList "app/backend/controller/yce/navlist"
	myregistry "app/backend/controller/yce/registry"
	myservice "app/backend/controller/yce/service"
	mytopology "app/backend/controller/yce/topology"
	myversion "app/backend/controller/yce/version"
	"github.com/kataras/iris"
	myoperationstat "app/backend/controller/yce/dashboard/operationstat"
)

type Router struct {
	Login            *mylogin.LoginController
	Logout           *mylogout.LogoutController
	Nav              *mynavList.NavListController
	ListDeploy       *mydeploy.ListDeploymentController
	Registry         *myregistry.ListRegistryController
	ListService      *myservice.ListServiceController
	InitService      *myservice.InitServiceController
	CreateService    *myservice.CreateServiceController
	ListEndpoints    *myendpoint.ListEndpointsController
	InitEndpoints    *myendpoint.InitEndpointsController
	CreateEndpoints  *myendpoint.CreateEndpointsController
	ListExtensions   *myextensions.ListExtensionsController
	InitDeploy       *mydeploy.InitDeploymentController
	CreateDeploy     *mydeploy.CreateDeploymentController
	RollingDeploy    *mydeploy.RollingDeploymentController
	RollbackDeploy   *mydeploy.RollbackDeploymentController
	ScaleDeploy      *mydeploy.ScaleDeploymentController
	DeleteDeploy     *mydeploy.DeleteDeploymentController
	LogsPod          *mydeploy.LogsPodController
	ListOperationLog *mydeploy.ListOperationLogController
	InitNamespace    *mynamespace.InitNamespaceController
	DeleteService    *myservice.DeleteServiceController
	DeleteEndpoint   *myendpoint.DeleteEndpointsController
	HistoryDeploy    *mydeploy.HistoryDeploymentController
	Topology         *mytopology.TopologyController
	Api              *mypath.ApisController
	Healthz          *myhealthz.HealthzController
	Version          *myversion.VersionController
	StatDeployment   *mydeploymentstat.StatDeploymentController
	StatResource     *myresourcestat.StatResourceController
	OperationStat    *myoperationstat.OperationStatController
}


func NewRouter() *Router {
	r := new(Router)
	r.Login = new(mylogin.LoginController)
	r.Logout = new(mylogout.LogoutController)
	r.Nav = new(mynavList.NavListController)
	r.ListDeploy = new(mydeploy.ListDeploymentController)
	r.Registry = new(myregistry.ListRegistryController)
	r.ListService = new(myservice.ListServiceController)
	r.InitService = new(myservice.InitServiceController)
	r.CreateService = new(myservice.CreateServiceController)
	r.ListEndpoints = new(myendpoint.ListEndpointsController)
	r.InitEndpoints = new(myendpoint.InitEndpointsController)
	r.CreateEndpoints = new(myendpoint.CreateEndpointsController)
	r.ListExtensions = new(myextensions.ListExtensionsController)
	r.InitDeploy = new(mydeploy.InitDeploymentController)
	r.CreateDeploy = new(mydeploy.CreateDeploymentController)
	r.RollingDeploy = new(mydeploy.RollingDeploymentController)
	r.RollbackDeploy = new(mydeploy.RollbackDeploymentController)
	r.ScaleDeploy = new(mydeploy.ScaleDeploymentController)
	r.DeleteDeploy = new(mydeploy.DeleteDeploymentController)
	r.LogsPod = new(mydeploy.LogsPodController)
	r.ListOperationLog = new(mydeploy.ListOperationLogController)
	r.InitNamespace = new(mynamespace.InitNamespaceController)
	r.DeleteService = new(myservice.DeleteServiceController)
	r.DeleteEndpoint = new(myendpoint.DeleteEndpointsController)
	r.HistoryDeploy = new(mydeploy.HistoryDeploymentController)
	r.Topology = new(mytopology.TopologyController)
	r.Api = new(mypath.ApisController)
	r.Healthz = new(myhealthz.HealthzController)
	r.Version = new(myversion.VersionController)
	r.StatDeployment = new(mydeploymentstat.StatDeploymentController)
	r.StatResource = new(myresourcestat.StatResourceController)
	r.OperationStat = new(myoperationstat.OperationStatController)

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
	//iris.API("/api/v1/organizations/:orgId/deployments/:deploymentName", *r.DeleteDeploy)
	iris.API("/api/v1/organizations/:orgId/pods/:podName/logs", *r.LogsPod)
	iris.API("/api/v1/organizations/:orgId/deployments/:deploymentName/delete", *r.DeleteDeploy)
	iris.API("/api/v1/organizations/:orgId/operationlog", *r.ListOperationLog)
	iris.API("/api/v1/registry/images", *r.Registry)
	iris.API("/api/v1/organizations/:orgId/users/:userId/services", *r.ListService)
	iris.API("/api/v1/organizations/:orgId/users/:userId/services/init", *r.InitService)
	iris.API("/api/v1/organizations/:orgId/users/:userId/services/new", *r.CreateService)
	iris.API("/api/v1/organizations/:orgId/services/:svcName", *r.DeleteService)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints", *r.ListEndpoints)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints/init", *r.InitEndpoints)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints/new", *r.CreateEndpoints)
	iris.API("/api/v1/organizations/:orgId/endpoints/:epName", *r.DeleteEndpoint)
	iris.API("/api/v1/organizations/:orgId/users/:userId/extensions", *r.ListExtensions)
	iris.API("/api/v1/organizations/init", *r.InitNamespace)
	iris.API("/api/v1/organizations/:orgId/datacenters/:dcId/deployments/:name/history", *r.HistoryDeploy)
	iris.API("/api/v1/organizations/:orgId/topology", *r.Topology)
	iris.API("/", *r.Api)
	iris.API("/version", *r.Version)
	iris.API("/healthz", *r.Healthz)
	iris.API("/api/v1/organizations/:orgId/deploymentstat", *r.StatDeployment)
	iris.API("/api/v1/organizations/:orgId/resourcestat", *r.StatResource)
	iris.API("/api/v1/organizations/:orgId/operationstat", *r.OperationStat)

	iris.StaticServe("../frontend", "/static")
}
