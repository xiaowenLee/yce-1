<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

绿色通道
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
通过引导式访问,使得用户能够便捷地发布应用及服务

###请求

* 请求方法: 
* 请求URL: 
* 请求头: Authorization:$SessionId, 从LocalStorage读 
* 请求参数: 
 详见应用发布与服务发布 

###页面设计 
1: 发布应用
2: 发布服务
3: 发布状态 

###程序实现逻辑:
无

###响应数据结构: 
返回码为0, 表示操作成功。
其他返回码表示出错。

### 备注
无


### 以下为旧版本, 无效///////////////////////////////////////////////////
绿色通道
=============

### 设计目的:
通过引导式访问,使得用户能够便捷地发布应用及服务

### 设计详情:
页面设计:
1: 发布应用
2: 发布服务
3: 发布状态 

使用模板: form-wizard.html

注: 这里可以选一个有趣的应用, 例如2048等等, 做一个发布会时的demo, 现场观众登录体验。或者其他的应用。


### API及数据结构:
尽量采用原来的API及数据结构:

* 发布应用: POST /api/v1/organizations/:orgId/users/:userId/deployments/new
* 发布服务: POST /api/v1/organizations/:orgId/users/:userId/services/new

具体:

1. 用户点击绿色通道的时候显示发布应用的页面时, 
   操作: 请求 GET /api/v1/organizations/:orgId/users/:userId/deployments/init
   目的: 获取数据中心列表, 所属组织名等
   
2. 用户填写应用名完毕时检查是否重名,
   操作: 请求 POST /api/v1/organizations/:orgId/users/:userId/deployments/check
   目的: 检查用户所填应用名是否已经被使用
   
3. 用户填写完毕,点击下一步时
   操作1: 请求 GET /api/v1/organizations/:orgId/users/:userId/services/init
   目的: 获取数据中心列表等
   
   操作2: 将用户填写的应用信息拼接为JSON(依实际开发习惯而定)
   目的: 为后面的提交做好准备
   
4. 显示服布服务页面, 初始化服务的某些信息
   操作1: 初始化数据中心, 默认选择跟创建应用同样的数据中心
   操作2: 初始化选择器, 应用前面创建应用的label属性里的key name及对应值, 赋予选择器的key和value
   操作3: 推荐nodePort
   
5. 用户填写服务名完毕时检查是否重名
   操作: 请求 POST /api/v1/organizations/:orgId/users/:userId/services/check
   目的: 检查用户所填服务名是否已经被使用

6. 用户填写nodePort完毕时检查是否重复
   操作：请求 POST /api/v1/organizations/:orgId/users/:userId/nodeport/check (后台待开发)
   目的：检查用户所填写nodePort是否已经被使用

7. 用户填写完毕, 点击上一步时
   操作1: 将当前的服务信息保存为JSON
   目的: 为后面的提交做好准备
  
   操作2: 请求GET /api/v1/organizations/:orgId/users/:userId/deployments/init
   目的: 为用户修改应用信息做好准备
   
   操作3: 按用户之前所填的JSON填写相应的文本框控件等(依实际开发习惯而定)
   目的: 同上

8. 用户所有信息填写完毕, 点击提交
   操作1: 请求 POST /api/v1/organizations/:orgId/users/:userId/deployments/new
   目的: 发布应用
   
   操作2: 请求 POST /api/v1/organizations/:orgId/users/:userId/services/new
   目的: 发布服务
   
   操作3: 跳转到第三页面(发布状态页面)

附:

   [初始化应用: deployments/init的文档地址](https://github.com/lth2015/yce/blob/master/doc/apis/deploy/create.md#应用发布准备)
   [检查应用名: deployments/check的文档地址](https://github.com/lth2015/yce/blob/master/doc/apis/deploy/check/README.md)
   [发布应用: deployments/new的文档地址](https://github.com/lth2015/yce/blob/master/doc/apis/deploy/create.md#应用发布请求提交)
   
   [初始化服务: services/init的文档地址](https://github.com/lth2015/yce/blob/master/doc/apis/service_endpoint/create_service.md#服务发布准备)
   [检查服务名: services/check的文档地址](https://github.com/lth2015/yce/blob/master/doc/apis/service_endpoint/check_service%26endpoint.md)
   [发布服务: services/new的文档地址](https://github.com/lth2015/yce/blob/master/doc/apis/service_endpoint/create_service.md#服务发布提交)
