<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

应用滚动升级
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
对该应用进行滚动升级

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{orgId}/deployments/{:deploymentName}/rolling
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON
```json
 {
  "appName": "nginx-test",
  "dcIdList": [1],
  "orgName":"ops",
  "userId": 1,
  "strategy": {
      "image": "nginx:1.9"
  },
  "comments": "Update the xxx function"
}
```

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 发布应用
YCE-->>K8s: 发出滚动升级请求 
YCE<<--K8s: 返回升级结果
YCE-->>MySQL: 插入操作记录
YCE<<--MySQL: 返回插入结果
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

### 备注
滚动升级策略两个参数: MaxUnavalable/MaxSurge参数在后台中采用默认值,均为2


### 以下为旧版本, 无效///////////////////////////////////////////////////

应用滚动升级
============


### 应用滚动升级

用户点击应用管理里,某个应用的"滚动升级"按钮时向后台发送数据:

请求的URL:

POST /api/v1/organizations/{orgId}/deployments/{:deploymentName}/rolling

请求头中包含: Authorization: ${sessionId}

其中: userId, orgId, sessionId在登录成功后从后台返回给浏览器, 前端存储在LocalStorage(目前为SessionStorage)里面

上传值:

* 升级策略(默认为滚动升级) 

* 最大不可用实例数(不允许为全部)

* 升级版本

* 升级间隔(预留) 

* 升级说明

提交的大概数据格式如下:

```json
{
  "appName": "nginx-test",
  "dcIdList": [1],
  "orgName":"ops",
  "userId": 1,
  "strategy": {
      "image": "nginx:1.9"
  },
  "comments": "Update the xxx function"
}
```

注意: 
    滚动升级策略两个参数: MaxUnavalable/MaxSurge参数在后台中采用默认值,均为2


升级记录需要写入数据库和Annotations

### 镜像搜索辅助(联想)
==================

略不同于应用发布里的镜像搜索联想,这里的联想仅允许确定的镜像,在不同的版本间进行滚动升级。例如,确定的nginx里,提供1.7.1 ~ 1.9.7的版本供升级。
