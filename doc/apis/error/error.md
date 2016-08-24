YceError
=========

为了简化日志打印同时返回前端结果的操作, 特定义YceError。

基本定义如下:

```golang

type Errno uintptr

const (
    ETIMEOUT Errno = 1 /*Operation Time Out*/
)

var errors = [...]string {
    ETIMEOUT: "Operation Time out",
}

func (e Errno) Error() string {
    if 0 <= int(e) && int(e) < len(errors) {
        return errors[e] 
    }
    New("error code not found").Error()
}


type yerror interface{
    Error() string
}

type YceError struct {
    Code int32 `json:"code"`
    Message string `json:"message"`
    Data []byte `json:"data"`
}

func New(text string) *YceError {
    return &YceError{errmsg: text}
}

func (e *YceError) Error() string {
    return e.errmsg
}

func (e *YceError) EncodeSelf() []byte {
    errJSON, _ := json.Marshal(e)
    return errJSON
}

```

使用方法:

在每个controller里定义成员如下:

```golang
type abcController struct {
    *iris.Context
    Logger *YceLogger
    //...
}
```

同时打印日志和返回错误:

```golang
func (a *abcController) LogAndResponse(level LogLevel, code int32, msg string) {
    a.Logger := New(level)
    a.Logger.Printf(msg)
    e := New(code, msg)
    a.Write(string(e.EncodeSelf()))
}

```

如果是定义了的错误码及错误说明,调用LogAndResponse(ERR, ETIMEOUT, ETIMEOUT.Error())
未定义的直接填