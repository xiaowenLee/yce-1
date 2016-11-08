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

* [检查重名]()

### 绿色通道

### 错误Error

### 镜像image

### 日志log

### 登录login

### 注销logout

### 导航栏navList

### 开放端口nodePort

### 个人中心personal

### 镜像仓库registry

### 服务及访问点service_endpoint

### 拓扑topoloty

### 用户user











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
