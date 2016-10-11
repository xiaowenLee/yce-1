package user

import (
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	"encoding/json"
	myerror "app/backend/common/yce/error"
)

type InitUserController struct {
	yce.Controller
}

type OrgNames struct {
	OrgNameList []string `json:"orgNameList"`
}

// get Organization Names
func (iuc *InitUserController) getOrgNames() string {
	orgNames := new(OrgNames)
	orgNames.OrgNameList, iuc.Ye = yceutils.GetOrgNameList()
	if iuc.Ye != nil {
		return ""
	}

	orgNamesJSON, err := json.Marshal(orgNames.OrgNameList)
	if err != nil {
		iuc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	orgNamesString := string(orgNamesJSON)

	return orgNamesString
}

func (iuc InitUserController) Get() {
	//TODO: rethink of session authroization. Here it is omitted.
	// SessionIdFromClient := iuc.RequestHeaders("Authrozation")

	orgNames := iuc.getOrgNames()
	if iuc.CheckError() {
		return
	}

	iuc.WriteOk(orgNames)
	log.Infoln("InitUserController Get Over")

}
