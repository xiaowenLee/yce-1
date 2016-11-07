<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

镜像列表
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
获取镜像列表,为创建应用等做准备

###请求

* 请求方法: GET 
* 请求URL: /api/v1/registry/images
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 

###页面设计 
无


###程序实现逻辑
无

###响应数据结构: 
```json
{
    "code":0,
    "message":"OK",
    "data":"[{"name":"busybox","tags":["v2.0","v1.0","v3.0","latest"]},{"name":"golang","tags":["1.6.2","latest"]},{"name":"memcached","tags":["1.4.24"]},{"name":"mysql","tags":["5.6"]},{"name":"nginx","tags":["1.7.9"]},{"name":"tomcat7","tags":["latest"]},{"name":"ubuntu","tags":["14.04"]}]"
}
```

### 备注
image的json结构：

```json
  {
    "name": "",
    "tags": [
      ""
    ]
  }

```


### 以下为旧版本, 无效///////////////////////////////////////////////////

##镜像管理
==========

用户点击镜像管理时请求后台数据:


请求的方法及URL: GET /api/v1/registry/images
请求头中包含: Authorization: ${sessionId} *暂时在Session Storage里*

返回值:

* 该组织下组织下的镜像列表

返回json示例：

```json
{
    "code":0,
    "message":"OK",
    "data":"[{"name":"busybox","tags":["v2.0","v1.0","v3.0","latest"]},{"name":"golang","tags":["1.6.2","latest"]},{"name":"memcached","tags":["1.4.24"]},{"name":"mysql","tags":["5.6"]},{"name":"nginx","tags":["1.7.9"]},{"name":"tomcat7","tags":["latest"]},{"name":"ubuntu","tags":["14.04"]}]"
}
```

image的json结构：

```json
  {
    "name": "",
    "tags": [
      ""
    ]
  }

```
