服务及访问点
--------------

用户点击扩展功能请求后台数据:

请求的方法及URL: GET /api/v1/organizations/{orgId}/users/{userId}/extensions

请求头中包含: Authorization: ${sessionId}

返回值:

* 该组织下数据中心里的服务及访问点列表

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
            "serviceList": {
                //该数据中心下服务列表，json为k8s原生ServiceList
            },
            "endpointsList": {
                //该数据中心下访问点列表, json为k8s原生EndpointsList
            }
    }]
}
```

根据列表页面的设计，要显示的信息及相关说明如下：

|信息：      |  说明：|
|:------------:|:--------------:|
|ID          |  数字，为页面显示ID|
|名称    |      data[].serviceList.serviceList.items[].metadata.name 或 data[].endpointsList.endpointsList.items[].metadata.name|
|类型       |   *如果是service显示为服务,如果是endpoints显示为访问点* 
|数据中心    |  data[].dcName, 需要为中文 |
|创建者 | data[].serviceList.items[].labels["author"] 或 data[].endpointsList.items.labels["author"]| 
|运行时长    |  data[].serviceList.items[].metadata.creationTimestamp，需要转化为天、分、时、秒（一级） |

需要详情相关定义


### 暂时不用
-------------
用户点击服务管理请求后台数据:

请求的方法及URL: GET /api/v1/organizations/{orgId}/users/{userId}/services

请求头中包含: Authorization: ${sessionId}

返回值:

* 该组织下数据中心里的服务列表

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
            "svclist": {
                //该数据中心下服务列表，json为k8s原生ServiceList
            }
    }]
}
```

根据服务列表页面的设计，要显示的信息及相关说明如下：

|信息：      |  说明：|
|:------------:|:--------------:|
|ID          |  数字，为页面显示ID|
|服务名称    |  data[].serviceList.items[].metadata.name |
|Selector    |  data[].serviceList.items[].spec.selector |
|数据中心    |  data[].dataCenter, 需要为中文 |
|地址及端口   |   data[].serviceList.items[].spec.ports[], 可为中文 |
|类型       |   data[].serviceList. 
|运行时长    |  data[].serviceList.items[].metadata.creationTimestamp，需要转化为天、分、时、秒（一级） |

用户点击服务名可以看到详情:

请求的方法及URL: GET /api/v1/organizatinos/{orgId}/users/{userId}/services/{servicesId}

请求头中包含: Authorization: ${sessionId}

返回值:

* 该服务的详细信息

返回json示例:

```json
{
    "code":0,
    "message":[
        "OK"
    ],
    "data": [{
            "dcId": 1,
            "dcName": "bangongwang",
            "service": {
                //该服务详情，json为k8s原生Service
            }
    }]
}
```

|信息：      |  说明：|
|:------------:|:--------------:|
|ID          |  数字，为页面显示ID|
|服务名称    |  data[].serviceList.items[].metadata.name |
|FQDN       | 暂时未知
|端口详情      |  data[].serviceList.items[].spec.ports[], 可为中文 |
|运行时长    |  data[].serviceList.items[].metadata.creationTimestamp，需要转化为
|相应访问点  | data[].endpointList.items[].subsets[]

