<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

查看应用详情
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
用户点击应用名可以查看到应用详情。这里的详情是实例的详情。

###请求

* 请求方法: GET 
* 请求URL: /api/v1/organizations/{orgId}/users/{uid}/deployments/{podId}
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
    无

###页面设计 
|信息：      |  说明：|
|:------------:|:--------------:|
|ID          |  数字，为页面显示ID|
|应用名称    |  data[].podList.items[].metadata.name |
|标签组合    |  data[].podList.items[].metadata.labels |
|数据中心    |  data[].dataCenter, 需要为中文 |
|当前状态    |  data[].podList.items[].status.phase, 需要为中文 |
|运行时长    |  data[].podList.items[].metadata.creationTimestamp，需要转化为天、分、时、秒（一级） |
|所属用户    |  data[].podList.items[].metadata.labels["maintainer"]  |
|所属组织    |  data[].podList.items[].metadata.labels["organzitions"]  |
|节点IP      | data[].podList.items[].status.hostIP  |
|应用IP      | data[].podList.items[].status.podIP  |
|镜像        | data[].podList.items[].status.containerStatuses.image  |
|重启次数    | data[].podList.items[].status.containerStatuses.restartCount  |

###程序实现逻辑:

```Sequence
Title: 查看应用详情
YCE-->>MySQL: 根据orgId请求获取orgName 
YCE<<--MySQL: 返回请求结果
YCE-->>K8s: 在对应的orgName命名空间里找到对应的deployment
YCE<<--K8s: 返回请求结果 
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

```json
{
    "code":0,
    "message":[
        "OK"
    ],
    "data": [{
            "dcId": 1,
            "dcName": "bangongwang",
            "podlist": {
                //该数据中心下的应用列列表，json为k8s原生[PodList](https://godoc.org/k8s.io/kubernetes/pkg/api#PodList)
            }
    }]
}
```

### 备注
应用管理列表的时候已经请求到了deployment的全部信息. 所以不需要通过本接口重复请求了。

### 以下为旧版本, 无效///////////////////////////////////////////////////
应用详情
==============

用户点击应用名时弹窗显示应用详情:

请求的方法及URL: GET /api/v1/organizations/{orgId}/users/{uid}/deployments/{podId}

请求头中包含: Authorization: ${sessionId} *暂时在Session Storage里*

返回值:

* 该组织下数据中心里的应用列表

返回json示例：

```json
{
    "code":0,
    "message":[
        "OK"
    ],
    "data": [{
            "dcId": 1,
            "dcName": "bangongwang",
            "podlist": {
                //该数据中心下的应用列列表，json为k8s原生[PodList](https://godoc.org/k8s.io/kubernetes/pkg/api#PodList)
            }
    }]
}
```

应用详情是在应用列表的基础上，对里面的应用信息进一步筛选, 然后显示在弹窗里

根据应用详情页面的设计，要显示的信息及相关说明如下：

|信息：      |  说明：|
|:------------:|:--------------:|
|ID          |  数字，为页面显示ID|
|应用名称    |  data[].podList.items[].metadata.name |
|标签组合    |  data[].podList.items[].metadata.labels |
|数据中心    |  data[].dataCenter, 需要为中文 |
|当前状态    |  data[].podList.items[].status.phase, 需要为中文 |
|运行时长    |  data[].podList.items[].metadata.creationTimestamp，需要转化为天、分、时、秒（一级） |
|所属用户    |  data[].podList.items[].metadata.labels["maintainer"]  |
|所属组织    |  data[].podList.items[].metadata.labels["organzitions"]  |
|节点IP      | data[].podList.items[].status.hostIP  |
|应用IP      | data[].podList.items[].status.podIP  |
|镜像        | data[].podList.items[].status.containerStatuses.image  |
|重启次数    | data[].podList.items[].status.containerStatuses.restartCount  |
|云盘        |  -  |
