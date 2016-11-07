<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

检查服务重名
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
查询服务是否重名(访问点和服务是同样的名字)

###请求

* 请求方法: POST
* 请求URL: /api/v1/organizations/:orgId/users/:userId/services/check
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "name": "xxx",
}
```

###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>K8s: 在K8s查询是否重名 
YCE<<--K8s: 返回查询结果
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
查重名除了要求数据库里不重复,还要求k8s里不重复

访问点检查重名
请求URL: POST /api/v1/organizations/:orgId/users/:userId/endpoints/check
请求头: Authorization:SessionId

携带数据:
```
{
    "name": "xxx",
}
```


### 以下为旧版本, 无效///////////////////////////////////////////////////
服务及访问点重名查询:
请求URL: POST /api/v1/organizations/:orgId/users/:userId/services/check
请求头: Authorization:SessionId

携带数据:
```
{
    "name": "xxx",
}
```

返回数据: code为1415, 表示该名称已被占用,不可使用; code为0表示该名称未被占用, 可以使用。

思路:

去每个数据中心获取所有的service的名称, 分别比较, 一旦发现重名即返回1415

请求URL: POST /api/v1/organizations/:orgId/users/:userId/endpoints/check
请求头: Authorization:SessionId

携带数据:
```
{
    "name": "xxx",
}
```

返回数据: code为1415, 表示该名称已被占用,不可使用; code为0表示该名称未被占用, 可以使用。

思路:

去每个数据中心获取所有的endpoints的名称, 分别比较, 一旦发现重名即返回1415
