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
    "data": [{
            "dcName": ["bangongwang", "dianxin"]
            "userName": "admin",
            "record": deployRecord mysql.Deployment
    }]
}
```
