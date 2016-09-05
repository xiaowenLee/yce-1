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