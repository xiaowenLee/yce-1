<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

镜像搜索辅助
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
给用户提供创建应用时的镜像列表, 用户可以按照关键字进行搜索

###请求

* 请求方法: GET 
* 请求URL: /api/v1/images
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON
```json
 {
    "code": 0,
    "message": "....",
    "data": [
        {
            "imageName": "",
            "imageTag": "",
            // 其他一些可有可无的信息,第一版不需考虑...
        },
        {
            //...
        }
    ]
}
```

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 发布应用
YCE-->>Registry: 请求获取镜像列表
YCE<<--Registry: 返回请求结果
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

### 备注
无
