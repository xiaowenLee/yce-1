应用列表
===========

用户点击应用列表时请求后台数据:


之前设计为：请求的方法及URL: GET /api/v1/organizations/{orgId}/users/{uid}/deployments*

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
            "dcId": "bangongwang",
            "podlist": {
                //该数据中心下的应用列列表，json为k8s原生[PodList](https://godoc.org/k8s.io/kubernetes/pkg/api#PodList)
            }
    }]
}
```

podList的json结构：

```json
{
  "kind": "PodList",
  "apiVersion": "v1",
  "metadata": {
    "selfLink": "/api/v1/namespaces/ops/pods",
    "resourceVersion": "16266621"
  },
  "items": [
    {
      "metadata": {
        "name": "nginx-test-1252813378-39wjk",
        "generateName": "nginx-test-1252813378-",
        "namespace": "ops",
        "selfLink": "/api/v1/namespaces/ops/pods/nginx-test-1252813378-39wjk",
        "uid": "217b98e7-64ee-11e6-b957-44a84240716a",
        "resourceVersion": "16258268",
        "creationTimestamp": "2016-08-18T02:47:40Z",
        "labels": {
          "name": "spec-template-metadata-labels",
          "pod-template-hash": "1252813378"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"ops\",\"name\":\"nginx-test-1252813378\",\"uid\":\"2179503e-64ee-11e6-b957-44a84240716a\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"16258237\"}}"
        }
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-jr2e0",
            "secret": {
              "secretName": "default-token-jr2e0"
            }
          }
        ],
        "containers": [
          {
            "name": "nginx-test",
            "image": "nginx:1.7.9",
            "resources": {},
            "volumeMounts": [
              {
                "name": "default-token-jr2e0",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "imagePullPolicy": "IfNotPresent"
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "172.21.1.21",
        "securityContext": {}
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2016-08-18T02:47:41Z"
          }
        ],
        "hostIP": "172.21.1.21",
        "podIP": "10.0.62.3",
        "startTime": "2016-08-18T02:47:40Z",
        "containerStatuses": [
          {
            "name": "nginx-test",
            "state": {
              "running": {
                "startedAt": "2016-08-18T02:47:41Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "nginx:1.7.9",
            "imageID": "docker://sha256:84581e99d807a703c9c03bd1a31cd9621815155ac72a7365fd02311264512656",
            "containerID": "docker://0ac0b1d2e4bc085a1655049e7c251c056bffe05ec26b60523bd34fd590bcc472"
          }
        ]
      }
    }
    ]
}
    //还有一些省略了
```

根据应用列表页面的设计，要显示的信息及相关说明如下：

    |信息：      |  说明：|
    |:------------:|:--------------:|
    |ID          |  数字，为页面显示ID|
    |应用名称    |  data[].podList.items[].metadata.name |
    |标签组合    |  data[].podList.items[].metadata.labels |
    |数据中心    |  data[].dataCenter, 需要为中文 |
    |副本个数    |  data[].podList.items[] 的元素个数 |
    |运行时长    |  data[].podList.items[].metadata.creationTimestamp，需要转化为天、分、时、秒 |

