<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

检查用户名重名
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
###目的
管理员创建用户时检查用户名是否重复

###请求

* 请求方法: POST
* 请求URL: /api/v1/user/check
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "userName": "xxx",
    "orgName": "yyy",   // 
    "orgId": "1"          //表示创建者(管理员)所在的组织,用来验证管理员会话, 从本地存储中获取
}
```


###页面设计 
无


###程序实现逻辑
```Title: 创建组织
YCE-->>MySQL: 在user表根据组织和可用用户名查询是否重复 
YCE<<--MySQL: 返回查询结果 
```

###响应数据结构: 
返回码为0, 表示可用。
返回码为1415表示用户名已存在

### 备注
无