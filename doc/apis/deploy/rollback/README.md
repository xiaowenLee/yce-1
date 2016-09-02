应用回滚
=========================

### 根据Deployment的名字回滚应用
-------------------------------------

请求地址:
```bash
POST /api/v1/organizations/{orgId}/deployments/{name}/rollback
```

请求头中包含: Authorization: ${sessionId}

请求体:

* 用户ID

* 镜像名称

* 回滚说明

* 回滚的修订版本号

* 数据中心ID


提交的大概数据格式如下:

```json
{
  "appName": "nginx-test",
  "dcId": 1,
  "userId": 1,
  "image": "xxxx:xxx",
  "revision": "xxxx",
  "comments": "Update the xxx function"
}
```

### 镜像搜索辅助(联想)
==================

### 镜像搜索辅助(联想)
==================

略不同于应用发布里的镜像搜索联想,这里的联想仅允许确定的镜像,在不同的版本间进行。

这个要讨论一下。