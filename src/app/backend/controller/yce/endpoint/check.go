package endpoint

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/api"
)

type CheckEndpointsController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string

	params *CheckEndpointsParams
	orgName string
}

type CheckEndpointsParams struct {
	Name string `json:"name"`
}

func (cec *CheckEndpointsController) check(c *client.Client, name string) {
	epList, err := c.Endpoints(cec.orgName).List(api.ListOptions{})
	if err != nil {
		cec.Ye = myerror.NewYceError(myerror.EKUBE_LIST_ENDPOINTS, "")
		return
	}

	for _, ep := range epList.Items {
		if name == ep.Name {
			cec.Ye = myerror.NewYceError(myerror.EYCE_EXISTED_NAME, "")
			return
		}
	}
}

func (cec *CheckEndpointsController) checkDuplicatedName() {
	for _, c := range cec.k8sClients {
		cec.check(c, cec.params.Name)
		if cec.Ye != nil {
			return
		}
	}
}

func (cec CheckEndpointsController)Post() {
	cec.params = new(CheckEndpointsParams)

	sessionIdFromClient := cec.RequestHeader("Authorization")
	orgId := cec.Param("orgId")

	cec.ValidateSession(sessionIdFromClient, orgId)
	if cec.CheckError() {
		return
	}

	err := cec.ReadJSON(cec.params)
	if err != nil {
		cec.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cec.CheckError() {
		return
	}

	cec.orgName, cec.Ye = yceutils.GetOrgNameByOrgId(orgId)
	if cec.CheckError() {
		return
	}

	dcIdList, ye := yceutils.GetDcIdListByOrgId(orgId)
	if ye != nil {
		cec.Ye = ye
	}
	if cec.CheckError() {
		return
	}

	cec.apiServers, cec.Ye = yceutils.GetApiServerList(dcIdList)
	if cec.CheckError() {
		return
	}

	cec.k8sClients, cec.Ye = yceutils.CreateK8sClientList(cec.apiServers)
	if cec.CheckError() {
		return
	}

	cec.checkDuplicatedName()
	if cec.CheckError() {
		return
	}

	cec.WriteOk("")
	log.Infoln("CheckEndpointsController Check Over!")
}
