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

