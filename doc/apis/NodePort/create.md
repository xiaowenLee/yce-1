
### NodePort创建

#### NodePort表的创建
当前数据中心的默认占用nodePort, 仅有yce暴露占用的nodePort。 
用户添加数据中心的同时, 去创建nodePort表, 默认是30000~32767, 并除去系统本身所占用的nodePort(yce-svc: 31080, 以后会指定为32080)

#### NodePort状态的变化
用户发布服务, 按填写的nodePort, 将对应的数据库记录的status改为INVALID, 表示已被占用.
删除服务会将对应的nodePort的status字段更改为VALID。


### NodePort检查
用户填写了nodePort, 去数据库里检查该nodePort是否为INVALID, 如果为INVALID表示已被占用, 如果为VALID表示可以被使用。
检查时检查的是dcId与nodePort的唯一性。如果是多个数据中心, 提交的应该是dcIdList与nodePort, 后端依次检查, 有一组被占用即不可用

### NodePort管理
管理员可以看到nodePort的列表, 仅显示已被占用的, 显示数据中心的内容


### NodePort推荐
//按首次适配方法, 按照选定的数据中心, 从前向后选择第一个status为VALID的nodePort并返回
按首次适配的方法, 选定第一个两个数据中心里均为VALID的nodePort并返回
推荐的时候改为VALID
假设没有两个数据中心的交集, 应该按选择的数据中心进行推荐(该数据中心里第一个为VALID的nodePort)。 


程序设计逻辑:

1. admin账户用例: 创建数据中心, 初始化nodePort表, 每个新插入表中的nodePort状态为VALID, 表示可用
2. admin账户用例: 查看nodePort被占用列表, 并提供详情展示comments等, 选择nodePort状态为INVALID的列出
3. admin账户用例: 删除数据中心, 创建时插入的nodePort表, 将每个nodePort的状态置为INVALID, 表示不可用
4. 普通账户用例: 填写服务nodePort, 检查nodePort是否可用。 可用为VALID
5. 普通账户用例: 发布服务, 初始化发布服务页面, 附带返回推荐nodePort。 
6. 普通账户用例: 发布服务, 更改nodePort状态为INVALID; 删除服务, 更改nodePort状态为VALID。

