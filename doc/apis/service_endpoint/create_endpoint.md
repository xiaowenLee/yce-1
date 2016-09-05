### 访问点发布准备

用户点击访问点发布(左侧菜单)时请求后台数据:

请求的URL:

GET /api/v1/organizations/{orgId}/users/{userId}/endpoints/init

请求头中包含: Authorization: ${sessionId}

其中: userId, orgId, sessionId在登录成功后从后台返回给浏览器, 前端存储在LocalStorage(目前为SessionStorage)里面

返回值:

* 组织名称

* 该组织下的数据中心列表



获取成功的大概数据格式如下:

```json

   {
      "code": 0,
      "message": "请求成功",
      "data": {
                "orgId":  "1",
                "orgName": "ops",
                "dataCenters": [
                {
                    "dcId": "1",
                    "name": "世纪互联",
                    "budget": 10000000,
                    "balance": 10000000
                },
                {
                    "dcId": "2",
                    "name": "电信机房",
                    "budget": 10000000,
                    "balance": 10000000
                },
                {
                    "dcId": "3",
                    "name": "电子城机房",
                    "budget": 10000000,
                    "balance": 10000000
                }
                ],
      }
   } 
    
```

这些将关系到用户部署访问点到哪个机房


用户需要确定:

* 访问点名:
* 访问点所属数据中心: 依据GET请求的返回值确定。卡片式
* 所属组织
* Labels
* 地址端口组:
  多条:
 
  IP:PORT
  
  
### 访问点发布提交
用户点击提交请求后台数据:
  
请求的URL:
  
POST /api/v1/organizations/{orgId}/users/{userId}/endpoints/new
  
请求头中包含: Authorization: ${sessionId}