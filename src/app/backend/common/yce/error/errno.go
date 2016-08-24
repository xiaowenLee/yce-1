package error


type Error struct {
	LogMsg string
	ErrMsg string
}

const (
	EMYSQL int32 = 1000
	EREDIS int32 = 1100
	EKUBE	int32 = 1200
	EIRIS	int32 = 1300
	EYCE	int32 = 1400
)

const Errors = map[int32] *Error {

	// 1000~1099 MySQL错误
	EMYSQL: &Error{
		LogMsg: "MySQL Error",
		ErrMsg: "MySQL数据库错误",
	},

	// 1100~1199 Redis错误
	EREDIS: &Error{
		LogMsg: "Redis Error",
		ErrMsg: "Redis数据库错误",
	},

	// 1200~1299 K8s错误
	EKUBE:  &Error{
		LogMsg: "Kubernetes Error",
		ErrMsg: "Kubernetes错误",
	},

	// 1300~1399 Iris错误
	EIRIS:  &Error{
		LogMsg: "Iris Error",
		ErrMsg: "Iris服务器错误",
	},

	// 1400~1499 YCE错误
	EYCE:  &Error{
		LogMsg: "YCE Internal Error",
		ErrMsg: "YCE内部错误",
	},

}





