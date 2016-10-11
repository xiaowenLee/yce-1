### 创建用户
-----------

#### 创建初始化
目的: 检查用户名是否重复
请求URL: /api/v1/organizations/1/user/init //organizations 表示创建者(管理员)所在的组织,用来验证管理员会话
请求头: Authorization:SessionId
请求方法: POST

携带数据:
```
{
    "name": "xxx",
    "orgName": "yyy"
}
```

返回在该组织里是否存在, "code": 1414 为已存在, "code": 0为未被占用.

程序实现逻辑:

根据orgName获得orgId

去user表里选择同时满足orgId和name的用户,如果有,返回存在,如果没有,返回不存在

#### 创建
请求URL: /api/v1/organization/1/user/new
请求头: Authorization:SessionId
请求方法: POST

携带数据:
```
{
    "name": "xxx",
    "password": "xxx", // 暂时有默认值
    "orgName": 1,       // 创建用户时选择
    "createdOp": 1
}
```

