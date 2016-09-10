显示当前命名空间下的集群拓扑图
===============

### 请求topology数据
---------------------------------

请求地址:
```bash
GET organizations/{orgId}/topology
```

请求头中包含: Authorization: ${sessionId}


返回值

* 当前命名空间下所有的: Pod, Node, ReplicaSet, Service

* Node, ReplicaSet, Service与Pod的对应关系

数据格式如下:

```json
{
    "items": {
        "16ad29ed-682e-11e6-b957-44a84240716a": {
            "kind": "Service",
            "apiVersion": "v1beta1",
            "metadata": {
                ...
            }
            ...
        },
        "16ad29ed-682e-11e6-b957-44a84240716a": {
            "kind": "Pod",
            "apiVersion": "v1beta1",
            "metadata": {
                ...
            }
            ...
        },
        "16ad29ed-682e-11e6-b957-44a84240716a": {
            "kind": "Node",
            "apiVersion": "v1beta1",
            "metadata": {
                ...
            }
            ...
        },
        "16ad29ed-682e-11e6-b957-44a84240716a": {
            "kind": "Pod",
            "apiVersion": "v1beta1",
            "metadata": {
                ...
            }
            ...
        }
    },
    "relations": [
      {
        "source": "a46216b5-75c5-11e6-b957-44a84240716a",
        "target": "a4659bf7-75c5-11e6-b957-44a84240716a"
      },
      {
        "source": "a9019e1d-323c-11e6-b9d6-44a84240716a",
        "target": "a4659bf7-75c5-11e6-b957-44a84240716a"
      }
    ]
}
```

