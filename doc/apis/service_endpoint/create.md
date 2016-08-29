### 服务发布准备

用户点击服务发布(左侧菜单)时请求后台数据:

请求的URL:

GET /api/v1/organizations/{orgId}/users/{userId}/services/init

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
                "orgName": "Ops",
                "nodePort": 30000,  
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

这些将关系到用户部署服务到哪个机房

重名校验在哪里做?前端OR后端?目前发布的重名校验呢?要统一

用户需要确定:

* 服务名:
* 服务所属数据中心:
* selector:
* 服务部署方式:

    1. withoutselector, 部署服务的时候自动创建Endpoint

    2. withselector, 部署服务的时候还需指定或创建Endpoint

* 服务类型:

    1. ClusterIP

    2. NodePort

    3. LoadBalancer
    
    4. External Name


仅允许ClusterIP和NodePort两种类型



