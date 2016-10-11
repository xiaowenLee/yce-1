package deployment

import (
	"testing"
	mysql "app/backend/common/util/mysql"
	config "app/backend/common/yce/config"
)

/*
func Test_NewDeployment(*testing.T) {
	d := NewDeployment("ncpay", "GET", "http://192.168.1.11:8080/namespaces/default/pods/", "shijihulian:dianxin", "null", "null", "null", 1, 2, VALID)
	fmt.Printf("%v\n", d)
}

func Test_QueryDeploymentById(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	d := new(Deployment)
	d.QueryDeploymentById(1)

	fmt.Printf("%v\n", d)

}

func Test_QueryDeploymentByAppName(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	app := "ncpay"
	deployments, _ := QueryDeploymentByAppName(app)

	fmt.Printf("%v\n", deployments)

	for _, d := range *deployments {
		// fmt.Printf("%s\n", d.Name)
		fmt.Printf("%v\n", d)
	}
}
*/
/*
func Test_InsertDeployment(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	d := NewDeployment("ncpay", "GET", "http://192.168.1.11:8080/namespaces/default/pods/", "shijihulian:dianxin", "null", "null", "null", 1, 2, 1)
	d.InsertDeployment(2)

	dp := new(Deployment)
	dp.QueryDeploymentById(1)
	fmt.Printf("%v\n", dp)
}

func Test_StatDeploymentByActionType(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	cnt, _ := StatDeploymentByActionType(2)
	fmt.Println(cnt)
}
*/

func Test_QueryOperationStat(*testing.T) {
	/*
	mysql.NewMysqlClient(config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.MAX_POOL_SIZE)
	*/

	config.Instance().Load()
	mysql.MysqlInstance().Open()

	// days := []string{"2016-10-11", "2016-10-10", "2016-10-09", "2016-10-08", "2016-10-07", "2016-10-06", "2016-10-05"}

	var orgId int32
	orgId = 1

	ops, _ := QueryOperationStat(orgId)

	for k, v := range ops {
		log.Infof("OperationStat: key=%s, value=%v", k, v)
	}


}
