package quota

import (
	"time"
)

type Quota struct {
	Id         int32   `json:"id"`
	Name       string  `json:"name"`
	Cpu        int32   `json:"cpu"`
	Mem        int32   `json:"mem"`
	Rbd        int32   `json:"rbd"`
	Price      float32 `json:"price"`
	CreatedAt  string  `json:"createdAt"`
	ModifiedAt string  `json:"modifiedAt"`
	ModifiedOp int     `json:"modifiedOp"`
	Comment    string  `json:"comment"`
}
