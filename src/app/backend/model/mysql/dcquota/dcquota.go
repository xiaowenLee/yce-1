package dcquota

import (
	"big"
	"cmd/compile/internal/big"
)

type DcQuota struct {
	Id          int32   `json:"id"`
	Price       big.Rat `json:"price"`
	DcId        int32   `json:"dcId"`
	OrgId       int32   `json:"orgId"`
	PodNumLimit int32   `json:"podNumLimit"`
	PodCpuMax   int32   `json:"podCpuMax"`
	PodMemMax   int32   `json:"podMemMax"`
	PodCpuMin   int32   `json:"podCpuMin"`
	PodMemMin   int32   `json:"podMemMin"`
	RbdQuota    int32   `json:"rbdQuota"`
	podRbdMax   int32   `json:"podRbdMax"`
	podRbdMin   int32   `json:"podRbdMin"`
	CreatedAt   string  `json:"createdAt`
	ModifiedAt  string  `json:"modifiedAt"`
	ModifedOp   int     `json:"modifiedOp"`
	Comment     string  `json:"comment"`
}
