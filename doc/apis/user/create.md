<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

创建用户
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
###目的
由管理员创建用户

###请求

* 请求方法: GET
* 请求URL: /api/v1/user/new
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "userName": "xxx",
    "password": "xxx",  // 暂时有默认值
    "orgName": "dev",       // 创建用户时选择
    "orgId": "1",          // 表示创建者所在的组织, 用来验证管理员会话 
    "op": "1"           // 管理员userId
}
```


###页面设计 
无


###程序实现逻辑
```Title: 创建用户
YCE-->>MySQL: 在user表中插入一条数据  
YCE<<--MySQL: 返回插入结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注
无





### 创建用户
-----------

#### 创建初始化
目的: 为创建用户做准备, 获取组织列表供管理员为用户选择
请求URL: /api/v1/user/init
请求头: Authorization:SessionId
请求方法: GET 

返回数据:
```
{
    "code": 0,
    "msg": "xxx",
    "data": [
            "dev",
            "ops"
        ] 
}
```


#### 用户名检查
目的: 当管理员输入用户名完毕后(离开输入框), 检查用户名是否重复
请求URL: /api/v1/user/check
请求头: Authorization:SessionId
请求方法: POST

携带数据:
```
{
    "userName": "xxx",
    "orgName": "yyy",   // 
    "orgId": "1"          //表示创建者(管理员)所在的组织,用来验证管理员会话, 从本地存储中获取
}
```

返回在该组织里是否存在, "code": 1415 为用户名已存在, 不能使用该名称, 需提示。 "code": 0为未被占用, 可以使用该名称, 无需提示。

程序实现逻辑:

根据orgName获得orgId

去user表里选择同时满足orgId和name的用户,如果有,返回存在,如果没有,返回不存在

#### 创建
请求URL: /api/v1/user/new
请求头: Authorization:SessionId
请求方法: POST

携带数据:
```
{
    "userName": "xxx",
    "password": "xxx",  // 暂时有默认值
    "orgName": "dev",       // 创建用户时选择
    "orgId": "1",          // 表示创建者所在的组织, 用来验证管理员会话 
    "op": "1"           // 管理员userId
}
```

#### 用户列表
请求URL: /api/v1/user
请求头: Authorization:SessionId
请求方法: GET 

返回数据:
```
{
    "code": 0,
    "msg": "...",
    "data": {
        "users": [{
            "id": 1,
            "name": "abc.de",
            ....
        }] 
        "orgList": [{
            "orgId": orgName,    // map例如: "1": "dev"
        }]
    }
}
```

列表显示内容:
ID, 用户名, 所属组织, 创建时间, 操作

数据均从data.users的数组里每个元素中取, 其中: 

* ID: id
* 用户名: name
* 所属组织: orgId 
* 创建时间: createdAt
* 操作: 更新、删除

所属组织名称根据orgId获取orgName

#### 删除用户
请求URL: /api/v1/user/delete
请求方法: POST
请求头: Authorization:SessionId
携带数据:
```
{
    "op": 1,
    "userName": "xxx"
}
```

#### 更新用户
请求URL: /api/v1/user/update
请求方法: POST
请求头: Authorization:SessionId

携带数据:
```
{
    "op": 1,              // 管理员userId
    "name": "xxx",
    "orgId": "3",         // 管理员所属orgId
    "password": "xxx",    // 更新的密码
    "orgName":  "xxx",    // 更新的组织, 目前不支持
}
```
