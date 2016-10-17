创建K8s的命名空间
===============

此项功能是用于后台管理员操作的功能,跟普通用户无关
----------------------------


### 为创建组织或更新组织做准备
为创建组织或更新组织准备: 请求数据中心列表、账户余额信息、资源套餐信息
请求URL: GET /api/v1/organization/init     // URL里面暂为organization
请求头: Authorization: ${SessionId}

返回数据:
```
{
    "code": 0,
    "msg": "...",
    "data": {
        "dcList": [{
            "id": 1,
            "name": "xx"
            ...
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
在创建命名空间的页面,在填写完组织名称后,发起后台校验,成功或失败。 组织名具有全局唯一性
请求URL: POST /api/v1/organization/check
请求头中包含: Authorization: ${sessionId}
请求体重提交: 

```
{
    "orgName": "xxx",   // 待创建的组织名称
    "orgId": "xxx"      // 管理员的orgId, 用来认证会话
}
```

返回值的格式:
返回在该组织里是否存在, "code": 1415 为用户名已存在, 不能使用该名称, 需提示。 "code": 0为未被占用, 可以使用该名称, 无需提示。
例如可用:
```json
{
    "code": 0,
    "message": "....",
    "data": 
}
```


### 创建命名空间
请求URL: POST /api/v1/organization/new
请求头中包含: Authorization: ${sessionId}
请求体为Json格式:
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

返回值:
* 组织Id
* 组织名称
* 该组织下每个数据中心的k8s:
    * 命名空间名称
    * 资源使用配额: CPU和内存
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
3. LimitRange和数据中心配额表下个版本再添加 //目前没实现

### 更新组织信息,主要是购买资源

资源及账户信息在list的时候返回
请求URL: POST /api/v1/organization/update
请求头: Authorization: SessionId

携带数据:
```
{
    "userId": xxx,        // 操作者的用户ID,这里是数字,不是字符串, 管理员
    "name":  "xxx",       // 组织名称
    "orgId": "xxx",       // 操作者的组织ID(管理员),不是新创建的那个组织ID
    "cpuQuota": 100,      // CPU配额
    "memQuota": 200,      // 内存配额
  //目前只返回上面内容的json, 即只更新配额 
    "dcIdList": [1, 2]    // 数据中心列表
    "budget": 10000,      // 预算
    "balance": 10000,     // 余额
}
```

资源配额(套餐)已经在组织管理返回了quotaPkg
    


### 获取组织列表

请求URL: GET /api/v1/organization
请求头: Authorization: ${sessionId}
请求返回:
```
{
    "code": 0,
    "msg": "...",
    "data": {
        "organizations": [{
            "id": 1,
            "name": "xxx",
            ... 
        }],
        "dcList": [{
            "dcId": dcName,    //  map例如: "1": "电信" 
        }]
        
    }
}
```

显示内容:
ID, 组织名称, CPU, MEM, 账户余额, 创建时间, 数据中心, 操作

数据均从data.organizations的数组里每个元素中取, 其中: 

* ID: id
* 组织名称: name
* CPU:  cpuQuota
* MEM:  memQuota
* 账户余额: balance
* 创建时间: createdAt
* 数据中心: dcIdList
* 操作: 更新、删除

dcIdList为JSON, 从里面取值到dcList里获取相应的名字对应起来

### 删除组织
请求URL: POST /api/v1/organization/delete 待讨论
请求头: Authorization: ${sessionId}

携带数据: 待讨论
```
{
    "orgName": "xxx",
    "orgId": "3",
    "op": 1
}

```

