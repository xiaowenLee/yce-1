<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

更新数据中心
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-10-27

目录
--------------
###目的
由管理员更新数据中心信息. 管理员首次登录可能需要对默认的数据中心名称进行更新。


###请求

* 请求方法: POST 
* 请求URL: /api/v1/datacenter/update
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON
```json
{
    "op": 1,              // 管理员userId
    "name": "xxx",
    "orgId": "3",         // 管理员所属orgId
    "nodePort": [
        30000,
        32767
    ] 
    "host": "xxx",
    "port": 8080,
    "secret": "xxx"       // 暂时为空
}
```


###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 更新数据中心 
YCE-->>MySQL: 更新数据中心记录
YCE<<--MySQL: 返回更新结果
```

说明: 收到POST请求, 将用户填写的数据中心信息写入MySQL里。

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注
无
