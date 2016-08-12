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
GET /api/v1/organizations/{id}
```

#### 修改某个组织

```bash
POST /api/v1/organizations/{id}

```

#### 查询某个组织下有关联的数据中心
```bash
GET /api/v1/organizations/{id}/datacenters
```

#### 删除某个组织

```bash
DELETE /api/v1/organizations/{id}
```

### 用户
----------------------------------------------------------

#### 新建用户

```bash
POST /api/v1/organizations/{id}/users
```

#### 查询用户: 用户名使用email, 不允许重复

查询用户列表

```bash
GET /api/v1/organizations/{id}/users
```

查询某个用户的信息

```bash
GET /api/v1/organizations/{id}/users/{id}
```

#### 修改某个用户

```bash
POST /api/v1/organizations/{id}/users/{id}
```

#### 删除某个用户

```bash
DELETE /api/v1/organizations/{id}/users/{id}
```

登录

```bash
POST /api/v1/users/{email}/login
```

登出

```bash
POST /api/v1/users/{id}/logout
```

修改密码

```bash
POST /api/v1/user/{id}/password
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
GET /api/v1/datacenters/{id}
```

#### 修改某个数据中心

```bash
POST /api/v1/datacenters/{id}
```

#### 删除某个数据中心

```bash
DELETE /api/v1/datacenters/{id}
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
POST /api/v1/organizations/{id}/dcquotas/
```

#### 查询某个组织下的配额列表(包括所关联的数据中心)

```bash
GET /api/v1/organizations/{id}/dcquotas/
```

#### 查询某个组织下的配额的详细信息

```bash
GET /api/v1/organizations/{id}/dcquotas/{id}
```

#### 修改某个组织下的某个配额的信息

```bash
POST /api/v1/organizations/{id}/dcquotas/{id}
```

#### 删除某个组织下的某个配额

```bash
DELETE /api/v1/organizations/{id}/dcquotas/{id}
```

### 云盘
----------------------------------------------------------

#### 新建云盘

```bash
POST /api/v1/organizations/{id}/rbds/
```

#### 查看云盘

查询列表:

```bash
GET /api/v1/organizations/{id}/rbds/
```

查询某个云盘的细节

```bash
GET /api/v1/organizations/{id}/rbds/{id}
```

#### 修改某个云盘的细节

```bash
POST /api/v1/organizations/{id}/rbds/{id}  //主要应用场景是云盘的扩容
```

#### 删除某个云盘

```bash
DELETE /api/v1/organizations/{id}/rbds/{id}
```

### 发布操作记录
----------------------------------------------------------

#### 新建发布操作

```bash
POST /api/v1/deployments
```

#### 查询过去的发布操作(倒序)

查询列表:

```bash
GET /api/v1/organizations/{id}/deployments
```

查询详细信息

```bash
GET /api/v1/organizations/{id}/deployments/{id}
```


**发布操作表没有删除和修改操作,只能读取和追加.**
