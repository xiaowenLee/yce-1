package utils

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log
type OrgIdAndNameType struct {
	OrgId int32 `json:"orgId"`
	OrgName string `json:"orgName"`
}

type DcIdAndNameType struct {
	DcId int32 `json:"dcId"`
	DcName string `json:"dcName"`
}


