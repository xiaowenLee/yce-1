<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

检查应用重名
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
###目的
用户创建应用时检查应用名是否重复

###请求

* 请求方法: POST
* 请求URL: /api/v1/organizations/{:orgId}/users/:userId/deployments/check
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
```Title: 创建组织
YCE-->>K8s: 在每个数据中心的k8s集群里查询是否有该应用名重复 
YCE<<--K8s: 返回查询结果 
```

###响应数据结构: 
返回码为0, 表示可用。
返回码为1415表示用户名已存在

### 备注