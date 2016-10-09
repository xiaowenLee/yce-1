package deployment

import (
	yce "app/backend/controller/yce"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	yceutils "app/backend/controller/yce/utils"
)

type DeploymentStatController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
}


func (dsc DeploymentStatController) Get() {
	sessionIdFromClient := d
}


