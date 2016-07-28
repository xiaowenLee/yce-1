package deployment

import (
	"fmt"
	"testing"
	mysql "app/backend/common/util/mysql"
)

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

func Test_InsertDeployment(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	d := NewDeployment("ncpay", "GET", "http://192.168.1.11:8080/namespaces/default/pods/", "shijihulian:dianxin", "null", "null", "null", 1, 2, VALID)
	d.InsertDeployment()

	dp := new(Deployment)
	dp.QueryDeploymentById(1)
	fmt.Printf("%v\n", dp)
}

func Test_QueryDeploymentByAppName(*testing.T) {
	app := "ncpay"
	deployments := QueryDeploymentByAppName(app)

	for d := range *deployments {
		fmt.Printf("%v\n", d)
	}
}