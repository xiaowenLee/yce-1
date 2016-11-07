<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

删除应用
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
由用户删除所发布的应用

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{ordId}/deployment/{deploymentName}
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON
```json
 {
    "userId": "1", 
    "dcIdList": [1], 
    "comments": "delete busybox-test-delete"
  } 
```

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 删除应用
YCE-->>MySQL: 根据orgId找到对应dcId里的orgName,
YCE<<--MySQL: 返回请求结果
YCE-->>K8s: 根据orgName和deploymentName获取应用的RS及Pod,并删除
YCE<<--K8s: 返回删除结果
YCE-->>MySQL: 插入操作记录
YCE<<--MySQL: 返回插入结果
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

### 备注
应该将删除Deployment和数据库插入做成事务,保持一致。但是不同于创建,应用删除后难以恢复。
无状态的能恢复,有状态的不易恢复。可能需要保存它的配置文件(yaml等)才能恢复。


### 以下为旧版本, 无效///////////////////////////////////////////////////
删除应用
------------


用户点击删除应用时提示拼接json, 点击确认删除时发送请求

请求的方法及URL: DELETE /api/v1/organizations/{ordId}/deployment/{deploymentName}

请求头中包含: Authorization: ${sessionId}

发送Json格式:

```json
  {
    "userId": "1", 
    "dcIdList": [1], 
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
