package template
import (
	"app/backend/controller/yce"
	myerror "app/backend/common/yce/error"

	mytemplate "app/backend/model/mysql/template"
	"encoding/json"
	"strconv"
)

type DeleteTemplateController struct {
	yce.Controller

	params *DeleteTemplateParams
}

type DeleteTemplateParams struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}

func (dtc *DeleteTemplateController)deleteDbItem(OrgId, Op int32) {

	t := new(mytemplate.Template)
	err := t.QueryTemplateById(dtc.params.Id)
	if err != nil {
		dtc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	err = t.DeleteTemplate(Op)
	if err != nil {
		dtc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return
	}

	return
}

func (dtc DeleteTemplateController) Post() {
	SessionIdFromClient := dtc.RequestHeader("Authorization")

	orgId := dtc.Param("orgId")
	userId := dtc.Param("userId")

	dtc.ValidateSession(SessionIdFromClient, orgId)
	if dtc.CheckError() {
		return
	}

	dtc.params = new(DeleteTemplateParams)

	err := dtc.ReadJSON(dtc.params)
	if err != nil {
		dtc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if dtc.CheckError() {
		return
	}

	oId, err := strconv.Atoi(orgId)
	if err != nil {
		dtc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	uId, err := strconv.Atoi(userId)
	if err != nil {
		dtc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	if dtc.CheckError() {
		return
	}

	dtc.deleteDbItem(int32(oId), int32(uId))
	if dtc.CheckError() {
		return
	}

	dtc.WriteOk("")
	log.Infof("DeleteTemplateController Success")
	return
}
