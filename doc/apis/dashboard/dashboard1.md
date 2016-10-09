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




#### 应用统计
请求URL: GET /api/v1/organizations/{:orgId}/deploymentstat
请求头: Authorization: SessionId

需要数据: 数据中心ID及名字, 和该数据中心下的deployment中,最近的rs名称和pod名称

备注:
获取该数据中心下的全部deployment,再找到NewRS,再找到pod。




#### 操作统计
请求URL: GET /api/v1/organizations/{:orgId}/operationstat
请求头: Authorization: SessionId

需要数据: 该数据中心下, 最近7天的五种主要操作(上线、滚升、回滚、扩容、删除)统计

备注:
按数据中心在deployment表里进行选择, 然后再统计不同操作的出现频次。现在存放的是数据中心列表dcIdList不好筛选





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




