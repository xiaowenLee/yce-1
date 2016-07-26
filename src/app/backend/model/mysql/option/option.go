package option

import (
	"time"
)

type OptionType int32 //操作类型

const (
	GET            = iota // 查询
	ONLINE                // 上线
	ROLLINGBACK           // 回滚
	ROLLINGUPGRADE        // 滚动升级
	CANCEL                // 取消上线,下线
	PAUSE                 // 暂停上线
	RESUME                // 恢复上线

)

type Option struct {
	Id             int32     `json:"id"`
	Name           string    `json:"name"`
	CreateTs       string    `json:"create_ts"`
	LastModifiedTs time.Time `json:"last_modified_ts"`
	LastModifiedOp int       `json:"last_modified_op"`
	Comment        string    `json:"comment"`
}
