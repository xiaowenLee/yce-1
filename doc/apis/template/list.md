<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

管理模板列表
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-14

目录
--------------
###目的
为用户返回模板列表

###请求

* 请求方法: GET 
* 请求URL: /api/v1/organizations/{orgId}/users/{userId}/templates
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
无


###页面设计 
无


###程序实现逻辑
```Title: 
管理模板列表
YCE-->>MySQL: 在template表中请求所有可用数据并返回  
YCE<<--MySQL: 返回请求结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

JSON
```json
{
   "code": xxx,
   "msg": "xxx",
   "data": {
      "templates": [{
        "id": 1,
        "name": "xxx",
        "deployment": { },
        "service": { }  ,
        "createdAt": "xxx",
        "modifiedOp": xx,
      }],
      "users": {
         1: "xxx",
         2: "yyy"
      }
   } 
}
```
### 备注

页面列表:

|信息：         |  说明：|
|:------------:|:--------------:|
|ID            |  数字，为页面显示ID|
|模板名称       |  data[].templates[every].name |
|应用信息       |  data[].templates[every].deployment, 点击看到详情 | 
|服务信息       |  data[].templates[every].service, 点击看到详情 |
|创建时间       |  data[].templates[every].createdAt |
|创建人员       |  data[].users[data[].templates[every].modifiedOp] |

