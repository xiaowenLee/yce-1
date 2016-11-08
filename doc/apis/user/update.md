<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

更新用户
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
###目的
由管理员更新用户


###请求

* 请求方法: POST 
* 请求URL: /api/v1/user/update
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "op": 1,              // 管理员userId
    "name": "xxx",
    "orgId": "3",         // 管理员所属orgId
    "password": "xxx",    // 更新的密码
    "orgName":  "xxx",    // 更新的组织, 目前不支持
}
```


###页面设计 
无


###程序实现逻辑
```Title: 更新用户 
YCE-->>MySQL: 在user表中更新对应用户记录  
YCE<<--MySQL: 返回更新结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注
普通用户修改个人信息用的是不同的API,但可以用同样的后台
