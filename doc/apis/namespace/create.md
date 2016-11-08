<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

创建组织
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
由管理员创建组织(K8s里的命名空间)

###请求

* 请求方法: POST
* 请求URL: /api/v1/organization/new
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "orgId": "xxx",       // 操作者的组织ID(管理员),不是新创建的那个组织ID
    "userId": xxx,        // 操作者的用户ID,这里是数字,不是字符串, 管理员
    "name": "dev",        // 组织名称也是K8s里namespace的名称
    "dcIdList": [1],      // 数据中心列表
    "cpuQuota": 0,        // blow 4 items will be modified in update.go, omitempty
    "memQuota": 0,
    "budget":   0,
    "balance":  0
}
```


###页面设计 
无


###程序实现逻辑
```Title: 创建组织
YCE-->>MySQL: 在organization表中插入一条数据  
YCE<<--MySQL: 返回插入结果 
YCE-->>k8s: 每个k8s集群创建namespace和resourceQuota
YCE<<--K8s: 返回创建结果
YCE-->>K8s: 为每个K8s集群limitrange
YCE<<--k8s: 返回创建结果
```

###响应数据结构: 
JSON
```json
{
    "code": 0,
    "message": "...",
    "data": {
        "orgId": "xxx",
        "orgName": "xxx",
        "budget": 1111,
        "balance": 1111,
        "dataCenters": [
            {
                "dcId": 1,
                "namespace": "xxx",
                "resourceQuota": {
                    "cpu": xxx,
                    "mem": xxx
                }

            },
            {
                "dcId": 3,
                "namespace": "xxx",
                "resourceQuota": {
                    "cpu": xxx,
                    "mem": xxx
                }
            }
            ...
        ]

    }
}
```


### 备注
在创建组织的时候应该添加limitrange
```Title: 创建组织
YCE-->>K8s: 为每个K8s集群limitrange
YCE<<--k8s: 返回创建结果
```

