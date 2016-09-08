应用实例日志
--------------
用户点击删除应用时提示拼接json, 点击确认删除时发送请求

请求的方法及URL: GET /api/v1/organizations/{ordId}/pods/{podName}/logs

请求头中包含: Authorization: ${sessionId}

发送Json格式:

```json
  {
    "userId": "1", 
    "dcId": [1], 
    "logOption": {
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


### 查看日志步骤

logOption ---> tail, timeStamp

dcId --> apiServer --> k8sClient --> namespace & podName ---> getLog

getLog --> restclient.Request

restclient.Request -->req.Stream() --> io.ReadCloser

byte, err := ioutil.ReadAll(io.ReadCloser)

writeBack string(byte)

返回json:

{
    "code":
    "message":
    "data": "logs" 
}




//DEPRECATED
io.Copy(io.Writer, io.ReadCloser)

bufio.Writer . out ---> io.Writer

iris.StreamWrite(cd func(writer *bufio.Writer))

先获取该数据中心里该命名空间下的该名称应用

dcId --> orgName --> deploymentName

然后删除两步:

先获取该deployment对应的所有replicaSet, 并依次删除

再将该deployment直接删除。

不等待,不处理删除失败恢复,仅返回操作结果