创建K8s的命名空间
===============

此项功能是用于后台管理员操作的功能,跟普通用户无关
----------------------------


### 为创建组织或更新组织做准备
为创建组织或更新组织准备: 请求数据中心列表、账户余额信息、资源套餐信息
请求URL: GET /api/v1/organizations/init
请求头: Authorization: ${SessionId}

返回数据:
```
{
    "code": 0,
    "msg": "...",
    "data": {
        "dcList": [{
            "dcId": 1,
            "dcName": "xx"
        }],
        "account": {        // 首次数据库里查找不到, 返回0, 要求充值 或者 首次注册赠送500元
            "budget": 100,
            "balance": 50
        },
        "quotaPkg": [{
            "cpu": 100,
            "mem": 200,
            "cost": 200
        }]
    }
}
```

### 校验新的名字是否可用(是否已经存在)

在创建命名空间的页面,在填写完组织名称后,发起后台校验,成功或失败,都在页面上给于提示. 逻辑是只要数据库里有就不建议取这个名字。

请求URL: POST /api/v1/organizations/check

请求头中包含: Authorization: ${sessionId}

请求体重提交: 
```
{
    "orgName": "xxx",   // 待创建的组织名称
    "orgId": "xxx"      // 管理员的orgId, 用来认证会话
}
```

返回值:

* 是否可用(数据库org表)

返回值的格式:

成功时:
```json
{
    "code": 0,
    "message": "....",
    "data": ","
}
```

失败时:
```json
{
    "code": 1414,
    "message": "该名称已被占用"
    "data": {
      "dcIdList": [
          1 
      ],
      "dcNameList": [
          "xxx" //已被这些数据中心所使用, 图片变灰不可用, 要求修改名字 
      ]
    }
}
```

程序逻辑实现:

在organization表中判断是否已经存在该组织名,如果存在,返回错误;如果不存在,返回成功 



### 创建命名空间

请求URL: POST /api/v1/organizations/new

请求头中包含: Authorization: ${sessionId}

请求体为Json格式:
```json
{
    "orgId": "xxx",       // 操作者的组织ID(管理员),不是新创建的那个组织ID
    "userId": xxx,        // 操作者的用户ID,这里是数字,不是字符串, 管理员
    "orgName": "dev",        // 组织名称也是K8s里namespace的名称
    "dcIdList": [1, 2]      // 数据中心列表
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


### 更新组织信息,主要是购买资源

资源及账户信息在list的时候返回
请求URL: POST /api/v1/organizations/update
请求头: Authorization: SessionId

携带数据:
```
{
    "orgId": "xxx",       // 操作者的组织ID(管理员),不是新创建的那个组织ID
    "userId": xxx,        // 操作者的用户ID,这里是数字,不是字符串, 管理员
    "orgName": xxx,       // 组织名称
    "dcIdList": [1, 2]    // 数据中心列表
    "orgName": "dev",     // 组织名称也是K8s里namespace的名称
    "cpuQuota": 100,      // CPU配额
    "memQuota": 200,      // 内存配额
    "budget": 10000,      // 预算
    "balance": 10000,     // 余额
}
```

    


### 获取组织列表

请求URL: GET /api/v1/organizations
请求头: Authorization: ${sessionId}
请求返回:
```
{
    "code": 0,
    "msg": "...",
    "data": [{
        "orgId": 1,
        "orgName": "xxx",
        "dcList": [{
            "dcId": 1,
            "dcName": "xxx" 
        }] 
        "cpuQuota": xxx,
        "memQuota": xxx,
        "budget": xxx,
        "balance": xxx
    }]
}
```
