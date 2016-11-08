<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

服务创建初始化
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-10-27

目录
--------------
###目的
为创建服务做准备, 返回数据中心列表以及推荐的nodePort等

###请求

* 请求方法: GET 
* 请求URL: /api/v1/organizations/:orgId/users/:userId/services/init
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 无 


###页面设计 
无


###程序实现逻辑:
```Sequence
Title: 添加数据中心 
YCE-->>MySQL: 查询可用数据中心记录, 拼成列表 
YCE<<--MySQL: 返回查询结果
YCE-->>MySQL: 查询推荐的nodePort
YCE<<--MySQL: 获取可用的nodePort
```

说明: 
去数据库里获取所有可用的数据中心拼接为列表, 并获取推荐的nodePort值

依据现在的nodePort推荐算法设计, 当选中单个数据中心时,返回该数据中心第一个可用的nodePort。当选中多个数据中心时,如果它们有交集,返回交集里第一个可用nodePort。
如果没有选择数据中心, 不推荐nodePort。

###响应数据结构: 
JSON
```json
{
    "orgId": "1",
    "orgName": "dev",
    "dataCenters": [
        {
            "id": 1,
            "name": "办公网",
            "host": "172.21.1.11",
            "port": 8080,
            "secret": "",
            "status": 1,
            "nodePort": {
                "nodePort": [
                    "30000",
                    "32767"
                ]
            },
            "createdAt": "2016-08-15T16:27:30Z",
            "modifiedAt": "2016-10-18T10:23:35    +08:00",
            "modifiedOp": 1,
            "comment": ""
        }
    ],
    "nodePort": {
        "port": 32085,
        "dcId": 1,
        "svcName": "mark",
        "status": 1,
        "createdAt": "2016-10-18T16:25:36+08:00",
        "M    odifiedAt": "2016-10-21T16:52:24+08:00",
        "modifiedOp": 8,
        "comment": ""
    }
}
```


### 备注
无







