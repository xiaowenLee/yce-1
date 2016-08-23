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
    New("err code not found")
}


type yerror interface{
    Error() string
}

type YceError struct {
    errmsg string 
}

func New(text string) *YceError {
    return &YceError{errmsg: text}
}

func (e *YceError) Error() string {
    return e.errmsg
}

func (e *YceError) EncodeSelf() []byte {
    errJSON, _ := json.Marshal(e.errmsg)
    return errJSON
}


func (x xxxController) LogAndResponse() {
    x.YceLogger.Logger.Printf()
    x.Write(string(e.EncodeSelf()))
}


```