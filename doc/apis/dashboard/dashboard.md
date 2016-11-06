<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

Dashboard显示信息
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-06

目录
--------------
###目的
在Dashboard上显示资源使用情况、应用情况以及最近操作情况的统计。反应大概情况。暂不支持详情查看。

目前将dashboard划分为三块, 基本内容包括:

* 资源总额/已用统计, 双环饼图。不同的数据中心不同的饼。可以看该组织的配额比例(Overall), 剩下的显示按数据中心使用的。注: 暂时不考虑云盘
* 应用(Deployment)拓扑统计, rs->pods示例, 在卡片上的生成树, 不同的数据中心不同的条带, 每条里多个卡片。单行不够显示换行。
* 操作统计, 柱状图。默认按数据中心显示最近7天内各种操作统计.

###请求
用户登录后默认进入Dashboard, 此时向后台发送请求并显示相关图表,需要一个默认刷新时间, 30s?

见详细设计文档



###页面设计 
见详细设计文档


###程序实现逻辑:
见详细设计文档



###响应数据结构: 
见详细设计文档


### 备注
上面设计的是基本内容API及数据结构等,下面是一些设想中的高级内容API及数据结构设计:

高级内容包括:
* 应用拓扑统计: 高级选项可以由用户定义Pinned卡片。
* 操作统计: 可以选择看该数据中心下具体的deployment相关操作。还可以按时间轴进行放缩。

#### 高级内容API及数据结构设计

##### 应用统计
高级功能:

该数据中心下该用户Pinned的5个项目
请求URL: GET /api/v1/organizations/{:orgId}/datacenters/{:dcId}/pinned
请求头: Authorization: SessionId

##### 操作统计
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

### 以下为旧版本, 无效///////////////////////////////////////////////////
--------

### 获取dashboard统计信息

请求URL: /api/v1/organizations/{:orgId}/dashboard 

请求头: Authorization: SessionId

请求返回Json:

```json
{
  "code":0
  "message": "OK",
  "data": {
      "datacenterStat": {       // 数据中心信息统计
        "totalActionStat": {
            "counts": {
              "online": 1,
              "rollingBack": 1,
              "rollingUpgrade": 1,
              "pause": 1,
              "resume": 1,
              "scaling": 1,
              "delete": 1,
              "cancel": 1
            }
           
      },
      "totalDeploymentStat": { // 当前发布总数
          "current": 10
      },
      "totalService": {       // 当前服务总数
          "current": 5 
      },
      "quotaStat":    {      // 当前配额统计
          "total": {
              "cpu": 10,
              "mem": 200
          },
          "used": {
              "cpu": 5,
              "mem": 100
          }, 
          "free": {
              "cpu": 5,
              "mem": 100,
          }
      }
    },
    
   "userStat": {                // 当前用户统计
        "organizations": "ops",
        "dataCenters": [
              "bangongwang",
              "dianxin"
        ],
        "userActionCounts": {   // 数据中心里用户操作信息
              "ActionType1": 10,
              "ActionType2": 20
           } 
        "accountAudit": [      // 账户收支统计
            "budget": 100,
            "balance": 50
        ]
   },
   
   "deploymentStat": [{         // 该组织下应用统计
          "dcIdList": [
              1,
              2
          ],
          "dcName": [
             "bangongwang",
             "dianxin"
          ],
          "actionStat": {
            "counts": {
              "online": 1,
              "rollingback": 1,
              "rollingupgrade": 1,
              "pause": 1,
              "resume": 1,
              "scaling": 1,
              "delete": 1,
              "cancel": 1,
            },
          }, 
          "health": [      // 健康状态百分比  running / all
              70,
              80
          ],
          "podStat": [     // 实例个数求它的长度即可
              {
                  "podName": "abc",
                  "restart": 10,
                  "state": "Running"
              } 
          ]
   }],
   "serviceStat": [
      {
          "dcIdList": [
              1,
              2
           ],
          "dcName": [
              "bangongwang",
              "dianxin"
          ]
          "svcName": "test-svc",
          "health": 100
      }       
   ]
} 

```


数据包括:
    数据中心统计, 包括:
        总操作统计: 各个操作(上线、回滚等)的次数
        总应用统计: 当前应用数量
        总服务统计: 当前服务数量
        配额统计: 
            已用:
                CPU:
                MEM:
            可用:
                CPU:
                MEM:
            总额:
                CPU:
                MEM:
    用户统计, 包括:
        组织:
        数据中心:
        用户操作统计:
        账户统计:
            预算:
            支出:
    应用统计, 按应用统计, 按数据中心区分
        数据中心:
        数据中心名称:
        操作统计:
        健康度统计: 运行/全部 百分比
        实例统计:
            实例名称:
            运行状态:
            重启次数:
    服务统计, 按数据中心区分
        数据中心:
        数据中心名称:
        服务名称:
        是否可用:


涉及到多个信息组,目前能想到的有:

###普通用户

应用发布信息:

发布次数、滚升次数、回滚次数、扩容次数、删除次数. 可按时间筛选

应用状态信息:

实例个数、健康状态(可运行实例数/总实例数)、实例重启次数

服务信息:

服务个数、健康状态

用户个人信息:

数据中心、组织、账户信息、配额信息

账户信息:

充值记录、当前余额、消费记录

配额信息:

配额总量、当前余额


### 管理员用户

数据中心信息:

数据中心位置、集群数量、节点网络、CPU、内存、硬盘使用情况

镜像仓库信息:

镜像仓库位置、网络、硬盘使用情况、同步情况

用户、组织审计信息:

