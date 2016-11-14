<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

更新模板
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-14

目录
--------------
###目的
由普通用户更新模板

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{orgId}/users/{userId}/templates/update
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
   "name": "xxx",
   "deployment": { }, //创建应用时生成的json
   "service": { } //发布服务时生成的json
}
```


###页面设计 
无


###程序实现逻辑
```Title: 
更新模板 
YCE-->>MySQL: 在template表中找到该记录并更新
YCE<<--MySQL: 返回更新结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注
在管理模板页面上有个更新,点击之后跳转到创建模板页,里面将导入该模板的信息,由用户进行修改。可以保存为这个模板,也可以保存为新模板。 两个按钮, 保存和取消。
如果是保存, 判断模板名是否已经有,有则更新对应记录,没有就插入新记录 。

更新采用的是重复时更新,否则插入。并且将模板名设置为唯一索引。
