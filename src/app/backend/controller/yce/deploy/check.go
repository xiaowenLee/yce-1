package deploy

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/api"
)

type CheckDeploymentController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string

	params  *CheckDeploymentParams
	orgName string
}

type CheckDeploymentParams struct {
	Name string `json:"name"`
}

func (cdc *CheckDeploymentController) check(c *client.Client, name string) {
	dpList, err := c.Deployments(cdc.orgName).List(api.ListOptions{})
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_DEPLOYMENTS, "")
		return
	}

	for _, dp := range dpList.Items {
		if name == dp.Name {
			cdc.Ye = myerror.NewYceError(myerror.EYCE_EXISTED_NAME, "")
			return
		}
	}
}

func (cdc *CheckDeploymentController) checkDuplicatedName() {
	for _, c := range cdc.k8sClients {
		cdc.check(c, cdc.params.Name)
		if cdc.Ye != nil {
			return
		}
	}
}

func (cdc CheckDeploymentController) Post() {
	cdc.params = new(CheckDeploymentParams)

	sessionIdFromClient := cdc.RequestHeader("Authorization")
	orgId := cdc.Param("orgId")

	cdc.ValidateSession(sessionIdFromClient, orgId)
	if cdc.CheckError() {
		return
	}

	err := cdc.ReadJSON(cdc.params)
	if err != nil {
		cdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cdc.CheckError() {
		return
	}

	cdc.orgName, cdc.Ye = yceutils.GetOrgNameByOrgId(orgId)
	if cdc.CheckError() {
		return
	}

	dcIdList, ye := yceutils.GetDcIdListByOrgId(orgId)
	if ye != nil {
		cdc.Ye = ye
	}
	if cdc.CheckError() {
		return
	}

	cdc.apiServers, cdc.Ye = yceutils.GetApiServerList(dcIdList)
	if cdc.CheckError() {
		return
	}

	cdc.k8sClients, cdc.Ye = yceutils.CreateK8sClientList(cdc.apiServers)
	if cdc.CheckError() {
		return
	}

	cdc.checkDuplicatedName()
	if cdc.CheckError() {
		return
	}

	cdc.WriteOk("")
	log.Infoln("CheckDeploymentController Check Over!")
}
