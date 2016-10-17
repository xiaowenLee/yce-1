### 创建数据中心
-----------

#### 创建初始化 
目的: 为创建数据中心做准备, 获取组织列表供管理员为数据中心选择
请求URL: /api/v1/datacenter/init
请求头: Authorization:SessionId
请求方法: GET 

内容需讨论

#### 数据中心检查
目的: 当管理员输入数据中心名完毕后(离开输入框), 检查数据中心名是否重复
请求URL: /api/v1/datacenter/check
请求头: Authorization:SessionId
请求方法: POST

携带数据:
```
{
    "name": "xxx",
    "orgId": "1"          //表示创建者(管理员)所在的组织,用来验证管理员会话, 从本地存储中获取
}
```

返回是否存在, "code": 1415 为数据中心名已存在, 不能使用该名称, 需提示。 "code": 0为未被占用, 可以使用该名称, 无需提示。

程序实现逻辑:

去datacenter表里选择满足name的数据中心,如果有,返回存在,如果没有,返回不存在

还应该检查Datacenter Host:Port对是否重复. 重复了会导致查询上的问题,但不影响实际部署 

#### 创建
请求URL: /api/v1/datacenter/new
请求头: Authorization:SessionId
请求方法: POST

携带数据:
```
{
    "name": "xxx",
    "nodePort" [
        "30000",
        "32767"
    ]
    "host": "192.168.1.110",
    "port": 8080,
    "orgId": "3",          // 表示创建者所在的组织, 用来验证管理员会话 
    "op": 1           // 管理员datacenterId
    //"secret": xxx       // 暂时空接口
}
```

#### 数据中心列表
请求URL: /api/v1/datacenter
请求头: Authorization:SessionId
请求方法: GET 

返回数据:
```
{
    "code": 0,
    "msg": "...",
    "data": {
        "datacenters": [{
            "id": 1,
            "name": "abc.de",
            ....
        }] 
    }
}
```

列表显示内容:
ID, 数据中心名, 地址, 端口, NodePort, 创建时间, 操作

数据均从data.datacenters的数组里每个元素中取, 其中: 

* ID: id
* 数据中心名: name
* 地址: host
* 端口: port
* NodePort: nodePort[0] - nodePort[1] //暂定
* 创建时间: createdAt
* 操作: 更新、删除


#### 删除数据中心
请求URL: /api/v1/datacenter/delete
请求方法: POST
请求头: Authorization:SessionId
携带数据:
```
{
    "op": 1,
    "orgId": "3",
    "name": "xxx"
}
```

#### 更新数据中心
请求URL: /api/v1/datacenter/update
请求方法: POST
请求头: Authorization:SessionId

携带数据:
```
{
    "op": 1,              // 管理员userId
    "name": "xxx",
    "orgId": "3",         // 管理员所属orgId
    "nodePort": [
        30000,
        32767
    ] 
    "host": "xxx",
    "port": 8080,
    "secret": "xxx"       // 暂时为空
}
```
