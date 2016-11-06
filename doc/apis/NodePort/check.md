<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

NodePort重名检测
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-10-26

目录
--------------
###目的
根据用户输入的nodePort以及选择的数据中心, 判断该nodePort是否重名, 即是否已经有应用占用了此端口.

###请求

* 请求方法: POST 
* 请求URL: /api/v1/nodeports/check
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON 
  
```json
   {
       "nodePort": 32380,
       "dcIdList": [
           1,
           2
       ],
       "orgId": "1"
   } 
```

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: NodePort占用列表
YCE-->>MySQL: 按nodePort和dcId进行查询
YCE<<--MySQL: 返回查询结果
```

说明: 收到POST请求, 去数据库里查询该nodePort和dcId组合是否存在, 如果不存在, 返回可用。 如果存在且为VALID, 返回可用, 如果存在且为INVALID, 返回已存在不可用
用户没有选择数据中心就填写nodePort不提供检查
用户选了多个数据中心的当且仅当所有数据中心里该端口均为可用时方可使用,其余全按不可用

###响应数据结构: 
返回码为0, 表示可用。
返回码为1415, 表示资源名重复.

### 备注
无