<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

创建服务
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
由用户创建服务

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{orgId}/users/{userId}/services/new
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 
JSON
```json
    {
      "serviceName": "test-service-withselector",
      "orgName": "ops",
      "dcIdList": [1],
      "service": {
                     "kind": "Service",
                     "apiVersion": "v1",
                     "metadata": {
                         "name": "1-test-nginx-service",
                         "labels": {
                             "name": "1-test-nginx-service",
                             "type": "service"
                         }
                     },
                     "spec": {
                         "type": "NodePort",
                         "selector": {
                             "name": "test-nginx-test"
                         },
                         "ports": [
                             {
                                 "protocol": "TCP",
                                 "port": 30091,
                                 "targetPort": 80,
                                 "nodePort": 32289
                             }
                         ]
                     }
                 }
    }
```


###页面设计 
无


###程序实现逻辑
```Title: 创建组织
YCE-->>k8s: 每个k8s集群创建service
YCE<<--K8s: 返回创建结果
```

###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。



### 备注


### 以下为旧版本, 无效///////////////////////////////////////////////////

### 服务发布准备

用户点击服务发布(左侧菜单)时请求后台数据:

请求的URL:

GET /api/v1/organizations/{orgId}/users/{userId}/services/init

请求头中包含: Authorization: ${sessionId}

其中: userId, orgId, sessionId在登录成功后从后台返回给浏览器, 前端存储在LocalStorage(目前为SessionStorage)里面

返回值:

* 组织名称

* 该组织下的数据中心列表



获取成功的大概数据格式如下:

```json

   {
      "code": 0,
      "message": "请求成功",
      "data": {
                "orgId":  "1",
                "orgName": "Ops",
                "nodePort": 30000,  
                "dataCenters": [
                {
                    "dcId": "1",
                    "name": "世纪互联",
                    "budget": 10000000,
                    "balance": 10000000
                },
                {
                    "dcId": "2",
                    "name": "电信机房",
                    "budget": 10000000,
                    "balance": 10000000
                },
                {
                    "dcId": "3",
                    "name": "电子城机房",
                    "budget": 10000000,
                    "balance": 10000000
                }
                ],
      }
   } 
    
```

这些将关系到用户部署服务到哪个机房


用户需要确定:

* 服务名:
* 服务所属数据中心: 依据GET请求的返回值确定。卡片式
* 服务类型:
    仅允许ClusterIP和NodePort两种类型，单选

    1. ClusterIP: 不允许用户填写NodePort或NodePort为0
    2. NodePort: 允许用户填写NodePort

* 选择器(selector): 
  开关。如果打开,用户需要填写一条记录。如果关闭
* 端口组:
  多条:
  
	|名称  |协议  |类型  |端口号 |
|:---:|:----:|:---:|:----:|
|     | tcp  | port| 80|
|  |          |targetPort| 8080|
|  |            |*nodeport*| *80*|

* 标签组
  多条:
  
  KEY:VALUE

### 服务发布提交
用户点击提交请求后台数据:

请求的URL:

POST /api/v1/organizations/{orgId}/users/{userId}/services/new

请求头中包含: Authorization: ${sessionId}

```json
    {
      "serviceName": "test-service-withselector",
      "orgName": "ops",
      "dcIdList": [1],
      "service": {
                     "kind": "Service",
                     "apiVersion": "v1",
                     "metadata": {
                         "name": "1-test-nginx-service",
                         "labels": {
                             "name": "1-test-nginx-service",
                             "type": "service"
                         }
                     },
                     "spec": {
                         "type": "NodePort",
                         "selector": {
                             "name": "test-nginx-test"
                         },
                         "ports": [
                             {
                                 "protocol": "TCP",
                                 "port": 30091,
                                 "targetPort": 80,
                                 "nodePort": 32289
                             }
                         ]
                     }
                 }
    }
```


### 服务修改和删除

按钮在服务及访问点管理的页面上。


