<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

应用扩容
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
对该应用进行扩(缩)容

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{orgId}/deployments/{deploymentName}/scale
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON
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

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 发布应用
YCE-->>K8s: 发出扩容请求 
YCE<<--K8s: 返回扩容结果
YCE-->>MySQL: 插入操作记录
YCE<<--MySQL: 返回插入结果
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

### 备注
无

### 以下为旧版本, 无效///////////////////////////////////////////////////
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

