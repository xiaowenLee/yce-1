<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

应用发布历史
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
为回滚提供方便, 返回应用发布的历史

###请求

* 请求方法: GET 
* 请求URL: /api/v1/organizations/{orgId}/operationlog
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 

  无

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 发布应用
YCE-->>MySQL: 根据应用名查询所有发布记录
YCE<<--MySQL: 返回查询结果
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

JSON
```json
 {
    "code":0,
    "message":[
        "OK"
    ],
    "data": {
      "operationLog":  [{
            "dcName": ["bangongwang", "dianxin"]
            "userName": "admin",
            "record": deployRecord mysql.Deployment
      }]
    }
} 
```
### 备注
无



### 以下为旧版本, 无效///////////////////////////////////////////////////


应用发布历史
===========

用户点击应用管理时请求后台数据:

请求的方法及URL: GET /api/v1/organizations/{orgId}/operationlog

请求头中包含: Authorization: ${sessionId} *暂时在Session Storage里*

返回值:

* 该组织下该用户的应用发布历史


```json
{
    "code":0,
    "message":[
        "OK"
    ],
    "data": {
      "operationLog":  [{
            "dcName": ["bangongwang", "dianxin"]
            "userName": "admin",
            "record": deployRecord mysql.Deployment
      }]
    }
}
```
