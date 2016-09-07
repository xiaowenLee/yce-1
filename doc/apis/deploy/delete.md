删除应用
------------


用户点击删除应用时提示拼接json, 点击确认删除时发送请求

请求的方法及URL: DELETE /api/v1/organizations/{ordId}/deployment/{deploymentName}

请求头中包含: Authorization: ${sessionId}

发送Json格式:

```json
  {
    "userId": "1", 
    "dcId": 1, 
    "comments": "delete busybox-test-delete"
  }
    
```

返回值:

* 操作结果 


### 删除步骤

先获取该数据中心里该命名空间下的该名称应用

dcId --> orgName --> deploymentName

然后删除两步:

先获取该deployment对应的所有replicaSet, 并依次删除

再将该deployment直接删除。

不等待,不处理删除失败恢复,仅返回操作结果
