<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

Dashboard资源统计信息
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-06

目录
--------------
###目的
在Dashboard上显示所有数据中心的资源使用情况。


###请求

* 请求方法: GET 
* 请求URL: /api/v1/organizations/{:orgId}/resourcestat
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 无


###页面设计 

以双环饼图显示,其中外环为CPU使用,内环为内存使用。
第一个饼图为总览, 统计所有数据中心的资源使用。后续的按数据中心分别统计。


###程序实现逻辑:

```Sequence
Title: 数据中心列表 
YCE-->>MySQL: 获取该组织所有的资源配额及对应的可用数据中心列表
YCE<<--MySQL: 返回请求结果
YCE-->>K8S: 获取各个应用已使用的资源情况 
YCE--<<K8S: 返回请求结果
```

###响应数据结构: 
JSON
```json
{
       "code": 0,
       "msg": "操作成功",
       "data": [{                      // 总览
               "dcId": 0,
               "dcName": "总览",
               "cpu": {
                   "total": 20,         // 总额
                   "used":  6           // 已用
               }, 
               "mem": {
                   "total": 64,
                   "used": 24
               }
           },
           {                            // 分数据中心显示
               "dcId": 1,
               "dcName": "办公网",
               "cpu": {
                   "total": 10,
                   "used": 2
               }, 
               "mem": {
                   "total": 32,
                   "used": 8 
               } 
           },
           {
               "dcId": 2,
               "dcName": "电信",
               "cpu": {
                   "total": 10,
                   "used": 4
               }, 
               "mem": {
                   "total": 32,
                   "used": 16 
               }
           }
       }]
}
```

### 备注
organization里存放的有CPU和MEM各自的总额。total会存在数据库里, used从k8s里读出动态计算.目前假定同一个组织下多个数据中心的配额均相同。
