dashboard显示信息:
--------
用户登录后默认进入Dashboard, 此时向后台发送请求并显示相关图表,需要一个默认刷新时间

### 获取dashboard统计信息

目前将dashboard划分为三块, 包括:

* 配额/已用比, 双环饼图。不同的数据中心不同的饼。可以看该组织的配额比例(Overall), 剩下的显示按数据中心使用的。
* 项目(Deployment)拓扑统计, rs->pods示例, 在卡片上的生成树, 不同的数据中心不同的条带, 每条里多个卡片。默认显示最近的5个(最终以屏幕显示方便决定)deployment。高级选项可以由用户定义Pinned卡片。
* 操作统计, 柱状图。默认显示所有数据中心总计操作统计,高级功能:可以选择看不同的数据中心,也可以选择看该数据中心下具体的deployment相关操作。还可以按时间轴进行放缩。

### API及数据结构设计

#### 资源配额/已用
注: 暂时不考虑云盘
请求URL: GET /api/v1/organizations/{:orgId}/resourcestat
请求头: Authorization: SessionId

需要数据: 该组织的budget和balance, 每个数据中心的budge和balance
organization里存放的有buget和balance以及CPU和MEM各自的总额。另外存放的是dcIdList。应该加上用了多少。
datacenters里没有总额和现在已用多少的资源配额。 应该加一条。
```
{
    "code": 0,
    "msg": "",
    "data": {
        "overall": {
            "cpu": {
                "quota": 10,
                "used": 6
            }, 
            "mem": {
                "quota": 10,
                "used": 4
            }
        },
        "datacenters": [{
            "dcId": 1,
            "dcName": "办公网",
            "cpu": {
                "quota": 10,
                "used": 6
            }, 
            "mem": {
                "quota": 10,
                "used": 4
            }
        }]
    }
}
```



该组织下所有数据中心配额比例总计:

请求URL: GET /api/v1/organizations/{:orgId}/resourcestat
请求头: Authorization: SessionId

返回数据定义:

该组织下某数据中心配额比例:
请求URL: GET /api/v1/organizations/{:orgId}/datacenters/{:dcId}/resourcestat
请求头: Authorization: SessionId

返回数据定义:

#### 项目拓扑卡片
请求URL: GET /api/v1/organizations/{:orgId}/deploymentstat
请求头: Authorization: SessionId

需要数据: 该数据中心下的deployment中,最近的rs和示例数量(仅)?
获取该数据中心下的全部deployment,再找到NewRS,再找到pod。
返回数据定义:
```
{
    "code": 0,
    "msg": "",
    "data": [{
        "dcId": 1,
        "dcName": "办公网",
        "deployments": [{
            "deploymentName": "yce",
            "rsName": "test-yce",
            "podName": [
                "test-yce-gxbvs",
                "test-yce-sdgsg"
            ] 
        }]
    }]
}
```




该数据中心下默认最近前5个项目:
请求URL: GET /api/v1/organizations/{:orgId}/datacenters/{:dcId}/defaultdeployment
请求头: Authorization: SessionId

返回数据定义:

高级功能:

该数据中心下该用户Pinned的5个项目
请求URL: GET /api/v1/organizations/{:orgId}/datacenters/{:dcId}/pinneddeployment
请求头: Authorization: SessionId

返回数据定义:


#### 操作统计
请求URL: GET /api/v1/organizations/{:orgId}/datacenters/{:dcId}/operationstat
请求头: Authorization: SessionId

需要数据: 该数据中心下: 
按数据中心在deployment里进行选择, 然后再统计不同操作的出现频次。现在存放的是数据中心列表dcIdList不好筛选


返回数据定义:
```
{
    "code": 0,
    "msg": "",
    "data": [{
        "dcId": 1,
        "dcName": "办公网",
        "create": [1,2,3,4,5,6,7],
        "scale": [2,3,4,5,6,7,8],
        "rollingupgrade": [1,2,3,4,5,6,7],
        "rollback": [1,2,3,4,5,6,7]
    }]

}
```


该组织下所有数据中心的操作统计:
请求URL: POST /api/v1/organizations/{:orgId}/operationstat
请求头: Authorization: SessionId

发送数据为时间区间, 默认为1个月:
```
{
   period: 1 
}
```

返回数据定义:
```
{
    "code": 0,
    "msg": "",
    "data": [{
        "dcId": 1,
        "create": 1,
        "scale": 2,
        "rollingupgrade": 3,
        "rollback": 4
    }]

}
```


高级选项:
该组织下某数据中心的操作统计:
请求URL: POST /api/v1/organizations/{:orgId}/datacenters/{:dcId}/operationstat
请求头: Authorization: SessionId

发送数据为时间区间, 默认为1个月:

返回数据定义:

该组织下某项目(deployment)的操作统计:
请求URL: POST /api/v1/organizations/{:orgId}/deployments/{:deployName}/operationstat
请求头: Authorization: SessionId

发送数据为时间区间以及数据中心Id, 默认为1个月

返回数据定义:



