<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

数据中心重名检测
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-10-27

目录
--------------
###目的
根据用户输入的数据中心名称, 判断该数据中心名称是否重名, 即是否已经有数据中心占用了此名称.

###请求

* 请求方法: POST 
* 请求URL: /api/v1/datacenter/check
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON 
  
```json
   
{
    "name": "xxx",
    "orgId": "1"          //表示创建者(管理员)所在的组织,用来验证管理员会话, 从本地存储中获取
}
```

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: NodePort占用列表
YCE-->>MySQL: 按数据中心名称进行查询
YCE<<--MySQL: 返回查询结果
```

说明: 收到POST请求, 去数据库里查询该数据中心名称是否存在, 如果不存在, 返回可用。 如果存在且为VALID, 返回可用, 如果存在且为INVALID, 返回已存在不可用

注: 还应该检查Datacenter, Host:Port对是否重复. 重复了会导致查询上的问题,但不影响实际部署 

###响应数据结构: 
返回码为0, 表示可用。
返回码为1415, 表示资源名重复.

### 备注
无
