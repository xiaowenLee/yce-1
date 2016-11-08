<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

为创建访问点做准备
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
为创建访问点准备

###请求

* 请求方法: GET
* 请求URL: /api/v1/organizations/{orgId}/users/{userId}/endpoints/init 
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 


###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>MySQL: 查询该组织下的数据中心列表 
YCE<<--MySQL: 返回查询结果
```


###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。
JSON
```
{
      "code": 0,
      "message": "请求成功",
      "data": {
                "orgId":  "1",
                "orgName": "ops",
                "dataCenters": [
                {
                    "dcId": "1",
                    "name": "世纪互联",
                    "budget": 10000000,
                    "balance": 10000000
                },
                {
                    "dcId": "2",
                    "name": "电信机房",
                    "budget": 10000000,
                    "balance": 10000000
                },
                {
                    "dcId": "3",
                    "name": "电子城机房",
                    "budget": 10000000,
                    "balance": 10000000
                }
                ],
      }
   } 
```


### 备注
无