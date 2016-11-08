<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

管理用户列表
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
###目的
由管理员以列表形式管理用户

###请求

* 请求方法: GET 
* 请求URL: /api/v1/user
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 



###页面设计 
无


###程序实现逻辑
```Title: 用户列表 
YCE-->>MySQL: 请求所有可用的用户  
YCE<<--MySQL: 返回请求结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。
JSON:
```json
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

### 备注
无