<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

NodePort占用列表
==============

Author: [maxwell92](github.com/maxwell92)

Last Revised: 2016-10-26

Content
--------------
####目的
为管理员列出所有被占用的NodePort. 被占用的含义是该nodePort在数据库表中的记录存在且status为INVALID.

####请求

* 请求方法: GET
* 请求URL: /api/v1/nodeports
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 无

####页面设计 
管理页面样式(参考应用管理). 显示表头: ID, NodePort, 数据中心, 服务名称, 创建时间, 其中:

* ID: id, 按序生成
* NodePort: data.nodePorts[0].port
* 数据中心: 根据data.nodePorts[0].dcId去data.dcList[]去对应的文字
* 服务名称: data.nodePorts[0].svcName
* 创建时间: data.nodePorts[0].createdAt 

####程序实现逻辑:

```Sequence
Title: NodePort占用列表
YCE-->>MySQL: 查询status为INVALID的nodePort
YCE<<--MySQL: 返回查询结果
YCE-->>MySQL: 查询dcId对应的dcName
YCE<<--MySQL: 返回查询结果
```

说明: 收到GET请求, 去数据库里查询所有status为INVALID的nodePort,并查询dcId与dcName的对应关系, 两者共同作为结果返回

####响应数据结构: 

JSON, 示例如下:

```json
  {
      "nodePorts": [
          {
              "port": 30080,
              "dcId": 1,
              "svcName": "play2048",
              "status": 0,
              "createdAt": "2016-10-18T10:01:04+08:00",
              "ModiAt": "2016-10-18T10:01:04+08:00",
              "modifiedOp": 7,
              "comment": ""
          },
          {
              "port": 32080,
              "dcId": 1,
              "svcName": "yce-test.ycetest",
              "status": 0,
              "createdAt": "2016-08-26T09:40:57+08:00",
              "ModifiedAt": "2016-08-26T09:52:27+08:00",
              "modifiedOp": 1,
              "comment": "yce-test service"
          }
      ],
      "dcList": {
          "1": "办公网",
          "2": "电信机房"
      }
  } 

```
#### 备注
无
