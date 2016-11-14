<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

另存为
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-11

目录
--------------
###目的
用户在管理模板的时候会更新某个模板中的部分信息,并保存为新模板,即另存为

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organization/{orgId}/users/{userId}/templates/new
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
{
   "name": "xxx",
   "deployment": { },//创建应用时生成的json
   "service": { } //发布服务时生成的json
}
```


###页面设计 
无


###程序实现逻辑
```Title: 
另存为(同创建模板) 
YCE-->>MySQL: 在template表中插入一条数据  
YCE<<--MySQL: 返回插入结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注
另存为的时候要求用户必须输入新的模板名称
暂无需实现