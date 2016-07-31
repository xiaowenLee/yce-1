package dcquota

import (
	localtime "app/backend/common/util/time"
)

const (
	DCQUOTA_SELECT = "SELECT id, dcId, orgId, podNumList, podCpuMax, podMemMax, podCpuMin, podMemMin, rbdQuota, podRbdMax, podRbdMin, price, status, createdAt, modifiedAt, modifiedOp, comment FROM dcquota WHERE id=?"
	DCQUOTA_INSERT = "INSERT INTO dcquota(dcId, orgId, podNumList, podCpuMax, podMemMax, podCpuMin, podMemMin, rbdQuota, podRbdMax, podRbdMin, price, status, createdAt, modifiedAt, modifiedOp, comment) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	DCQUOTA_UPDATE = "UPDATE dcquota SET"

	VALID = 1
	INVALID = 0
)

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
	Status 		int32  `json:"status"`
	CreatedAt   string `json:"createdAt`
	ModifiedAt  string `json:"modifiedAt"`
	ModifedOp   int32  `json:"modifiedOp"`
	Comment     string `json:"comment"`
}

func NewDcQuota(id, dcId, orgId, podNumLimit, podCpuMax, podMemMax, podCpuMin, podMemMin, rbdQuota, PodRbdMax, podRbdMin, modifiedOp int32, price, comment string) *DcQuota {
	return &DcQuota{
		Id:          id,
		DcId:        dcId,
		OrgId:       orgId,
		PodNumLimit: podNumLimit,
		PodCpuMax:   podCpuMax,
		PodMemMax:   podMemMax,
		PodCpuMin:   podCpuMin,
		PodMemMin:   podMemMin,
		RbdQuota:    rbdQuota,
		PodRbdMax:   PodRbdMax,
		PodRbdMin:   podRbdMin,
		Price:       price,
		Status: VALID,
		CreatedAt:   localtime.NewLocalTime().String(),
		ModifiedAt:  localtime.NewLocalTime().String(),
		ModifedOp:   modifiedOp,
		Comment:     comment,
	}
}


