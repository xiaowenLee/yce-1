package error

type Error struct {
	LogMsg string
	ErrMsg string
}

const (
	EOK    int32 = 0

	EMYSQL int32 = 1000
	EMYSQL_QUERY int32 = 1001
	EMYSQL_INSERT int32 = 1002

	EREDIS int32 = 1100
	EREDIS_GET int32 = 1101

	EKUBE  int32 = 1200
	EKUBE_CLIENT int32 = 1201
	EKUBE_CREATE_DEPLOYMENT = 1202
	EKUBE_LIST_PODS = 1203

	EIRIS  int32 = 1300

	EYCE   int32 = 1400
	EYCE_LOGIN int32 = 1401
	EYCE_SESSION int32 = 1402
	EYCE_SESSION_DEL int32 = 1403
	EYCE_ORG_EXIST int32 = 1404


	EREGISTRY int32 = 1500
	EREGISTRY_GET int32 = 1501

	EJSON int32 = 1600

)

var Errors = map[int32]*Error{

	EOK: &Error{
		LogMsg: "OK",
		ErrMsg: "操作成功",
	},

	// 1000~1099 MySQL错误
	EMYSQL: &Error{
		LogMsg: "MySQL Error",
		ErrMsg: "MySQL数据库错误",
	},
	EMYSQL_QUERY: &Error{
		LogMsg: "MySQL Query Error",
		ErrMsg: "MySQL查询出错",
	},
	EMYSQL_INSERT: &Error{
		LogMsg: "MySQL Insert Error",
		ErrMsg: "MySQL插入出错",
	},

	// 1100~1199 Redis错误
	EREDIS: &Error{
		LogMsg: "Redis Error",
		ErrMsg: "Redis数据库错误",
	},

	// 1200~1299 K8s错误
	EKUBE: &Error{
		LogMsg: "Kubernetes Error",
		ErrMsg: "Kubernetes错误",
	},
	EKUBE_CLIENT: &Error{
		LogMsg: "Kubernetes Create Client Error",
		ErrMsg: "创建K8s客户端出错",
	},
	EKUBE_CREATE_DEPLOYMENT: &Error{
		LogMsg: "Kubernetes Create Deployment Error",
		ErrMsg: "创建Deployment出错",
	},
	EKUBE_LIST_PODS: &Error{
		LogMsg: "Kubernetes List Pods Error",
		ErrMsg: "获取Pod列表出错",
	},

	// 1300~1399 Iris错误
	EIRIS: &Error{
		LogMsg: "Iris Error",
		ErrMsg: "Iris服务器错误",
	},

	// 1400~1499 YCE错误
	EYCE: &Error{
		LogMsg: "YCE Internal Error",
		ErrMsg: "YCE内部错误",
	},

	EYCE_LOGIN: &Error{
		LogMsg: "Can't Find the User",
		ErrMsg: "用户名密码错误",
	},

	EYCE_SESSION: &Error{
		LogMsg: "Can't Find the Session",
		ErrMsg: "请重新登录",
	},

	EYCE_SESSION_DEL: &Error{
		LogMsg: "Delete Session Error",
		ErrMsg: "退出遇到问题",
	},
	EYCE_ORG_EXIST: &Error{
		LogMsg: "The organization exists",
		ErrMsg: "组织已经存在",
	},

	// 1500~1599 Registr错误
	EREGISTRY: &Error{
		LogMsg: "Registry is empty, no docker image",
		ErrMsg: "私有镜像仓库为空",
	},
	EREGISTRY_GET: &Error{
		LogMsg: "Can't Get value from redis",
		ErrMsg: "不能检索镜像仓库",
	},

	// 1600 Json错误
	EJSON: &Error{
		LogMsg: "Json Marshal/Unmarshal Error",
		ErrMsg: "Json序列化错误",
	},
}
