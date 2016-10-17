绿色通道
=============

### 设计目的:
通过引导式访问,使得用户能够便捷地发布应用及服务

### 设计详情:
页面设计:
1: 发布应用
2: 发布服务
3: Dashboard或拓扑图

使用模板: form-wizard.html

注: 这里可以选一个有趣的应用, 例如2048等等, 做一个发布会时的demo, 现场观众登录体验。或者其他的应用。

页面选择简单的部署设置: 选择镜像部署, 然后绑定实例发布服务。

### API及数据结构:
尽量采用原来的API及数据结构:

* 发布应用: POST /api/v1/deployment/new
* 发布服务: POST /api/v1/service/new  POST /api/v1/endpoint/new
* Dashboard: GET /api/v1/resourcestat GET /api/v1/deploymentstat GET /api/v1/operationstat
* 拓扑图: GET /api/v1/organizations/:orgId/topology 

前后端都需要对输入以及输出的东西负责和控制, 划分边界。比如输入合法性检查等


