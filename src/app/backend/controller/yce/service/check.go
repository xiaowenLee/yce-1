package service

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/api"
)

type CheckServiceController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string

	params *CheckServiceParams
	orgName string
}

type CheckServiceParams struct {
	Name string `json:"name"`
}

func (csc *CheckServiceController) check(c *client.Client, name string) {
	svcList, err := c.Services(csc.orgName).List(api.ListOptions{})
	if err != nil {
		csc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_SERVICE, "")
		return
	}

	for _, svc := range svcList.Items {
		if name == svc.Name {
			csc.Ye = myerror.NewYceError(myerror.EYCE_EXISTED_NAME, "")
			return
		}
	}
}

func (csc *CheckServiceController) checkDuplicatedName() {
	for _, c := range csc.k8sClients {
		csc.check(c, csc.params.Name)
		if csc.Ye != nil {
			return
		}
	}
}

func (csc CheckServiceController) Post() {
	csc.params = new(CheckServiceParams)
	sessionIdFromClient := csc.RequestHeader("Authorization")
	orgId := csc.Param("orgId")

	csc.ValidateSession(sessionIdFromClient, orgId)
	if csc.CheckError() {
		return
	}

	err := csc.ReadJSON(csc.params)
	if err != nil {
		csc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if csc.CheckError() {
		return
	}

	csc.orgName, csc.Ye = yceutils.GetOrgNameByOrgId(orgId)
	if csc.CheckError() {
		return
	}

	dcIdList, ye := yceutils.GetDcIdListByOrgId(orgId)
	if ye != nil {
		csc.Ye = ye
	}
	if csc.CheckError() {
		return
	}

	csc.apiServers, csc.Ye = yceutils.GetApiServerList(dcIdList)
	if csc.CheckError() {
		return
	}

	csc.k8sClients, csc.Ye = yceutils.CreateK8sClientList(csc.apiServers)
	if csc.CheckError() {
		return
	}

	csc.checkDuplicatedName()
	if csc.CheckError() {
		return
	}

	csc.WriteOk("")
	log.Infoln("CheckServiceController Check Over!")
}
