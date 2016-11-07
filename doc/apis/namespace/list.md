<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

管理组织列表
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
由管理员删除组织(除了数据库记录外,还有k8s的命名空间)


###请求

* 请求方法: GET
* 请求URL: /api/v1/organization
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 


###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>MySQL: 查询得到所有可用的组织
YCE<<--MySQL: 返回查询结果
```


###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。
JSON
```json
{
    "code": 0,
    "msg": "...",
    "data": {
        "organizations": [{
            "id": 1,
            "name": "xxx",
            ... 
        }],
        "dcList": [{
            "dcId": dcName,    //  map例如: "1": "电信" 
        }]
        
    }
}
```

### 备注
无