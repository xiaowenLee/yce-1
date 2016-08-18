应用列表
===========

用户点击应用列表时请求后台数据:


*请求的方法及URL: GET /api/v1/organizations/{orgName}/users/{uid}/deployments*

暂时先用GET /api/v1/organizations/{orgName}/deployments

请求头中包含: Authorization: ${sessionId}

其中: uid, orgId, sessionId从LocalStorage里面获取后

返回值:

* 该组织下数据中心里的应用列表

大概数据结构：

```
{
    "code":{},
    "message":[
        ""
    ],
    "data": [{
            "dcId": "",
            "podlist": {
                //该数据中心下的应用列列表，json为k8s原生[PodList](https://godoc.org/k8s.io/kubernetes/pkg/api#PodList)
            }
    }]
}
```


