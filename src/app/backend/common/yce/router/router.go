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
	mynavList "app/backend/controller/yce/navlist"
	myregistry "app/backend/controller/yce/registry"
	myservice "app/backend/controller/yce/service"
	mytopology "app/backend/controller/yce/topology"
	myversion "app/backend/controller/yce/version"
	myuser "app/backend/controller/yce/user"
	mynamespace "app/backend/controller/yce/namespace"
	"github.com/kataras/iris"
	myoperationstat "app/backend/controller/yce/dashboard/operationstat"
	mydatacenter "app/backend/controller/yce/datacenter"
	mynodeport "app/backend/controller/yce/nodeport"
	mytemplate "app/backend/controller/yce/template"
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
	CheckService     *myservice.CheckServiceController
	ListEndpoints    *myendpoint.ListEndpointsController
	InitEndpoints    *myendpoint.InitEndpointsController
	CreateEndpoints  *myendpoint.CreateEndpointsController
	CheckEndpoints   *myendpoint.CheckEndpointsController
	ListExtensions   *myextensions.ListExtensionsController
	InitDeploy       *mydeploy.InitDeploymentController
	CreateDeploy     *mydeploy.CreateDeploymentController
	RollingDeploy    *mydeploy.RollingDeploymentController
	RollbackDeploy   *mydeploy.RollbackDeploymentController
	ScaleDeploy      *mydeploy.ScaleDeploymentController
	DeleteDeploy     *mydeploy.DeleteDeploymentController
	CheckDeploy      *mydeploy.CheckDeploymentController
	LogsPod          *mydeploy.LogsPodController
	ListOperationLog *mydeploy.ListOperationLogController
	DeleteService    *myservice.DeleteServiceController
	DeleteEndpoint   *myendpoint.DeleteEndpointsController
	HistoryDeploy    *mydeploy.HistoryDeploymentController
	Topology         *mytopology.TopologyController
	Api              *mypath.ApisController
	Healthz          *myhealthz.HealthzController
	Version          *myversion.VersionController
	StatDeployment   *mydeploymentstat.StatDeploymentController
	StatResource     *myresourcestat.StatResourceController
	InitUser 	 *myuser.InitUserController
	CreateUser 	 *myuser.CreateUserController
	CheckUser 	 *myuser.CheckUserController
	ListUser         *myuser.ListUserController
	DeleteUser	 *myuser.DeleteUserController
	UpdateUser       *myuser.UpdateUserController
	CheckNamespace	 *mynamespace.CheckNamespaceController
	InitNamespace    *mynamespace.InitNamespaceController
	ListNamespace    *mynamespace.ListNamespaceController
	CreateNamespace  *mynamespace.CreateNamespaceController
	UpdateNamespace  *mynamespace.UpdateNamespaceController
	DeleteNamespace  *mynamespace.DeleteNamespaceController
	OperationStat    *myoperationstat.OperationStatController
	CreateDatacenter *mydatacenter.CreateDatacenterController
	ListDatacenter   *mydatacenter.ListDatacenterController
	DeleteDatacenter *mydatacenter.DeleteDatacenterController
	UpdateDatacenter *mydatacenter.UpdateDatacenterController
	CheckDatacenter  *mydatacenter.CheckDatacenterController
	ListNodePort 	 *mynodeport.ListNodePortController
	CheckNodePort    *mynodeport.CheckNodePortController
	CheckTemplate    *mytemplate.CheckTemplateController
	CreateTemplate   *mytemplate.CreateTemplateController
	UpdateTemplate   *mytemplate.UpdateTemplateController
	DeleteTemplate   *mytemplate.DeleteTemplateController
	ListTemplate     *mytemplate.ListTemplateController
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
	r.CheckService = new(myservice.CheckServiceController)
	r.ListEndpoints = new(myendpoint.ListEndpointsController)
	r.InitEndpoints = new(myendpoint.InitEndpointsController)
	r.CreateEndpoints = new(myendpoint.CreateEndpointsController)
	r.CheckEndpoints = new(myendpoint.CheckEndpointsController)
	r.ListExtensions = new(myextensions.ListExtensionsController)
	r.InitDeploy = new(mydeploy.InitDeploymentController)
	r.CreateDeploy = new(mydeploy.CreateDeploymentController)
	r.RollingDeploy = new(mydeploy.RollingDeploymentController)
	r.RollbackDeploy = new(mydeploy.RollbackDeploymentController)
	r.ScaleDeploy = new(mydeploy.ScaleDeploymentController)
	r.DeleteDeploy = new(mydeploy.DeleteDeploymentController)
	r.CheckDeploy = new(mydeploy.CheckDeploymentController)
	r.LogsPod = new(mydeploy.LogsPodController)
	r.ListOperationLog = new(mydeploy.ListOperationLogController)
	r.DeleteService = new(myservice.DeleteServiceController)
	r.DeleteEndpoint = new(myendpoint.DeleteEndpointsController)
	r.HistoryDeploy = new(mydeploy.HistoryDeploymentController)
	r.Topology = new(mytopology.TopologyController)
	r.Api = new(mypath.ApisController)
	r.Healthz = new(myhealthz.HealthzController)
	r.Version = new(myversion.VersionController)
	r.StatDeployment = new(mydeploymentstat.StatDeploymentController)
	r.StatResource = new(myresourcestat.StatResourceController)
	r.InitUser = new(myuser.InitUserController)
	r.CheckUser = new(myuser.CheckUserController)
	r.CreateUser = new(myuser.CreateUserController)
	r.ListUser = new(myuser.ListUserController)
	r.DeleteUser = new(myuser.DeleteUserController)
	r.UpdateUser = new(myuser.UpdateUserController)
	r.CheckNamespace = new(mynamespace.CheckNamespaceController)
	r.InitNamespace = new(mynamespace.InitNamespaceController)
	r.ListNamespace = new(mynamespace.ListNamespaceController)
	r.CreateNamespace = new(mynamespace.CreateNamespaceController)
	r.UpdateNamespace = new(mynamespace.UpdateNamespaceController)
	r.DeleteNamespace = new(mynamespace.DeleteNamespaceController)
	r.OperationStat = new(myoperationstat.OperationStatController)
	r.CreateDatacenter = new(mydatacenter.CreateDatacenterController)
	r.DeleteDatacenter = new(mydatacenter.DeleteDatacenterController)
	r.ListDatacenter = new(mydatacenter.ListDatacenterController)
	r.UpdateDatacenter = new(mydatacenter.UpdateDatacenterController)
	r.CheckDatacenter = new(mydatacenter.CheckDatacenterController)
	r.ListNodePort = new(mynodeport.ListNodePortController)
	r.CheckNodePort = new(mynodeport.CheckNodePortController)
	r.CheckTemplate = new(mytemplate.CheckTemplateController)
	r.CreateTemplate = new(mytemplate.CreateTemplateController)
	r.UpdateTemplate = new(mytemplate.UpdateTemplateController)
	r.DeleteTemplate = new(mytemplate.DeleteTemplateController)
	r.ListTemplate = new(mytemplate.ListTemplateController)

	return r
}

func (r *Router) Registe() {

	iris.API("/api/v1/users/login", *r.Login)
	iris.API("/api/v1/users/logout", *r.Logout)
	iris.API("/api/v1/organizations/:orgId/users/:userId/navList", *r.Nav)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments", *r.ListDeploy)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments/init", *r.InitDeploy)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments/new", *r.CreateDeploy)
	iris.API("/api/v1/organizations/:orgId/users/:userId/deployments/check", *r.CheckDeploy)
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
	iris.API("/api/v1/organizations/:orgId/users/:userId/services/check", *r.CheckService)
	iris.API("/api/v1/organizations/:orgId/services/:svcName", *r.DeleteService)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints", *r.ListEndpoints)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints/init", *r.InitEndpoints)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints/new", *r.CreateEndpoints)
	iris.API("/api/v1/organizations/:orgId/users/:userId/endpoints/check", *r.CheckEndpoints)
	iris.API("/api/v1/organizations/:orgId/endpoints/:epName", *r.DeleteEndpoint)
	iris.API("/api/v1/organizations/:orgId/users/:userId/extensions", *r.ListExtensions)
	iris.API("/api/v1/organizations/:orgId/datacenters/:dcId/deployments/:name/history", *r.HistoryDeploy)
	iris.API("/api/v1/organizations/:orgId/users/:userId/topology", *r.Topology)
	iris.API("/", *r.Api)
	iris.API("/version", *r.Version)
	iris.API("/healthz", *r.Healthz)
	iris.API("/api/v1/organizations/:orgId/deploymentstat", *r.StatDeployment)
	iris.API("/api/v1/organizations/:orgId/resourcestat", *r.StatResource)
	iris.API("/api/v1/user/init", *r.InitUser)
	iris.API("/api/v1/user/check", *r.CheckUser)
	iris.API("/api/v1/user/new", *r.CreateUser)
	iris.API("/api/v1/user", *r.ListUser)
	iris.API("/api/v1/user/delete", *r.DeleteUser)
	iris.API("/api/v1/user/update", *r.UpdateUser)
	iris.API("/api/v1/organization/check", *r.CheckNamespace)
	iris.API("/api/v1/organization/init", *r.InitNamespace)
	iris.API("/api/v1/organization", *r.ListNamespace)
	iris.API("/api/v1/organization/new", *r.CreateNamespace)
	iris.API("/api/v1/organization/update", *r.UpdateNamespace)
	iris.API("/api/v1/organization/delete", *r.DeleteNamespace)
	iris.API("/api/v1/organizations/:orgId/operationstat", *r.OperationStat)
	iris.API("/api/v1/datacenter/new", *r.CreateDatacenter)
	iris.API("/api/v1/datacenter/delete", *r.DeleteDatacenter)
	iris.API("/api/v1/datacenter", *r.ListDatacenter)
	iris.API("/api/v1/datacenter/update", *r.UpdateDatacenter)
	iris.API("/api/v1/datacenter/check", *r.CheckDatacenter)
	iris.API("/api/v1/nodeports", *r.ListNodePort)
	iris.API("/api/v1/nodeports/check", *r.CheckNodePort)
	iris.API("/api/v1/organizations/:orgId/users/:userId/templates/check", *r.CheckTemplate)
	iris.API("/api/v1/organizations/:orgId/users/:userId/templates/new", *r.CreateTemplate)
	iris.API("/api/v1/organizations/:orgId/users/:userId/templates/update", *r.UpdateTemplate)
	iris.API("/api/v1/organizations/:orgId/users/:userId/templates/delete", *r.DeleteTemplate)
	iris.API("/api/v1/organizations/:orgId/users/:userId/templates", *r.ListTemplate)


	iris.StaticServe("../frontend", "/static")
}
