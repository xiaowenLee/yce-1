<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

删除组织
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
由管理员删除组织(除了数据库记录外,还有k8s的命名空间)


###请求

* 请求方法: POST
* 请求URL: /api/v1/organization/delete
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "orgName": "xxx",
    "orgId": "3",
    "op": 1
}
```

###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>K8s: 删除对应的命名空间
YCE<<--K8s: 返回删除结果
YCE-->>MySQL: 将organization表中对应记录的status改为INVALID
YCE<<--MySQL: 返回查询结果
```


###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。


### 备注
无