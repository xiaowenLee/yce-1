package template

import (
	"app/backend/controller/yce"
	myerror "app/backend/common/yce/error"

	//"k8s.io/kubernetes/pkg/apis/extensions"
	//"k8s.io/kubernetes/pkg/api"
	mytemplate "app/backend/model/mysql/template"
	mydeployment "app/backend/model/yce/deploy"
	myservice "app/backend/model/yce/service"
	"encoding/json"
	"strconv"
)

type CreateTemplateController struct {
	yce.Controller

	params *CreateTemplateParams
}

type CreateTemplateParams struct {
	Name string `json:"name"`
	Deployment *mydeployment.CreateDeployment `json:"deployment",omitempty`
	Service  *myservice.CreateService `json:"service",omitempty`
}

func (ctc *CreateTemplateController)createDbItem(OrgId, Op int32) {

	deployment, err := json.Marshal(ctc.params.Deployment)
	if err != nil {
		ctc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	service, err := json.Marshal(ctc.params.Service)
	if err != nil {
		ctc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}

	t := mytemplate.NewTemplate(ctc.params.Name, OrgId, string(deployment), string(service), "", Op, "")
	err = t.InsertTemplate(Op)
	if err != nil {
		ctc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}

	return
}

func (ctc CreateTemplateController) Post() {
	SessionIdFromClient := ctc.RequestHeader("Authorization")

	orgId := ctc.Param("orgId")
	userId := ctc.Param("userId")

	ctc.ValidateSession(SessionIdFromClient, orgId)
	if ctc.CheckError() {
		return
	}

	ctc.params = new(CreateTemplateParams)
	//ctc.params.Deployment = new(extensions.Deployment)
	//ctc.params.Service = new(api.Service)

	err := ctc.ReadJSON(ctc.params)
	if err != nil {
		ctc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if ctc.CheckError() {
		return
	}

	oId, err := strconv.Atoi(orgId)
	if err != nil {
		ctc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	uId, err := strconv.Atoi(userId)
	if err != nil {
		ctc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	if ctc.CheckError() {
		return
	}

	ctc.createDbItem(int32(oId), int32(uId))
	if ctc.CheckError() {
		return
	}

	ctc.WriteOk("")
	log.Infof("CreateTemplateController Success")
	return
}


