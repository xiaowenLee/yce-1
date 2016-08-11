应用发布页面数据交互说明
============

### 点击应用发布(左侧菜单)时请求后台数据:

请求的URL: GET /api/v1/organizations/{orgId}/users/{uid}/deployments/new

请求头中包含: Authorization: ${x-auth-token}

其中: uid, orgId, x-auth-token在登录成功后从后台返回给浏览器, 前端存储在LocalStorage里面

返回值:

* 组织名称

* 该组织下的数据中心列表

* 该组织的配额和预算

* 标准配额列表

大概的数据格式如下:

```json
{
    "orgId":  "1",
    "orgName": "Ops",
    "dataCenter": [
        {
            "dcId": "1",
            "name": "世纪互联",
            "budget": 10000000,
            "balance": 10000000
        },
        {
            "dcId": "2",
            "name": "电信机房",
            "budget": 10000000,
            "balance": 10000000
        },
        {
            "dcId": "3",
            "name": "电子城机房",
            "budget": 10000000,
            "balance": 10000000
        }
    ],
    "dcQuotas": {
        "dcId": "1"
        "PodMax": 1 
        // 第一版用不到...
    }
}
```

### 在应用发布页面中,点击镜像输入框后,弹出选择镜像的窗口

弹出框上面有搜索框,支持输入辅助,就是可以根据用户的输入筛选镜像列表

在点击输入框后,前台要向后台发送请求:

请求的URL: GET /api/v1/images/

请求头包含: Authorization: ${x-auth-token}

返回值:

```json
[
    {
        "imageName": "",
        "imageTag": "",
        // 其他一些可有可无的信息,第一版不需考虑...
    },
    {
    }
]
```

输入辅助是根据imageName来筛选的



### 应用发布请求提交

请求的URL: POST /api/v1/depolyments

请求头包含: Authorization: ${x-auth-token}

POST数据格式:

```json
{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "lifecycle": {  
              "preStop": {  //  启动准备输入框
                "exec": {
                  "command": [
                    "echopreStop"
                  ]
                }
              },
              "postStart": { // 优雅的停止输入框
                "exec": {
                  "command": [
                    "echopostStart"
                  ]
                }
              }
            },
            "readinessProbe": {  // 可读性检查
              "failureThreshold": 0,
              "successThreshold": 0,
              "periodSeconds": 2,  // 每隔多长时间探测
              "timeoutSeconds": 0,
              "initialDelaySeconds": 3, // 启动等待时间
              "httpGet": {
                "httpHeaders": null,
                "scheme": "",
                "host": "",
                "port": 11001,  // 端口
                "path": "http://api/v1/readiness" // 路径
              }
            },
            "name": "nginx-test", // 名称,跟应用名称一样
            "image": "nginx:1.7.9",  // 镜像名称
            "command": [  // 启动命令输入框的
              "echo"
            ],
            "args": [ // 参数输入框的
              "abc"
            ],
            "ports": null,  // 端口
            "env": [  // 环境变量列表
              {
                "value": "good",
                "name": "magic"
              },
              {
                "value": "mushroom",
                "name": "sheep"
              }
            ],
            "resources": {
              "requests": null
            },
            "livenessProbe": {  // 健康检查
              "failureThreshold": 0,
              "successThreshold": 0,
              "periodSeconds": 2, // 每隔多长时间
              "timeoutSeconds": 0,
              "initialDelaySeconds": 3,  // 启动等待时间
              "httpGet": {
                "httpHeaders": null,
                "scheme": "",
                "host": "",
                "port": 11000,  // 端口输入框
                "path": "http://api/v1/healthz" // 请求路径的输入框
              }
            }
          }
        ]
      },
      "metadata": { // 这部分看跟下面的metadata同样,见下面metadata
        "labels": {
          "maintainer": "liyao",
          "appname": "nginx-test"
        },
        "name": "nginx-test"
      }
    },
    "replicas": 3 // 副本个数
  },
  "metadata": {
    "labels": { // 标签,支持多个
      "maintainer": "liyao",
      "appname": "nginx-test"
    },
    "namespace": "default", // 组织名称
    "name": "nginx-test"  // 应用名称那个输入框
  },
  "kind": "Deployment", // 这是写死的
  "apiVersion": "extensions/v1beta1" // 这个默认写死的
}
```


