package deploy

import (
	myoption "app/backend/model/mysql/option"
)

const (

	// Create a deployment
	CREATE_TYPE  = myoption.ONLINE
	CREATE_VERBE = "POST"
	CREATE_URL   = "/api/v1/organization/<orgId>/users/<userId>/deployments"

	// Delete a deployment
	DELETE_TYPE  = myoption.DELETE
	DELETE_VERBE = "DELETE"
	DELETE_URL   = "/api/v1/organization/<orgId>/deployments/<deploymentName>"

	// List the history of a deployment
	REVISION_ANNOTATION string = "deployment.kubernetes.io/revision"

	// List operation log
	SELECT_DEPLOYMENT = "SELECT id, name, actionType, actionVerb, actionUrl, actionAt, actionOp, dcList, success, reason, json, comment FROM deployment WHERE orgId=? ORDER BY id DESC LIMIT 30"
	SELECT_USER       = "SELECT name FROM user WHERE id=?"
	SELECT_DATACENTER = "SELECT name FROM datacenter WHERE id=?"

	// Rollback a deployment
	ROLLBACK_ACTION_TYPE                = myoption.ROLLINGBACK
	ROLLBACK_ACTION_VERBE               = "POST"
	ROLLBACK_ACTION_URL                 = "/api/v1/organizations/<orgId>/deployments/<name>/rollback"
	ROLLBACK_REVISION_ANNOTATION string = "deployment.kubernetes.io/revision"
	ROLLBACK_IMAGE                      = "image"
	ROLLBACK_USERID                     = "userId"
	ROLLBACK_CHANGE_CAUSE        string = "kubernetes.io/change-cause"

	// Rolling update a deployment
	ROLLING_TYPE           = myoption.ROLLINGUPGRADE
	ROLLING_VERBE          = "POST"
	ROLLING_URL            = "/api/v1/organization/<orgId>/deployments/<deploymentName>/rolling"
	ROLLING_MAXUNAVAILABLE = 2
	ROLLING_MAXSURGE       = 2

	// Scale a deployment
	SCALE_ACTION_TYPE  = myoption.SCALING
	SCALE_ACTION_VERBE = "POST"
	SCALE_ACTION_URL   = "/api/v1/organizations/<orgId>/datacenters/<dcId>/deployments/<name>/scale"
)
