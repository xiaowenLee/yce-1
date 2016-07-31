package quota

import (
	localtime "app/backend/common/util/time"
)

const (
	VALID = 1
	INVALID = 0
)

type Quota struct {
	Id         int32   `json:"id"`
	Name       string  `json:"name"`
	Cpu        int32   `json:"cpu"`
	Mem        int32   `json:"mem"`
	Rbd        int32   `json:"rbd"`
	Price      string  `json:"price"`
	Status     int32   `json:"status"`
	CreatedAt  string  `json:"createdAt"`
	ModifiedAt string  `json:"modifiedAt"`
	ModifiedOp int     `json:"modifiedOp"`
	Comment    string  `json:"comment"`
}

func NewQuota(name, price, comment string, cpu, mem, rbd, modifiedOp int32) {
	return &Quota{
		Name: name,
		Cpu: cpu,
		Mem: mem,
		Rbd: rbd,
		Price: price,
		Status: VALID,
		CreatedAt: localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment: comment,
	}
}

func QueryQuotaById() {

}

func InsertQuota() {

}

func DeleteQuota() {

}

func UpdateQuota() {

}
