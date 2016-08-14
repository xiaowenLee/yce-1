登录页面数据交互说明
============

### 点击登录按钮

请求的URL: POST /api/v1/users/{email/login

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
        "token": "sfssfd-afds-asdf-af32s"
    }
}
```
