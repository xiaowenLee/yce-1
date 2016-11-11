package template

import (
	myerror "app/backend/common/yce/error"
	"app/backend/controller/yce"

	mytemplate "app/backend/model/mysql/template"
	"encoding/json"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"strconv"
)

type UpdateTemplateController struct {
	yce.Controller

	params *UpdateTemplateParams
}

type UpdateTemplateParams struct {
	Name       string                 `json:"name"`
	Deployment *extensions.Deployment `json:"deployment"`
	Service    *api.Service           `json:"service"`
}

func (utc *UpdateTemplateController) updateDbItem(OrgId, Op int32) {

	deployment, err := json.Marshal(utc.params.Deployment)
	if err != nil {
		utc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	service, err := json.Marshal(utc.params.Service)
	if err != nil {
		utc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}

	t := mytemplate.NewTemplate(utc.params.Name, OrgId, string(deployment), string(service), "", Op, "")
	err = t.UpdateTemplate(Op)
	if err != nil {
		utc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}

	return
}

func (utc UpdateTemplateController) Post() {
	SessionIdFromClient := utc.RequestHeader("Authorization")

	orgId := utc.Param("orgId")
	userId := utc.Param("userId")

	utc.ValidateSession(SessionIdFromClient, orgId)
	if utc.CheckError() {
		return
	}

	utc.params = new(UpdateTemplateParams)
	utc.params.Deployment = new(extensions.Deployment)
	utc.params.Service = new(api.Service)

	err := utc.ReadJSON(utc.params)
	if err != nil {
		utc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if utc.CheckError() {
		return
	}

	oId, err := strconv.Atoi(orgId)
	if err != nil {
		utc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	uId, err := strconv.Atoi(userId)
	if err != nil {
		utc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	if utc.CheckError() {
		return
	}

	utc.updateDbItem(int32(oId), int32(uId))
	if utc.CheckError() {
		return
	}

	utc.WriteOk("")
	log.Infof("UpdateTemplateController Success")
	return
}
