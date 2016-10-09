1, 资源

 ```
   {
       "code": 0,
       "msg": "操作成功",
       "data": {
           "overall": {                 // 总览
               "cpu": {
                   "total": 20,         // 总额
                   "used":  6           // 已用
               }, 
               "mem": {
                   "total": 64,
                   "used": 24
               }
           },
           "datacenters": [{            // 分数据中心显示
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
           }]
       }
   }
 ```
   
2, 应用

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

3, 操作

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