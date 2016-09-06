应用列表
===========

用户点击手动扩容时时请求后台数据:

请求的方法及URL: POST /api/v1/organizations/{orgId}/deployments/{deploymentName}/scale

请求头中包含: Authorization: ${sessionId} 

传递参数JSON:

```json
    {
      "newSize": 5,
      "dcIdList": [
        1
       ],
      "userId": 1,
      "comments": "scale to 5 instances"
    }

```

返回值:

操作结果

