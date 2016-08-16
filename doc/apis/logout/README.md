退出登录数据交互说明
======================

### 点击退出按钮

请求的URL: POST /api/v1/users/{username}/logout
 
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