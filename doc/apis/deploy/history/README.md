查看发布历史
============

### 查看某个Deployment下的发布历史
--------------------------------------

请求URL: 

```bash
GET /api/v1/organizations/{orgId}/datacenters/{dcId}/deployments/{name}/history
```

请求头中包含: Authorization: ${sessionId}


返回值——该Deployment下得RS列表:

每个RS中包括:

* Revision: 修订版本

* Name: 名称

* Namespace: 命名空间

* Image: 镜像名称

* Selector: 选择器

* Replicas: 当前/目标副本数


返回值格式:

```json
{
  
    "code": 0,
    "message": "....",
    "data": [
      {
        "name": "xxx",
        "namespace": "xxx",
        "Image": "xxxx:xx",
        "Selector": "xxxxx":"xxx", "xxx":"xxx",
        "Replicas": {
          "Current": xxx,
          "Desired": xxx
        }
      },
      {
        "name": "xxx",
        "namespace": "xxx",
        "Image": "xxxx:xx",
        "Selector": "xxxxx":"xxx", "xxx":"xxx",
        "Replicas": {
          "Current": xxx,
          "Desired": xxx
        }
      }
    ]
```
