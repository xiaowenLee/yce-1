<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

发布应用
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-22

目录
--------------
###目的
由用户填写应用相关信息,按照deployment形式进行发布

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organization/{orgId}/users/{userId}/deployments/new
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
JSON
```json
{
  "appName": "xxx", // 应用名
  "orgName": "xxx", // 组织名
  "dcIdList": [1,2],   // 数据中心
  "deployment": {
    "kind": "Deployment",               // 可忽略,建议保留,值为固定值
    "apiVersion":"extensions/v1beta1",  // 可忽略,建议保留,值为固定值
    "metadata": {
        "name": "xxx", // 应用名          
        "namespace": "xxx", // 组织名 
        "labels": { // 标签组,类型为map
            "name": "xxx",  //key默认值name, value为应用名 
            "version": "xxx", //key默认值为version, value为镜像被冒号:分割的最后一段。例如img.reg.3g:15000/nginx:1.7.9,此处为1.7.9
            "author": "xxx", //key默认值为author, value为从本地读取的userName
            "xxx": "xxx" //key为用户自填, value也为用户自填
        }
    },
    "spec": {
        "replicas": 3, // 副本数 
        "template": {
            "spec": {
                "volumes": [
                    "name": "xxx", //存储卷名字
                    "hostPath": {
                        "path": "xxx"  //存储卷宿主文件
                    }
                ],
                "containers": [
                    "name": "xxx", // 应用名 
                    "image": "xxx", //镜像
                    "env": [        // 环境变量
                        {
                            "name": "xxx", //环境变量的key
                            "value": "xxx" //环境变量的value
                        } 
                    ]
                    "ports": [             //开放端口
                        {
                            "name": "xxx", //名称 
                            "containerPort": 8080, //端口 
                            "protocol": "xxx" //协议
                        } 
                    ],  
                    "command": [ // 执行命令
                        "xxx"    
                    ],
                    "args": [    // 参数
                        "xxx" 
                    ],
                    "livenessProbe": { //健康检查
                        "httpGet": {
                            "path": "xxx", // URL 
                            "port": 8080  // Port
                        },
                        "periodSeconds": 10, // 每隔
                        "initialDelaySeconds": 30 //生效时间
                    },
                    "volumeMounts": [
                        {
                            "name": "xxx", //存储卷名字 
                            "mountPath": "xxx", //存储卷应用目录
                            "readOnly": true //存储卷是否只读
                        }
                    ],
                    "lifecycle": {
                        "postStart": { //启动准备
                            "exec": {
                                "command": [
                                    "xxx"  //启动准备脚本
                                ] 
                            } 
                        },
                        "preStop": { //优雅停止
                            "exec": {
                                "command": [
                                    "xxx"  //优雅停止脚本
                                ] 
                            } 
                        }
                    }
                ] 
            } 
        }
    }
  }
}
```

###页面设计 

假设最终生成的json名为data

基本配置:


| 项目 | 变量 | 说明 |
|:------:|:--:|:----:|
|应用名|data.appName, data.deployment.metadata.name, data.deployment.metadata.labels["name"], data.deployment.spec.template.spec.containers[0].name|上述四处值一致|
|组织名|data.orgName, data.deployment.metadata.namespace| 上述值均相同, 从本地读取orgName|
|数据中心|data.dcIdList|可多选|
|镜像|data.deployment.spec.template.spec.containers[0].image||
|副本数|data.deployment.spec.replicas|整型|
|环境变量NAME|data.deployment.spec.template.spec.containers[0].env[$index].name||
|环境变量NAME|data.deployment.spec.template.spec.containers[0].env[$index].name||
|环境变量VALUE|data.deployment.spec.template.spec.containers[0].env[$index].value||
|开放端口名称|data.deployment.spec.template.spec.containers[0].ports[$index].name||
|开放端口端口|data.deployment.spec.template.spec.containers[0].ports[$index].port|整型|
|开放端口协议|data.deployment.spec.template.spec.containers[0].ports[$index].protocol||



高级配置:

|项目|变量|说明|
|:-:|:--:|:-:|
|标签组Key|name|默认|
|标签组Value|data.deployment.metadata.labels["name"]|同应用名|
|标签组Key|version|默认|
|标签组Value|data.deployment.metadata.labels["version"]| data.deployment.spec.template.spec.containers[0].image里由冒号分割的最后一段|
|标签组Key|author|默认|
|标签组Value|data.deployment.metadata.labels["author"]|从本地读取userName|
|标签组Key||用户填写Key|
|标签组Value|data.deployment.metadata.labels["xxx"]|用户填写Value
|健康检查URL|data.deployment.spec.template.spec.containers[0].livenessProbe.httpGet.path||
|健康检查Port|data.deployment.spec.template.spec.containers[0].livenessProbe.httpGet.port|整型|
|健康检查每隔|data.deployment.spec.template.spec.containers[0].livenessProbe.periodSeconds|整型|
|健康检查生效时间|data.deployment.spec.template.spec.containers[0].livenessProbe.initialDelaySeconds|整型|
|存储卷名称|data.deployment.spec.template.spec.volumes[$index].name, data.deployment.spec.template.spec.containers[0].volumeMounts[$index].name|二者值一致,index一致|
|存储卷应用目录|data.deployment.spec.template.spec.containers[0].volumeMounts[$index].mountPath|index同上|
|存储卷宿主文件|data.deployment.spec.template.spec.volumes[$index].hostPath.path|index同上|
|存储卷只读|data.deployment.spec.template.spec.containers[0].volumeMounts[$index].readOnly|index同上|
|启动准备脚本|data.deployment.spec.template.spec.containers[0].lifecycle.postStart.exec.command[0]||
|优雅停止脚本|data.deployment.spec.template.spec.containers[0].lifecycle.preStop.exec.command[0]||
|执行命令脚本|data.deployment.spec.template.spec.containers[0].command[0]||
|参数脚本|data.deployment.spec.template.spec.containers[0].args[0]||


###程序实现逻辑:

```Sequence
Title: 发布应用
YCE-->>K8s: 发布Deployment
YCE--<<K8s: 返回发布结果
YCE-->>MySQL: 插入发布记录
YCE<<--MySQL: 返回插入结果
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

### 备注
应该将创建Deployment和数据库插入做成事务,保持一致。
部分项目的中文名称可以继续优化




### 以下为旧版本, 无效///////////////////////////////////////////////////
应用发布
============
  JSON
```json
  {
    "dataCenters": [
        {
            "dcId": 1
        },
        {
            "dcId": 3
        },
        {
            "dcId": 5
        }
    ]
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

### 应用发布准备

用户点击应用发布(左侧菜单)时请求后台数据:

请求的URL:
```bash
GET /api/v1/organizations/{orgId}/users/{userId}/deployments/init
```

请求头中包含: Authorization: ${sessionId}

其中: userId, orgId, sessionId在登录成功后从后台返回给浏览器, 前端存储在LocalStorage(目前为SessionStorage)里面

返回值:

* 组织名称

* 该组织下的数据中心列表

* 该组织的配额和预算

* 标准配额列表

获取成功的大概数据格式如下:

```json
{
    "code": 0,
    "message": "....",
    "data": {
        "orgId":  "1",
        "orgName": "Ops",
        "dataCenters": [
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
}
```


### 镜像搜索辅助

在应用发布页面中,点击镜像输入框后,弹出选择镜像的窗口

弹出框上面有搜索框,支持输入辅助,就是可以根据用户的输入筛选镜像列表

在点击输入框后,前台要向后台发送请求:

请求的URL: GET /api/v1/images/

请求头包含: Authorization: ${sessionId}

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

请求的URL: POST /api/v1/organization/{orgId}/users/{userId}/deployments/new

请求头包含: Authorization: ${sessionId}

POST数据格式(data里面的是实例,用于讲解跟页面的输入框的关系,更严谨的定义看后面)

```json
{
    "dataCenters": [
        {
            "dcId": 1
        },
        {
            "dcId": 3
        },
        {
            "dcId": 5
        }
    ]
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