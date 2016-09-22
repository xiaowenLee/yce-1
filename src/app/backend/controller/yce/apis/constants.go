package apis

const (
	APIS = `
	{
		path: [
			"path",
			"version",
			"healthz",
			"/api/v1/users/login",
			"/api/v1/users/logout",
			"/api/v1/organizations/:orgId/users/:userId/deployments",
			"/api/v1/organizations/:orgId/users/:userId/deployments/init",
			"/api/v1/organizations/:orgId/users/:userId/deployments/new",
			"/api/v1/organizations/:orgId/deployments/:deploymentName/rolling",
			"/api/v1/organizations/:orgId/deployments/:deploymentName/rollback",
			"/api/v1/organizations/:orgId/deployments/:deploymentName/scale",
			"/api/v1/organizations/:orgId/pods/:podName/logs",
			"/api/v1/organizations/:orgId/deployments/:deploymentName/delete",
			"/api/v1/organizations/:orgId/operationlog",
			"/api/v1/registry/images",
			"/api/v1/organizations/:orgId/users/:userId/services",
			"/api/v1/organizations/:orgId/users/:userId/services/init",
			"/api/v1/organizations/:orgId/users/:userId/services/new",
			"/api/v1/organizations/:orgId/services/:svcName",
			"/api/v1/organizations/:orgId/users/:userId/endpoints",
			"/api/v1/organizations/:orgId/users/:userId/endpoints/init",
			"/api/v1/organizations/:orgId/users/:userId/endpoints/new",
			"/api/v1/organizations/:orgId/endpoints/:epName",
			"/api/v1/organizations/:orgId/users/:userId/extensions",
			"/api/v1/organizations/init",
			"/api/v1/organizations/:orgId/datacenters/:dcId/deployments/:name/history",
			 "/api/v1/organizations/:orgId/topology",
			 "/static"
		]
	}
	`
)
