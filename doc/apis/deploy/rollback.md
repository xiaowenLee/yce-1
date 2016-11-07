<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

应用回滚
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
对该应用进行回滚

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/:orgId/deployments/:deploymentName/rollback 
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON
```json
{
    "appName": "xxx",
	"dcIdList": [1],
	"userId": "xxx",
	"image": "xxx",
	"revision": "xxx",
	"comments": "xxx"
}
```

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 发布应用
YCE-->>K8s: 发出回滚请求 
YCE<<--K8s: 返回回滚结果
YCE-->>MySQL: 插入操作记录
YCE<<--MySQL: 返回插入结果
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

### 备注
无