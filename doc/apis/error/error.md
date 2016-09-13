YceError
=========

为了简化日志打印同时返回前端结果的操作, 特定义YceError。

基本定义及关键方法如下:

```golang

type YceError struct {
    Code int32 `json:"code"`
    Message string `json:"message"`
    Data string `json:"data,omitempty"`
}


func (e *YceError) EncodeSelf() []byte {
    errJSON, _ := json.Marshal(e)
    return errJSON
}

```

同时打印日志和返回错误:

```golang

func (a *abcController) Response() {
    log.Errorf("%s\n", error.Errors[error.EMYSQL])
    a.Write(string(error.NewYceError(error.EMYSQL, error.Errors[error.EMYSQL], "").EncodeSelf()))
}

```

### 错误类别及处理规则

1. 普通错误

普通错误包括输入值错误, 输入合法性校验在前端完成。例如应用名的校验, 除了要满足系统的要求外, 还需跟k8s的要求一致。

创建k8s资源失败, 要返回失败的错误信息。

对数据库的读操作: 如查询失败, 这些无需重启, 只需返回失败的错误信息即可。

对数据库的写操作: 如插入、删除、更新失败, 如果对应着k8s资源的操作, 例如: 创建K8s资源成功, 写入数据库失败的情况, 怎么处理? 数据库与k8s对象操作强一致

如果独立于k8s的操作, 例如更改用户组织、数据中心、预算及配额等, 返回错误信息。 

对于redis的读写操作同上。如果读写失败, 均表示用户登录出错, 需要用户多次尝试重新登录。

2. 内存错误

致命错误值引发panic的错误, 例如数组越界、空指针错误、死锁等, 这类错误应该被捕捉到,并返回出错信息。为了减少该类错误, 需要减少使用定长数组,减少使用下标,转为使用切片。尽量少返回或传递nil值,如有需要,返回特殊值。比如返回带有特殊值的结构体来代替返回nil

临界资源加锁要谨慎。

减少此类错误导致的崩溃。一旦发生, 立即重启。


3. 系统错误

mysql连接错误、redis连接错误、k8s连接错误, iris错误等。

redis连接超时, 关系到用户登录会话, 所以需要重启系统。

mysql连接超时提示操作超时, 返回超时信息, 并重启。

iris运作错误系统必须重启

k8s连接错误可能需要分数据中心进行返回, 而不能因为其中一个k8s连接失败而导致其他k8s的值也无法显示。






