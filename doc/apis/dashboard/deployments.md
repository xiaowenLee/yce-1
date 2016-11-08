<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

Dashboard应用统计信息
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-06

目录
--------------
###目的
在Dashboard上显示所有数据中心的应用情况。


###请求

* 请求方法: GET 
* 请求URL: /api/v1/organizations/{:orgId}/deploymentstat
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 无


###页面设计 

以拓扑卡片的形式显示该数据中心下的应用情况。不同的数据中心应间隔开。


###程序实现逻辑:

```Sequence
Title: 数据中心列表 
YCE-->>MySQL: 获取该组织对应的可用数据中心列表
YCE<<--MySQL: 返回请求结果
YCE-->>K8S: 获取各个数据中心的Deployment
YCE--<<K8S: 返回请求结果
YCE-->>K8S: 根据Deployment获取最新使用的RS
YCE<<--K8S: 返回请求结果
YCE-->>K8S: 根据RS获取目前运行的所有pod
YCE<<--K8S: 返回请求结果
```

###响应数据结构: 
JSON
```json
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

### 备注
获取该数据中心下的全部deployment,再找到NewRS,再找到pod。
