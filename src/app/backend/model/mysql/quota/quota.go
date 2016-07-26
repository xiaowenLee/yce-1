package quota

import (
	"time"
)

type Quota struct {
	Id             int32     `json:"id"`
	Name           string    `json:"name"`
	Cpu            int32     `json:"cpu"`
	Mem            int32     `json:"mem"`
	Rbd            int32     `json:"rbd"`
	Price          float32   `json:"price"`
	CreatedTs      time.Time `json:"created_ts"`
	LastModifiedTs time.Time `json:"last_modified_ts"`
	LastModifiedOp int       `json:"last_modified_op"`
	Comment        string    `json:"comment"`
}
