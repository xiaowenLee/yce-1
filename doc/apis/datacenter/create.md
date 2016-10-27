<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

添加数据中心
==============

Author: [maxwell92](github.com/maxwell92)

Last Revised: 2016-10-27

Content
--------------
###目的
由管理员添加新的数据中心, 注意这个数据中心必须是已经部署了Kubernetes的集群 

###请求

* 请求方法: POST 
* 请求URL: /api/v1/datacenter/new
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON
```json
  {
    "name": "xxx",
    "nodePort": [
        "30000",
        "32767"
    ],
    "host": "192.168.1.110",
    "port": 8080,
    "orgId": "3",          // 表示创建者所在的组织, 用来验证管理员会话 
    "op": 1           // 管理员datacenterId
    //"secret": xxx       // 暂时空接口
  }
```

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 添加数据中心 
YCE-->>MySQL: 插入新的数据中心记录
YCE<<--MySQL: 返回插入结果
```

说明: 收到POST请求, 将用户填写的数据中心信息写入MySQL里。写入前进行检查, 如果该名字已经存在,将用新的信息覆盖,并更改status为VALID。如果该名字不存在, 则插入新记录。 

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注
无
