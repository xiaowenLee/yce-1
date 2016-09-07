package option

type OptionType int32 //操作类型

const (
	GET            = iota + 1 // 查询
	ONLINE                    // 上线
	ROLLINGBACK               // 回滚
	ROLLINGUPGRADE            // 滚动升级
	SCALING			  // 扩容
	CANCEL                    // 取消上线,下线
	PAUSE                     // 暂停上线
	RESUME                    // 恢复上线
	DELETE			  // 删除

)

type Option struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	CreateAt   string `json:"createAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int    `json:"modifiedOp"`
	Comment    string `json:"comment"`
}
