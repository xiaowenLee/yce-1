发布应用重名检测
=============

请求URL: POST /api/v1/organizations/{:orgId}/users/:userId/deployments/check
请求头: Authorization:SessionId

携带数据:
```
{
    "name": "xxx",
}
```

返回数据: code为1415, 表示该名称已被占用,不可使用; code为0表示该名称未被占用, 可以使用。

思路:

去每个数据中心获取所有的deployment的名称, 分别比较, 一旦发现重名即返回1415
