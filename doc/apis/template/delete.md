<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

删除模板
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-14

目录
--------------
###目的
由普通用户删除模板

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{orgId}/users/{userId}/templates/delete
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
  "name": "xxx", // templateName
  "id": 1,    // templateId
}
```


###页面设计 
无


###程序实现逻辑
```Title: 
删除模板 
YCE-->>MySQL: 在template表中将对应记录的status变为INVALID
YCE<<--MySQL: 返回删除结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注