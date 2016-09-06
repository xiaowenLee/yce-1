### 访问点删除

用户在服务管理页点击删除时向后台发送数据:

请求的URL:

DELETE /api/v1/organizations/{orgId}/datacenters/{dcId}/users/{userId}/endpoints/{epName}

请求头中包含: Authorization: ${sessionId} 

其中: userId, orgId, sessionId在登录成功后从后台返回给浏览器, 前端存储在LocalStorage(目前为SessionStorage)里面


返回值:

* 操作结果

