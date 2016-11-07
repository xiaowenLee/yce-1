<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

用户登录
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
用户登录

###请求

* 请求方法: POST
* 请求URL: /api/v1/users/login
* 请求头: 
* 请求参数: 
JSON
```json
{
  "username": "xxx",
  "password": "xxx"

}
```

###页面设计 
无


###程序实现逻辑
无

###响应数据结构: 
JSON
```json
{
    "code": 0,
    "message": ""
    "data": {
        "userId": "12",
        "userName": "liyao.miao",
        "orgId": "2",
        "sessionId": "sfssfd-afds-asdf-af32s"
    }
}
```


### 备注
无


### 以下为旧版本, 无效///////////////////////////////////////////////////

登录页面数据交互说明
============

### 点击登录按钮

请求的URL: POST /api/v1/users/login

数据通过表单提交: username=${username}  password=${password}

返回值:

* 返回码:是否通过登录验证

* 出错信息

* 用户的ID

* 用户的姓名

* 用户所在的组织ID

* 用户的访问令牌(使用angular本地缓存)

数据格式如下:

```json
{
    "code": 0,
    "message": ""
    "data": {
        "userId": "12",
        "userName": "lidawei",
        "orgId": "2",
        "sessionId": "sfssfd-afds-asdf-af32s"
    }
}
```
