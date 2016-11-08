<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

镜像管理
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
提供镜像管理功能


###请求

* 请求方法: GET
* 请求URL: /api/v1/registry/images 
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 


###页面设计 
无


###程序实现逻辑
组织名具有全局唯一性
```Title: 检查组织重名
YCE-->>Docker Registry: 查询镜像仓库里所有的镜像
YCE<<--Docker Registry: 返回查询结果
```


###响应数据结构: 
返回码为0, 表示可用。
其他返回码表示出错。
JSON
```json
{
    "code": 0,
    "msg": "...",
    "data": {
        "organizations": [{
            "id": 1,
            "name": "xxx",
            ... 
        }],
        "dcList": [{
            "dcId": dcName,    //  map例如: "1": "电信" 
        }]
        
    }
}
```

### 备注
无



### 以下为旧版本, 无效///////////////////////////////////////////////////

镜像管理相关
======

### 查询所有镜像

在发布镜像、回滚、滚动升级等需要填入镜像时都统一调用这个URL

请求地址: GET /api/v1/registry/images

请求头中包含: Authorization: ${sessionId}

返回值:

* 镜像列表

* 每个镜像的名称和标签列表

Json格式如下:

```json
[
    { 
         "name": "ubuntu",
         "tags": [
            "14.04",
            "15.04",
            "latest"
         ]
    
    },
    {
         "name": "nginx",
         "tags": [
            "1.3.9",
            "1.7.1",
            "1.7.9"
         ]
    }
]

```