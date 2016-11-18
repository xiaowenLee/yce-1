package dcquota

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log



type DcQuota struct {
	Id          int32  `json:"id"`
	DcId        int32  `json:"dcId"`
	OrgId       int32  `json:"orgId"`
	PodNumLimit int32  `json:"podNumLimit"`
	PodCpuMax   int32  `json:"podCpuMax"`
	PodMemMax   int32  `json:"podMemMax"`
	PodCpuMin   int32  `json:"podCpuMin"`
	PodMemMin   int32  `json:"podMemMin"`
	RbdQuota    int32  `json:"rbdQuota"`
	PodRbdMax   int32  `json:"podRbdMax"`
	PodRbdMin   int32  `json:"podRbdMin"`
	Price       string `json:"price"`
	Status      int32  `json:"status"`
	CreatedAt   string `json:"createdAt`
	ModifiedAt  string `json:"modifiedAt"`
	ModifiedOp  int32  `json:"modifiedOp"`
	Comment     string `json:"comment"`
}

