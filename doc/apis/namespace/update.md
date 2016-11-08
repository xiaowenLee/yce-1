<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

更新组织
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
由管理员更新组织信息

###请求

* 请求方法: POST
* 请求URL: /api/v1/organization/update
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "userId": xxx,        // 操作者的用户ID,这里是数字,不是字符串, 管理员
    "name":  "xxx",       // 组织名称
    "orgId": "xxx",       // 操作者的组织ID(管理员),不是新创建的那个组织ID
    "cpuQuota": 100,      // CPU配额
    "memQuota": 200,      // 内存配额
  //目前只返回上面内容的json, 即只更新配额 
    "dcIdList": [1, 2]    // 数据中心列表
    "budget": 10000,      // 预算
    "balance": 10000,     // 余额
}
```

###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>MySQL: 按照新的信息更新organization表里对应的记录
YCE<<--MySQL: 返回更新结果
```


###响应数据结构: 
"code": 1415 为用户名已存在, 不能使用该名称, 需提示。 
"code": 0为未被占用, 可以使用该名称, 无需提示。


### 备注
资源配额(套餐)已经在组织管理返回了quotaPkg
