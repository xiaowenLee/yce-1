package namespace

import (
	yce "app/backend/controller/yce"
	myerror "app/backend/common/yce/error"
	yceutils "app/backend/controller/yce/utils"
	"github.com/kubernetes/kubernetes/pkg/util/json"
	"strconv"
)

type CheckNamespaceController struct {
	yce.Controller

	params *CheckNamespaceParams
}

type CheckNamespaceParams struct {
	OrgName string `json:"orgName"`
	OrgId   string `json:"orgId"`
}

func (cnc *CheckNamespaceController) checkDuplicatedName() string {
	org, ye := yceutils.QueryDuplicatedOrgName(cnc.params.OrgName)
	// not found
	if ye != nil {
		return ""
	}

	// found
	orgId := strconv.Itoa(int(org.Id))
	dcList, ye := yceutils.GetDatacenterListByOrgId(orgId)
	// found but error
	if ye != nil {
		cnc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	dcListJSON, err := json.Marshal(dcList)
	if err != nil {
		cnc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	dcNameString := string(dcListJSON)
	return dcNameString
}



func (cnc CheckNamespaceController) Post() {
	cnc.params = new(CheckNamespaceParams)


	// TODO rethink of this authentication way
	sessionIdFromClient := cnc.RequestHeader("Authorization")

	err := cnc.ReadJSON(cnc.params)
	if err != nil {
		cnc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}
	if cnc.CheckError() {
		return
	}

	// validate admin's session
	cnc.ValidateSession(sessionIdFromClient, cnc.params.OrgId)
	if cnc.CheckError() {
		return
	}



	//checkDuplicatedName
	checkResult := cnc.checkDuplicatedName()
	if cnc.CheckError() {
		return
	}

	cnc.WriteOk(checkResult)
	log.Infoln("CheckNamespaceController Post Over!")
}
