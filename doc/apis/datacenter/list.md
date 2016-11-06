<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

数据中心列表
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-10-27

目录
--------------
###目的
为管理员列出所有的可用数据中心列表。


###请求

* 请求方法: GET 
* 请求URL: /api/v1/datacenter
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 无


###页面设计 
ID, 数据中心名, 地址, 端口, NodePort, 创建时间, 操作

数据均从data.datacenters的数组里每个元素中取, 其中: 

* ID: id
* 数据中心名: name
* 地址: host
* 端口: port
* NodePort: nodePort[0] - nodePort[1] //暂定
* 创建时间: createdAt
* 操作: 更新、删除

###程序实现逻辑:

```Sequence
Title: 数据中心列表 
YCE-->>MySQL: 请求可用的数据中心列表
YCE<<--MySQL: 返回请求结果
```

说明: 收到GET请求, 在数据库里查找所有status为VALID的数据中心, 并组成列表返回 

###响应数据结构: 
JSON
```json
{
    "datacenters": [
        {
            "id": 1,
            "name": "湖南",
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
            "modifiedAt": "2016-08-15T16:27:30Z",
            "modifiedOp": 1,
            "comment": ""
        }
    ]
}
```

返回码为0, 表示可用。
其他返回码表示出错。

### 备注
无


