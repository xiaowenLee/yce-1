<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

拓扑
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
###目的
为管理员和普通用户提供应用拓扑图显示

###请求

* 请求方法: GET 
* 请求URL: organizations/{orgId}/topology
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 

###页面设计 
无


###程序实现逻辑

* 当前命名空间下所有的: Pod, Node, ReplicaSet, Service

* Node, ReplicaSet, Service与Pod的对应关系


###响应数据结构: 
JSON
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



### 备注
URL有问题, 需勘误

管理员用户看到的跟普通用户的略不同



### 以下为旧版本, 无效///////////////////////////////////////////////////

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

