package nodeport

import (
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type CheckNodePortController struct {
	yce.Controller

	params *CheckNodePortParams
}

type CheckNodePortParams struct {
	Port    int32 `json:"nodePort"`
	DcIdList []int32 `json:"dcIdList"`
	OrgId   string `json:"orgId"`
}

func (cnc *CheckNodePortController) checkDuplicatedNodePort() {
	for _, dcId := range cnc.params.DcIdList {
		log.Infof("CheckDuplicatedNodePort: port=%d, dcId=%d", cnc.params.Port, dcId)
		_, ye := yceutils.QueryDuplicatedNodePort(cnc.params.Port, dcId)
		if ye != nil {
			cnc.Ye = ye
			return
		}
	}
	return
}

func (cnc CheckNodePortController) Post() {
	cnc.params = new(CheckNodePortParams)

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

	cnc.checkDuplicatedNodePort()
	if cnc.CheckError() {
		return
	}

	cnc.WriteOk("")
	log.Infoln("CheckNodePortController Post Over!")
}
