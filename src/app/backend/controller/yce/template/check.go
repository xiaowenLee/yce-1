package template

import (
	yceutils "app/backend/controller/yce/utils"
	"app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	"strconv"
)

type CheckTemplateController struct {
	yce.Controller

	params *CheckTemplateParams
}

type CheckTemplateParams struct {
	Name string `json:"name"`
}

func (ctc *CheckTemplateController) CheckDuplicatedName(orgId int32) {
	t, ye := yceutils.QueryDuplicatedTemplateName(ctc.params.Name, orgId)
	// not found
	if t == nil && ye != nil {
		return
	}

	if t != nil && ye != nil {
		ctc.Ye = ye
		return
	}

	if t != nil && ye == nil {
		return
	}

	// TODO: rewrite logical
}

func (ctc CheckTemplateController) Post() {
	SessionIdFromClient := ctc.RequestHeader("Authorization")

	orgId := ctc.Param("orgId")
	//userId := ctc.Param("userId")

	ctc.ValidateSession(SessionIdFromClient, orgId)
	if ctc.CheckError() {
		return
	}

	params := new(CheckTemplateParams)

	err := ctc.ReadJSON(params)
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
	if ctc.CheckError() {
		return
	}

	ctc.CheckDuplicatedName(int32(oId))
	if ctc.CheckError() {
		return
	}

	ctc.WriteOk("")
	log.Infof("CheckTemplateController check Ok: name=%s, orgId=%d", ctc.params.Name, oId)
	return
}
