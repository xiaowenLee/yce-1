应用发布相关定义
============

Controller定义在src/backend/controller/yce/deploy里面，例如：

应用管理（列表）：ListDeployController的定义在src/backend/controller/yce/deploy/list.go里

Model定义在src/backend/model/yce/deploy里面，例如：

应用管理（列表）：相关定义在src/backend/model/yce/deploy/list.go里

使用说明在当前目录，例如：

应用管理（列表）：使用说明在List.md里

下同：

应用管理(列表)：ListDeployController

应用详情：DescribeDeployController

应用发布：CreateDeployController

应用滚升：RollupDeployController

应用回滚：RolldownDeployController

应用撤销：DrainDeployController

应用扩容：ScaleupDeployController 

应用缩容：ScaledownDeployController 

发布暂停：PauseDeployController

发布恢复：ResumeDeployController

发布历史：HistoryDeployController 


应用发布页面数据交互说明
============

### 点击应用发布(左侧菜单)时请求后台数据:

请求的URL: GET /api/v1/organizations/{orgId}/users/{uid}/deployments/init

请求头中包含: Authorization: ${sessionId}

其中: userId, orgId, sessionId在登录成功后从后台返回给浏览器, 前端存储在LocalStorage里面

返回值:

* 组织名称

* 该组织下的数据中心列表

* 该组织的配额和预算

* 标准配额列表

大概的数据格式如下:

```json
{
  "quotas": [
    {
      "comment": "",
      "modifiedOp": 3,
      "modifiedAt": "2016-08-15T20:58:16+08:00",
      "id": 1,
      "name": "2C4G50G",
      "cpu": 2,
      "mem": 4,
      "rbd": 50,
      "price": "1000",
      "status": 1,
      "createdAt": "2016-08-15T16:32:32Z"
    },
    {
      "comment": "",
      "modifiedOp": 1,
      "modifiedAt": "2016-08-15T16:32:32Z",
      "id": 2,
      "name": "4C8G100G",
      "cpu": 4,
      "mem": 8,
      "rbd": 100,
      "price": "1800",
      "status": 1,
      "createdAt": "2016-08-15T16:32:32Z"
    },
    {
      "comment": "",
      "modifiedOp": 1,
      "modifiedAt": "2016-08-15T16:32:32Z",
      "id": 3,
      "name": "4C16G200G",
      "cpu": 4,
      "mem": 16,
      "rbd": 200,
      "price": "2860",
      "status": 1,
      "createdAt": "2016-08-15T16:32:32Z"
    }
  ],
  "dataCenters": [
    {
      "comment": "",
      "modifiedOp": 1,
      "id": 1,
      "name": "办公网",
      "host": "172.21.1.11",
      "port": 8080,
      "secret": "",
      "status": 1,
      "createdAt": "2016-08-15T16:27:30Z",
      "modifiedAt": "2016-08-15T16:27:30Z"
    },
    {
      "comment": "",
      "modifiedOp": 2,
      "id": 2,
      "name": "电信机房",
      "host": "10.149.149.3",
      "port": 8080,
      "secret": "",
      "status": 0,
      "createdAt": "2016-08-15T17:46:42+08:00",
      "modifiedAt": "2016-08-15T20:36:45+08:00"
    }
  ],
  "orgName": "ops",
  "orgId": "1"
}
```


### 在应用发布页面中,点击镜像输入框后,弹出选择镜像的窗口

弹出框上面有搜索框,支持输入辅助,就是可以根据用户的输入筛选镜像列表

在点击输入框后,前台要向后台发送请求:

请求的URL: GET /api/v1/images/

请求头包含: Authorization: ${x-auth-token}

返回值:

```json
{
    "code": 0,
    "message": "....",
    "data": [
        {
            "imageName": "",
            "imageTag": "",
            // 其他一些可有可无的信息,第一版不需考虑...
        },
        {
            //...
        }
    ]
}
```

输入辅助是根据imageName来筛选的



### 应用发布请求提交

请求的URL: POST /api/v1/organization/{orgId}/users/{userId}/deployments

请求头包含: Authorization: ${sessionId}

POST数据格式(data里面的是实例,用于讲解跟页面的输入框的关系,更严谨的定义看后面)

```json
{
    "dcIdList": [1, 2],
    "appName": "nginx-elb",
    "deployment": {
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
}
```

一个标准的Deployment对象定义
```json
{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "lifecycle": {
              "preStop": {
                "exec": {
                  "command": [
                    "echopreStop"
                  ]
                }
              },
              "postStart": {
                "exec": {
                  "command": [
                    "echopostStart"
                  ]
                }
              }
            },
            "readinessProbe": {
              "failureThreshold": 0,
              "successThreshold": 0,
              "periodSeconds": 2,
              "timeoutSeconds": 0,
              "initialDelaySeconds": 3,
              "httpGet": {
                "httpHeaders": null,
                "scheme": "",
                "host": "",
                "port": 11001,
                "path": "http://api/v1/readiness"
              }
            },
            "name": "nginx-test",
            "image": "nginx:1.7.9",
            "command": [
              "echo"
            ],
            "args": [
              "abc"
            ],
            "ports": null,
            "env": [
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
            "livenessProbe": {
              "failureThreshold": 0,
              "successThreshold": 0,
              "periodSeconds": 2,
              "timeoutSeconds": 0,
              "initialDelaySeconds": 3,
              "httpGet": {
                "httpHeaders": null,
                "scheme": "",
                "host": "",
                "port": 11000,
                "path": "http://api/v1/healthz"
              }
            }
          }
        ]
      },
      "metadata": {
        "labels": {
          "maintainer": "liyao",
          "appname": "nginx-test"
        },
        "name": "nginx-test"
      }
    },
    "replicas": 3
  },
  "metadata": {
    "labels": {
      "maintainer": "liyao",
      "appname": "nginx-test"
    },
    "namespace": "default",
    "name": "nginx-test"
  },
  "kind": "Deployment",
  "apiVersion": "extensions/v1beta1"
}
```
