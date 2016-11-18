<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

删除服务
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-17

目录
--------------
###目的
由用户删除服务


###请求

* 请求方法: POST
* 请求URL: /api/v1/organizations/{orgId}/datacenters/{dcId}/users/{userId}/services/{svcName}
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 



###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>K8s: 删除对应的服务
YCE<<--K8s: 返回删除结果
```


###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。


### 备注
暂时方案: 目前在删除服务后还会去删除相关的访问点,如果能找到同名的访问点,就删除。如果找不到,则只删除服务。

### 以下为旧版本, 无效///////////////////////////////////////////////////


### 服务删除

用户在服务管理页点击删除时向后台发送数据:

请求的URL:

//GET /api/v1/organizations/{orgId}/datacenters/{dcId}/services/{serviceName}
//DELETE /api/v1/organizations/{orgId}/datacenters/{dcId}/users/{userId}/services/{svcName}
DELETE /api/v1/organizations/{orgId}/datacenters/{dcId}/users/{userId}/services/{svcName}

请求头中包含: Authorization: ${sessionId}, NodePort: ${nodePort}

其中: userId, orgId, sessionId在登录成功后从后台返回给浏览器, 前端存储在LocalStorage(目前为SessionStorage)里面


返回值:

* 操作结果

