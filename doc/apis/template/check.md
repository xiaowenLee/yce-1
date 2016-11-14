<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

重名检查
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-14

目录
--------------
###目的
由用户创建模板时进行重名检查

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{orgId}/users/{userId}/templates/check
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
   "name": "xxx" //templateName
}
```


###页面设计 
无


###程序实现逻辑
```Title: 
重名检查 
YCE-->>MySQL: 在template表中查询该模板名未被使用
YCE<<--MySQL: 返回查询结果 
```

###响应数据结构: 
返回码为0, 表示可用。
返回码1415表示出错。

### 备注