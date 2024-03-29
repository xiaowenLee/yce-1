package template


import (
	"app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	"strconv"
	mytemplate "app/backend/model/mysql/template"
	"github.com/kubernetes/kubernetes/pkg/util/json"
	yceutils "app/backend/controller/yce/utils"
)

type ListTemplateController struct {
	yce.Controller

	params *ListTemplateParams
}

type ListTemplateParams struct {
	Templates []mytemplate.Template `json:"templates"`
	Users     map[int32]string `json:"users"`
}

func (ltc *ListTemplateController) getTemplateList(orgId int32) string {
	templates, err := mytemplate.QueryAllTemplatesByOrgId(orgId)
	if err != nil {
		ltc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return ""
	}
	ltc.params.Templates = templates

	ltc.params.Users, ltc.Ye = yceutils.GetUsersByOrgId(orgId)
	if ltc.Ye != nil {
		return ""
	}

	templateJSON, err := json.Marshal(ltc.params)
	if err != nil {
		ltc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	templateString := string(templateJSON)
	return templateString

}

func (ltc ListTemplateController) Get() {
	SessionIdFromClient := ltc.RequestHeader("Authorization")

	orgId := ltc.Param("orgId")
	//userId := ltc.Param("userId")

	ltc.ValidateSession(SessionIdFromClient, orgId)
	if ltc.CheckError() {
		return
	}

	ltc.params = new(ListTemplateParams)

	oId, err := strconv.Atoi(orgId)
	if err != nil {
		ltc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}
	if ltc.CheckError() {
		return
	}

	data := ltc.getTemplateList(int32(oId))
	if ltc.CheckError() {
		return
	}

	ltc.WriteOk(data)
	log.Infof("ListTemplateController check Ok: len(templates)=%d", len(ltc.params.Templates))
	return
}
