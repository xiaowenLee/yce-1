package namespace

import (
	mylog "app/backend/common/util/log"
)

var log = mylog.Log

type AccountType struct {
	Budget float32 `json:"budget"`
	Balance float32 `json:"balance"`
}

type QuotaPkgType struct {
	Name string `json:"name"`
	Cpu int32 `json:"cpu"`
	Mem int32 `json:"mem"`
	Cost float32 `json:"cost"`
}
