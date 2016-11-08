<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

用户退出
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
用户退出系统,注销会话

###请求

* 请求方法: POST
* 请求URL: /api/v1/users/logout
* 请求头: Authorization:$SessionId, 从LocalStorage读  
* 请求参数: 


###页面设计 
无


###程序实现逻辑
无

###响应数据结构: 
JSON
```json
{
   "code": 0,
    "message": "ok",
    "data": ""
}
```


### 备注
用户退出应该有传递参数过去


### 以下为旧版本, 无效///////////////////////////////////////////////////

退出登录数据交互说明
======================

### 点击退出按钮

//请求的URL: POST /api/v1/users/{username}/logout
请求URL: POST /api/v1/users/logout
 
请求头中: Authorization: ${sessionId}

返回值:

* 返回码: 是否退出成功,退出功能不做是否成功检查

* 出错信息

浏览器接收到返回值后清空本地缓存

服务器端清除该session

数据格式如下:

```json
{
    "code": 0,
    "message": "ok",
    "data": ""
}
```