<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

Dashboard操作统计信息
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-06

目录
--------------
###目的
在Dashboard上显示最近的操作情况。


###请求

* 请求方法: GET 
* 请求URL: /api/v1/organizations/{:orgId}/operationstat
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 无


###页面设计 
以柱形图的方式显示最近一周的操作统计。


###程序实现逻辑:

```Sequence
Title: 数据中心列表 
YCE-->>MySQL: 获取最近不同操作统计情况
YCE<<--MySQL: 返回请求结果
```

###响应数据结构: 
JSON
```json
{
    "code": 0,
    "msg": "操作成功", 
    "data": {
        "statistics": {
            "online": [2,2,1,4,5,1,7],          // 发布 
            "scale": [2,3,0,5,0,7,8],           // 扩容
            "rollingupgrade": [1,2,1,4,0,6,7],  // 滚动升级 
            "rollback": [1,0,3,7,5,8,7],        // 回滚
            "delete": [0,0,0,0,0,0,0]           // 删除
        },
        "date": ["2016-10-10", "2016-10-09", "2016-10-08", "2016-10-07", "2016-10-06", "2016-10-05", "2016-10-04", "2016-10-03"]
    } 
}
```

### 备注
以后应添加按时间范围进行筛选,同时不同数据中心间有间隔。
