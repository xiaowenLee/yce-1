<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

为创建组织或更新组织做准备
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
为创建组织或更新组织准备

###请求

* 请求方法: GET
* 请求URL: /api/v1/organization/init 
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 


###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>MySQL: 查询数据中心列表, 账户余额信息, 资源套餐信息 
YCE<<--MySQL: 返回查询结果
```


###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。
JSON
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


### 备注
无