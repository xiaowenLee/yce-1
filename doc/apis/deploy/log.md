<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

查看应用日志
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
由用户查看应用某实例的日志

###请求

* 请求方法: POST 
* 请求URL: /api/v1/organizations/{ordId}/pods/{podName}/logs
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
  JSON
```json
 {
    "userId": "1",     // 现有代码有的是string, 有的是int32
    "dcIdList": [1], 
    "logOption": {     //可选
      "Container": "",//暂时不做
      "Follow": false, //暂时不做, 页面开关,默认为关闭
      "Previous": false,//暂时不做
      "SinceSeconds": 0,//暂时不做
      "SinceTime": nil, //暂时不做
      "Timestamps": true, //时间戳,默认打开
      "TailLines": 100 //用户设定
      "LimitBytes": 0, //暂时不做
    } 
 }
```

###页面设计 
无

###程序实现逻辑:

```Sequence
Title: 发布应用
YCE-->>K8s: 根据orgName获取deployment,再找到某个pod的日志
YCE<<--K8s: 返回查找结果
```

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

```json
{
    "code": 0,
    "message":"xx",
    "data": "logs" 
}
```

### 备注
logOption ---> tail, timeStamp

dcId --> apiServer --> k8sClient --> namespace & podName ---> getLog

getLog --> restclient.Request

restclient.Request -->req.Stream() --> io.ReadCloser

byte, err := ioutil.ReadAll(io.ReadCloser)

writeBack string(byte)

### 以下为旧版本, 无效///////////////////////////////////////////////////

应用实例日志
--------------
用户点击删除应用时提示拼接json, 点击确认删除时发送请求

请求的方法及URL: GET /api/v1/organizations/{ordId}/pods/{podName}/logs

请求头中包含: Authorization: ${sessionId}

发送Json格式:

```json
  {
    "userId": "1",     // 现有代码有的是string, 有的是int32
    "dcIdList": [1], 
    "logOption": {     //可选
      "Container": "",//暂时不做
      "Follow": false, //暂时不做, 页面开关,默认为关闭
      "Previous": false,//暂时不做
      "SinceSeconds": 0,//暂时不做
      "SinceTime": nil, //暂时不做
      "Timestamps": true, //时间戳,默认打开
      "TailLines": 100 //用户设定
      "LimitBytes": 0, //暂时不做
    } 
  }
    
```

返回值:

* 操作结果 

返回json:

{
    "code":
    "message":
    "data": "logs" 
}



### 查看日志步骤
logOption ---> tail, timeStamp

dcId --> apiServer --> k8sClient --> namespace & podName ---> getLog

getLog --> restclient.Request

restclient.Request -->req.Stream() --> io.ReadCloser

byte, err := ioutil.ReadAll(io.ReadCloser)

writeBack string(byte)
