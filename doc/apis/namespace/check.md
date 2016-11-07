<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

检查组织重名
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
管理员输入组织名称后检查是否重复

###请求

* 请求方法: POST
* 请求URL: /api/v1/organization/check
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "orgName": "xxx",   // 待创建的组织名称
    "orgId": "xxx"      // 管理员的orgId, 用来认证会话
}
```

###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>MySQL: 在后台查询是否重名 
YCE<<--MySQL: 返回查询结果
```


###响应数据结构: 
"code": 1415 为用户名已存在, 不能使用该名称, 需提示。 
"code": 0为未被占用, 可以使用该名称, 无需提示。
JSON
```json
{
    "code": 0,
    "message": "....",
    "data": 
}
```


### 备注
无
