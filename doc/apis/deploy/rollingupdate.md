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
  "dcId": 1,
  "userId": 1,
  "strategy": {
      "maxUnavailable": 3,
      "image": "nginx:1.9",
      "updateInterval": 2,
  }
  "comments": "Update the xxx function"
}
```

升级记录需要写入数据库和Annotations

### 镜像搜索辅助(联想)
==================

略不同于应用发布里的镜像搜索联想,这里的联想仅允许确定的镜像,在不同的版本间进行滚动升级。例如,确定的nginx里,提供1.7.1 ~ 1.9.7的版本供升级。
