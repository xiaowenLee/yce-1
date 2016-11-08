<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

创建用户准备
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
###目的
为管理员创建用户做准备

###请求

* 请求方法: GET
* 请求URL: /api/v1/user/init
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 



###页面设计 
无


###程序实现逻辑
```Title: 创建用户准备
YCE-->>MySQL: 请求查询可用组织列表 
YCE<<--MySQL: 返回查询结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。
JSON:
```json
{
    "code": 0,
    "msg": "xxx",
    "data": [
            "dev",
            "ops"
        ] 
}
```

### 备注
无