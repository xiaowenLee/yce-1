<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

创建模板
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-15

目录
--------------
###目的
由普通用户创建(或保存)模板

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{orgId}/users/{userId}/templates/new
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
创建模板
YCE-->>MySQL: 在template表中插入一条数据  
YCE<<--MySQL: 返回插入结果 
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。

### 备注

template表设计

主键: 自增长id
外键: orgId
模板名与应用名一致,没有的时候与服务名一致


|列:           |  数据类型：| 说明:   |  示例:       |
|:------------:|:------- :|:-------:|:-----------:|
|ID            | INT      | 自增主键 | 1           |
|name          | VARCHAR  | 和deployment名一致 | tomcat-gw   |
|orgId         | INT      | 外键     | 1           |
|status        | INT      |         |  1          |
|deployment    | VARCHAR  | json    |             |
|service       | VARCHAR  | json    |             |
|endpoints     | VARCHAR  | 预留json |             |
|createdAt     | VARCHAR  |         |             |
|modifiedAt    | VARCHAR  |         |             |
|modifiedOp    | VARCHAR  |         |             |
|comment      | VARCHAR  |         |             |


DAO:
```golang
    type Template struct {
        Id  int32 `json:"id"` 
        Name string `json:"name"
        OrgId int32 `json:"orgId"`
        Status int32 `json:"status"`
        Deployment string `json:"deployment"`
        Service string `json:"service"`
        Endpoints string `json: "endpoints"`
        CreatedAt string `json:"createdAt"` 
        ModifiedAt string `json:"modifiedAt"`
        ModifiedOp int32 `json:"modifiedOp"`
        Comment string `json:"comment"`
    }

```

两个按钮, 保存以及取消。

允许只存储deployment或service。在模板管理内将依照deployment.appName和service.serviceName作为判断是否存在对应模板的依据。