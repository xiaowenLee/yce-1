<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

删除访问点
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
由用户删除访问点


###请求

* 请求方法: POST
* 请求URL: /api/v1/organizations/{orgId}/datacenters/{dcId}/users/{userId}/endpoints/{epName} 
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 



###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>K8s: 删除对应的命名空间
YCE<<--K8s: 返回删除结果
```


###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。


### 备注
无 未实现

### 以下为旧版本, 无效///////////////////////////////////////////////////
### 访问点删除

用户在服务管理页点击删除时向后台发送数据:

请求的URL:

DELETE /api/v1/organizations/{orgId}/datacenters/{dcId}/users/{userId}/endpoints/{epName}

请求头中包含: Authorization: ${sessionId} 

其中: userId, orgId, sessionId在登录成功后从后台返回给浏览器, 前端存储在LocalStorage(目前为SessionStorage)里面


返回值:

* 操作结果

