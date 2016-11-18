package quota

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log


type Quota struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Cpu        int32  `json:"cpu"`
	Mem        int32  `json:"mem"`
	Rbd        int32  `json:"rbd"`
	Price      string `json:"price"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

