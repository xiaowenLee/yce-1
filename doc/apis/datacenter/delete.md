<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

删除数据中心
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-10-27

目录
--------------
###目的
由管理员删除该数据中心信息, 即将其变为不可用。


###请求

* 请求方法: POST 
* 请求URL: /api/v1/datacenter/delete
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
    JSON
```json
{
    "op": 1,          // 管理员userId
    "orgId": "3",     // 管理员所属组织orgId
    "name": "xxx"
}
```

###页面设计 
无

###程序实现逻辑:

```sequence
Title: 删除数据中心 
YCE-->>MySQL: 请求将该数据中心变为不可用 
YCE<<--MySQL: 返回请求结果
YCE-->>MySQL: 请求将该数据中心里的组织变为不可用
YCE<<--MySQL: 返回请求结果
YCE-->>MySQL: 请求将该数据中心里的组织下设的用户变为不可用
YCE<<--MySQL: 返回请求结果
YCE-->>MySQL: 请求将该数据中心里的nodePort变为不可用
YCE<<--MySQL: 返回请求结果
YCE-->>K8s: 删除该数据中心里所有的资源(namespace、limitrange、resoursquota、deployment、service、endpoints)
YCE<<--K8s: 返回删除结果
```

说明: 收到POST请求, 去数据库里按照上述列表将对应的记录改为不可用。对于dcIdList里包含此数据中心的记录(利用JSON_EXTRACT进行检查), 利用UPDATE + JSON_SET进行更新. 对于只在这个数据中心里的记录, 将其变为不可用INVALID。

###响应数据结构: 

返回码为0, 表示可用。
其他返回码表示出错。

### 备注
目前只实现了将该数据中心变为不可用、以及将该数据中心里的nodePort变为不可用。
