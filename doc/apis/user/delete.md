<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

删除用户
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
###目的
由管理员删除用户

###请求

* 请求方法: POST 
* 请求URL: /api/v1/user/delete
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
    "op": 1,
    "userName": "xxx"
}
```


###页面设计 
无


###程序实现逻辑
```Title: 删除用户 
YCE-->>MySQL: 在user表中更新对应用户记录为INVALID  
YCE<<--MySQL: 返回更新结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注
无