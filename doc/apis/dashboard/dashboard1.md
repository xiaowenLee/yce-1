dashboard显示信息:
--------
用户登录后默认进入Dashboard, 此时向后台发送请求并显示相关图表,需要一个默认刷新时间, 30s?

### 获取dashboard统计信息

目前将dashboard划分为三块, 基本内容包括:

* 资源总额/已用统计, 双环饼图。不同的数据中心不同的饼。可以看该组织的配额比例(Overall), 剩下的显示按数据中心使用的。注: 暂时不考虑云盘
* 应用(Deployment)拓扑统计, rs->pods示例, 在卡片上的生成树, 不同的数据中心不同的条带, 每条里多个卡片。单行不够显示换行。
* 操作统计, 柱状图。默认按数据中心显示最近7天内各种操作统计。

高级内容包括:
* 应用拓扑统计: 高级选项可以由用户定义Pinned卡片。
* 操作统计: 可以选择看该数据中心下具体的deployment相关操作。还可以按时间轴进行放缩。

### 基本内容API及数据结构设计

#### 资源配额/已用
请求URL: GET /api/v1/organizations/{:orgId}/resourcestat
请求头: Authorization: SessionId

需要数据: 该组织(CPU/MEM)的total和used, 每个数据中心的ID和名字,以及分别拥有的资源total和used

备注:
organization里存放的有CPU和MEM各自的总额。total会存在数据库里, used从k8s里读出动态计算.

示例数据:

 ```
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




#### 应用统计
请求URL: GET /api/v1/organizations/{:orgId}/deploymentstat
请求头: Authorization: SessionId

需要数据: 数据中心ID及名字, 和该数据中心下的deployment中,最近的rs名称和pod名称

备注:
获取该数据中心下的全部deployment,再找到NewRS,再找到pod。

示例数据:

```
{
    "code": 0,
    "msg": "操作成功",
    "data": [{
        "dcId": 1, 
        "dcName": "办公网",
        "deployments": [{
            "deploymentName": "test1",
            "rsName": "test1-deployment-3123123004",
            "podName": [
                "test1-deployment-gxbvs",
                "test1-deployment-sdgsg"
            ] 
        },
        {
            "deploymentName": "test2",
            "rsName": "test2-deployment-2650739725",
            "podName": [
                "test2-deployment-gxavs",
                "test2-deployment-gxavg",
                "test2-deployment-gxavj",
                "test2-deployment-gxavk",
                "test2-deployment-sdesl"
            ] 
        },
        {
            "deploymentName": "test3",
            "rsName": "test3-deployment-3123122311",
            "podName": [
                "test3-deployment-sdgsy"
            ] 
        },
        {
            "deploymentName": "test4",
            "rsName": "test4-deployment-3112523004",
            "podName": [
                "test4-deployment-gxbas",
                "test4-deployment-gxxjs",
                "test4-deployment-sdasg"
            ] 
        },
        {
            "deploymentName": "test5",
            "rsName": "test5-deployment-1209123004",
            "podName": [
                "test5-deployment-txbvs",
                "test5-deployment-txbns",
                "test5-deployment-txjsk",
                "test5-deployment-stgsg"
            ] 
        }
        ]
    },
    {
        "dcId": 2,
        "dcName": "电信",
        "deployments": [{
            "deploymentName": "test1",
            "rsName": "test1-deployment-1209112354",
            "podName": [
                "test1-deployment-gxivs",
                "test1-deployment-sdosg"
            ] 
        },
        {
            "deploymentName": "test2",
            "rsName": "test2-deployment-1209123784",
            "podName": [
                "test2-deployment-gxaas",
                "test2-deployment-adgsg"
            ] 
        }]
    }]
}
```




#### 操作统计
请求URL: GET /api/v1/organizations/{:orgId}/operationstat
请求头: Authorization: SessionId

需要数据: 该数据中心下, 最近7天的五种主要操作(上线、滚升、回滚、扩容、删除)统计

备注:
按数据中心在deployment表里进行选择, 然后再统计不同操作的出现频次。现在存放的是数据中心列表dcIdList不好筛选

示例数据:
```
{
    "code": 0,
    "msg": "操作成功", 
    "data": [{
        "dcId": 1,
        "dcName": "办公网",                  // 该数据中心下所有应用的每日操作频率统计, 按日期的顺序显示 
        "online": [2,2,1,4,5,1,7],          // 发布 
        "scale": [2,3,0,5,0,7,8],           // 扩容
        "rollingupgrade": [1,2,1,4,0,6,7],  // 滚动升级 
        "rollback": [1,0,3,7,5,8,7],        // 回滚
        "delete": [0,0,0,0,0,0,0]           // 删除
    },
    {
        "dcId": 2,
        "dcName": "电信",
        "online": [1,2,3,4,0,6,7],
        "scale": [2,3,0,5,6,7,8],
        "rollingupgrade": [0,2,3,4,5,6,7],
        "rollback": [1,2,0,4,5,6,7],
        "delete": [0,0,0,0,0,0,0]           
    }
    ]

}
```




### 高级内容API及数据结构设计

#### 应用统计
高级功能:

该数据中心下该用户Pinned的5个项目
请求URL: GET /api/v1/organizations/{:orgId}/datacenters/{:dcId}/pinned
请求头: Authorization: SessionId

#### 操作统计
选择看该数据中心下具体的deployment相关操作。
该组织下某项目(deployment)的操作统计:
请求URL: POST /api/v1/organizations/{:orgId}/operationstat
请求头: Authorization: SessionId
携带数据:
```
{
   "period": 7,      //为空表示默认7天
   "dcIdList": [1],
   "deploymentName": "test1-abc"
}
```

还可以按时间轴进行放缩。
请求URL: POST /api/v1/organizations/{:orgId}/operationstat
请求头: Authorization: SessionId
携带数据:
```
{
   "period": 1,      
   "dcIdList": [1],
   "deploymentName":  //为空表示全部应用
}
```




