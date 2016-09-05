NodePort
----------------

NodePort是Kubernetes里与Service紧密相关的概念, 从Kubernetes的文档可知:

Kubernetes master会从标记域里随机选择一个数作为端口, 默认的域范围是30000 ~ 32767。选定之后,每个Node将会代理这个端口到Service里,它可以从Service的spec.ports[*].nodePort域得到。

如果手动从这个域里指定一个NodePort,那么系统将为分配这个数作为端口号,如果冲突将导致分配失败。

NodePort使得开发者可以自行设定负载均衡器,或者配置Kubernetes不完全支持的云环境,甚至直接暴露多个Node的IP。


### 场景说明

现在采用的是一个数据中心一套Kubernetes集群, 每个Kubernetes集群里不能有冲突的NodePort端口。


### 数据库

用户填写了自己的NodePort后,需要去数据库进行验证这个NodePort是否已被占用。检查需要在发布提交前完成。如果不检查,用户需尝试多个NodePort以找到可用的NodePort。

为了简化操作, 可能需要推荐未使用的端口。数据范围在填表时进行检查,默认为30000 ~ 32767, 超出要求重填。

建nodeport表,里面列有:Id, port, dcId, svcId, status, createdAt, modifiedAt, modifiedOp, comment

*将NodePort作为Datacenter的新字段,数据类型为VARCHAR(255),存储的值为{"nodePort": [""]}形式。处理不方便,不用*

*MySQL的数据类型Set最多仅有64个, 不能采用Set*

*Enum最多可达65536个,值来自于创建表时显式指定的值, 不能用Enum*

删除的时候采用更改status字段从INVALID变为VALID而不真正删除

插入的时候采用 Insert On Duplicate Key Update。如果记录存在且为INVALID, 则返回错误。如果记录存在且为VALID, 则更新它的svcName, status, modifiedAt, modifiedOp字段,如果记录不存在,则插入新记录。

表里采用了双主键(port 和 dcId)以唯一确定一条记录


### 方法

```golang

// 查询是否存在该port和dcId组合, 如果存在,返回Nil, 如果不存在,返回err
func (np *NodePort) QueryNodePortByPortAndDcId(port, dcId int32) error 

// 根据NodePort号和所属DcId号查找相应的serviceId, 存在返回ServiceId和Nil, 不存在返回""和err
func (np *NodePort) QueryServiceNameByPortAndDcId(port, dcId int32) (string, error) 

// 插入port, 如果已存在且INVALID,返回err; 如果存在且为VALID, 更新记录的svcName, status, modifiedAt, modifiedOp, 并返回Nil。如果不存在,插入新记录并返回Nil, 插入失败返回err
func (np *NodePort) InsertNodePort(op int32) error 

//插入时,如果记录存在,就更新里面的一些字段,如果不存在则插入新记录。作用在具有唯一索引或主键。
func (np *NodePort) InsertOnDuplicateKeyUpdateStatus(svcName string, status int32, op int32) error

// 更新port对应的信息, 该记录不存在或该记录存在更新成功返回nil, 更新失败返回err
func (np *NodePort) UpdateNodePort(op int32) error

// 删除port对应的信息(修改status为VALID),  该记录不存在或该记录存在删除成功返回nil, 删除失败返回err
func (np *NodePort) DeleteNodePort(op int32) error

```