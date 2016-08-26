创建K8s的命名空间
===============

此项功能是用于后台管理员操作的功能,跟普通用户无关
----------------------------

### 校验新的名字是否可用(是否已经存在)

在创建命名空间的页面,在填写完组织名称后,发起后台校验,成功或失败,都在页面上给于提示

请求URL: POST /api/v1/organizations/init

请求头中包含: Authorization: ${sessionId}

请求体重提交: {"name": "xxx", "orgId": "xxx"}

返回值:

* 是否可用(数据库org表)

返回值的格式:

成功时:
```json
{
    "code": 0,
    "message": "....",
    "code": ""
}
```

失败时:
```json
{
    "code": 1517,
    "message": "该名称已被占用"

}
```

程序逻辑实现:

在organization表中判断是否已经存在该组织名,如果存在,返回错误,如果不存在,返回成功


### 创建命名空间

请求URL: POST /api/v1/organizations/

请求头中包含: Authorization: ${sessionId}

请求体为Json格式:
```json
{
    "name": "dev",        // 组织名称也是K8s里namespace的名称
    "cpuQuota": 100,      // CPU配额
    "memQuota": 200,      // 内存配额
    "budget": 10000,      // 预算
    "balance": 10000,     // 余额
    "dcList": [1, 2]      // 数据中心列表
}
```

返回值:

* 组织Id

* 组织名称

* 该组织下每个数据中心的k8s:
    ** 命名空间名称
    ** 资源使用配额: CPU和内存

* 预算

* 余额

返回值格式:

创建成功时:
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

程序实现逻辑:

1. 先在organization表中新增一条数据,

2. FOREACH每个数据中心,分别创建namespace和resourceQuota

3. LimitRange和数据中心配额表下个版本再添加