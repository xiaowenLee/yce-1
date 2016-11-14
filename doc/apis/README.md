<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

容器云对外API定义(第一版)
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-08

目录
--------------
### [控制台Dashboard](https://github.com/maxwell92/yce/blob/master/doc/apis/dashboard/dashboard.md)

* [资源统计](https://github.com/maxwell92/yce/blob/master/doc/apis/dashboard/resources.md)
* [应用统计](https://github.com/maxwell92/yce/blob/master/doc/apis/dashboard/deployments.md)
* [操作统计](https://github.com/maxwell92/yce/blob/master/doc/apis/dashboard/operations.md)

### [数据中心Datacenter](https://github.com/maxwell92/yce/tree/master/doc/apis/datacenter)

* [检查重名](https://github.com/maxwell92/yce/blob/master/doc/apis/datacenter/check.md) 
* [添加组织](https://github.com/maxwell92/yce/blob/master/doc/apis/datacenter/create.md)
* [初始化创建](https://github.com/maxwell92/yce/blob/master/doc/apis/datacenter/init.md)
* [更新组织](https://github.com/maxwell92/yce/blob/master/doc/apis/datacenter/update.md)
* [删除组织](https://github.com/maxwell92/yce/blob/master/doc/apis/datacenter/delete.md)
* [管理组织列表](https://github.com/maxwell92/yce/blob/master/doc/apis/datacenter/list.md)

### [应用Deployment](https://github.com/maxwell92/yce/tree/master/doc/apis/deploy)

* [检查重名](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/check.md)
* [发布应用](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/create.md)
* [删除应用](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/delete.md)
* [应用详情](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/describe.md)
* [发布历史](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/history.md)
* [镜像搜索辅助](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/image.md)
* [初始化创建](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/init.md)
* [应用日志](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/log.md)
* [应用回滚](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/rollback.md)
* [滚动升级](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/rollingupdate.md)
* [扩容](https://github.com/maxwell92/yce/blob/master/doc/apis/deploy/scale.md)

### [绿色通道](https://github.com/maxwell92/yce/tree/master/doc/apis/EasyAccess)

* [绿色通道](https://github.com/maxwell92/yce/blob/master/doc/apis/EasyAccess/easyaccess.md)

### [错误Error](https://github.com/maxwell92/yce/tree/master/doc/apis/error)

* [错误码](https://github.com/maxwell92/yce/blob/master/doc/apis/error/error.md)

### [镜像image](https://github.com/maxwell92/yce/tree/master/doc/apis/image)

* [镜像](https://github.com/maxwell92/yce/blob/master/doc/apis/image/list.md)

### [日志log](https://github.com/maxwell92/yce/tree/master/doc/apis/log)

* [日志处理](https://github.com/maxwell92/yce/blob/master/doc/apis/log/log.md)

### [登录login](https://github.com/maxwell92/yce/tree/master/doc/apis/login)

* [登录](https://github.com/maxwell92/yce/blob/master/doc/apis/login/README.md)

### [注销logout](https://github.com/maxwell92/yce/tree/master/doc/apis/logout)

* [注销](https://github.com/maxwell92/yce/blob/master/doc/apis/logout/README.md)

### [组织namespace](https://github.com/maxwell92/yce/tree/master/doc/apis/namespace)

* [检查重名](https://github.com/maxwell92/yce/blob/master/doc/apis/namespace/check.md)
* [创建组织](https://github.com/maxwell92/yce/blob/master/doc/apis/namespace/create.md)
* [删除组织](https://github.com/maxwell92/yce/blob/master/doc/apis/namespace/delete.md)
* [创建初始化](https://github.com/maxwell92/yce/blob/master/doc/apis/namespace/init.md)
* [组织列表](https://github.com/maxwell92/yce/blob/master/doc/apis/namespace/list.md)
* [更新组织](https://github.com/maxwell92/yce/blob/master/doc/apis/namespace/update.md)

### [导航栏navList](https://github.com/maxwell92/yce/tree/master/doc/apis/navList)

* [导航栏](https://github.com/maxwell92/yce/blob/master/doc/apis/navList/README.md)

### [开放端口nodePort](https://github.com/maxwell92/yce/tree/master/doc/apis/NodePort)

* [检查重复](https://github.com/maxwell92/yce/blob/master/doc/apis/NodePort/check.md)
* [添加端口](https://github.com/maxwell92/yce/blob/master/doc/apis/NodePort/create.md)
* [端口列表](https://github.com/maxwell92/yce/blob/master/doc/apis/NodePort/list.md)

### [个人中心personal](https://github.com/maxwell92/yce/tree/master/doc/apis/personal)

### [镜像仓库registry](https://github.com/maxwell92/yce/tree/master/doc/apis/registry)

* [镜像仓库](https://github.com/maxwell92/yce/blob/master/doc/apis/registry/README.md)

### [服务及访问点service_endpoint](https://github.com/maxwell92/yce/tree/master/doc/apis/service_endpoint)

* [检查重名](https://github.com/maxwell92/yce/blob/master/doc/apis/service_endpoint/check_service%26endpoint.md)
* [创建访问点](https://github.com/maxwell92/yce/blob/master/doc/apis/service_endpoint/create_endpoint.md)
* [创建服务](https://github.com/maxwell92/yce/blob/master/doc/apis/service_endpoint/create_service.md)
* [删除服务](https://github.com/maxwell92/yce/blob/master/doc/apis/service_endpoint/delete_service.md)
* [删除访问点](https://github.com/maxwell92/yce/blob/master/doc/apis/service_endpoint/delete_endpoints.md)
* [服务列表](https://github.com/maxwell92/yce/blob/master/doc/apis/service_endpoint/extensions.md)
* [初始化访问点创建](https://github.com/maxwell92/yce/blob/master/doc/apis/service_endpoint/init_endpoint.md)
* [初始化服务创建](https://github.com/maxwell92/yce/blob/master/doc/apis/service_endpoint/init_service.md)

### [拓扑topoloty](https://github.com/maxwell92/yce/tree/master/doc/apis/topology)

* [应用拓扑](https://github.com/maxwell92/yce/blob/master/doc/apis/topology/README.md)

### [用户user](https://github.com/maxwell92/yce/tree/master/doc/apis/user)

* [创建用户](https://github.com/maxwell92/yce/blob/master/doc/apis/user/create.md)
* [检查重名](https://github.com/maxwell92/yce/blob/master/doc/apis/user/check.md)
* [删除用户](https://github.com/maxwell92/yce/blob/master/doc/apis/user/delete.md)
* [初始化创建](https://github.com/maxwell92/yce/blob/master/doc/apis/user/init.md)
* [用户列表](https://github.com/maxwell92/yce/blob/master/doc/apis/user/list.md)
* [更新用户](https://github.com/maxwell92/yce/blob/master/doc/apis/user/update.md)

### [模板template](https://github.com/maxwell92/yce/tree/master/doc/apis/template)

* [创建模板](https://github.com/maxwell92/yce/blob/master/doc/apis/template/create.md)
* [检查重名](https://github.com/maxwell92/yce/blob/master/doc/apis/template/check.md)
* [模板列表](https://github.com/maxwell92/yce/blob/master/doc/apis/template/list.md)
* [更新模板](https://github.com/maxwell92/yce/blob/master/doc/apis/template/update.md)
* [删除模板](https://github.com/maxwell92/yce/blob/master/doc/apis/template/delete.md)

### 以下为旧版本, 无效///////////////////////////////////////////////////

容器云对外API定义(第一版)
==========================================================

### 组织
----------------------------------------------------------

#### 新建组织

```bash
POST /api/v1/organizations/  
```

#### 查询组织

查询列表:

```bash
GET /api/v1/organizations/
```

查询具体某个组织的细节
```bash
GET /api/v1/organizations/{orgId}
```

#### 修改某个组织

```bash
POST /api/v1/organizations/{orgId}

```

#### 查询某个组织下有关联的数据中心
```bash
GET /api/v1/organizations/{orgId}/datacenters
```

#### 删除某个组织

```bash
DELETE /api/v1/organizations/{orgId}
```

### 用户
----------------------------------------------------------

#### 新建用户

```bash
POST /api/v1/organizations/{orgId}/users
```

#### 查询用户: 用户名使用email, 不允许重复

查询用户列表

```bash
GET /api/v1/organizations/{orgId}/users
```

查询某个用户的信息

```bash
GET /api/v1/organizations/{orgId}/users/{userId}
```

#### 修改某个用户

```bash
POST /api/v1/organizations/{orgId}/users/{userId}
```

#### 删除某个用户

```bash
DELETE /api/v1/organizations/{orgId}/users/{userId}
```

登录

```bash
POST /api/v1/users/login
```

登出

```bash
POST /api/v1/users/{userId}/logout
```

修改密码

```bash
POST /api/v1/user/{userId}/password
```

### 数据中心表
----------------------------------------------------------

#### 新建数据中心

```bash
POST /api/v1/datacenters
```

#### 查询数据中心

查询数据中心列表

```bash
GET /api/v1/datacenters/
```

查询某个数据中心细节

```bash
GET /api/v1/datacenters/{dcId}
```

#### 修改某个数据中心

```bash
POST /api/v1/datacenters/{dcId}
```

#### 删除某个数据中心

```bash
DELETE /api/v1/datacenters/{dcId}
```


### 配额标准
----------------------------------------------------------

#### 新建配额标准

```bash
POST /api/v1/quotas/
```

#### 查询现有配额标准

查询现有配额列表

```bash
GET /api/v1/quotas/
```

查询现有某个配额的细节

```bash
GET /api/v1/quotas/{id}
```

#### 修改某个配额标准

```bash
POST /api/v1/quotas/{id}
```

#### 删除某个配额标准

```bash
DELETE /api/v1/quotas/{id}
```

### 数据中心配额: 某个组织在某个数据中心有多少配额
----------------------------------------------------------

#### 新建数据中心配额

```bash
POST /api/v1/organizations/{orgId}/dcquotas/
```

#### 查询某个组织下的配额列表(包括所关联的数据中心)

```bash
GET /api/v1/organizations/{orgId}/dcquotas/
```

#### 查询某个组织下的配额的详细信息

```bash
GET /api/v1/organizations/{orgId}/dcquotas/{id}
```

#### 修改某个组织下的某个配额的信息

```bash
POST /api/v1/organizations/{orgId}/dcquotas/{id}
```

#### 删除某个组织下的某个配额

```bash
DELETE /api/v1/organizations/{orgId}/dcquotas/{id}
```

### 云盘
----------------------------------------------------------

#### 新建云盘

```bash
POST /api/v1/organizations/{orgId}/rbds/
```

#### 查看云盘

查询列表:

```bash
GET /api/v1/organizations/{orgId}/rbds/
```

查询某个云盘的细节

```bash
GET /api/v1/organizations/{orgId}/rbds/{id}
```

#### 修改某个云盘的细节

```bash
POST /api/v1/organizations/{orgId}/rbds/{id}  //主要应用场景是云盘的扩容
```

#### 删除某个云盘

```bash
DELETE /api/v1/organizations/{orgId}/rbds/{id}
```

### 发布操作记录
----------------------------------------------------------

#### 新建发布操作

```bash
POST /api/v1/organizations/{orgId}/users/{userId}/deployments
```

#### 查询过去的发布操作(倒序)

查询列表:

```bash
GET /api/v1/organizations/{id}/users/{userId}/deployments
```

查询详细信息

```bash
GET /api/v1/organizations/{id}/users/{userId}/deployments
```


**发布操作表没有删除和修改操作,只能读取和追加.**


### 镜像
----------------------------------------------------------

#### 查询仓库里的镜像

查询列表：

```bash
GET /api/v1/registry/images
```
*暂时不按组织做租户*
