package deployment

import (
	mysql "app/backend/common/util/mysql"
	"fmt"
	"testing"
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

	d := NewDeployment("ncpay", "GET", "http://192.168.1.11:8080/namespaces/default/pods/", "shijihulian:dianxin", "null", "null", "null", 1, 2, 1)
	d.InsertDeployment(2)

	dp := new(Deployment)
	dp.QueryDeploymentById(1)
	fmt.Printf("%v\n", dp)
}

func Test_QueryDeploymentByAppName(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	app := "ncpay"
	deployments := QueryDeploymentByAppName(app)

	fmt.Printf("%v\n", deployments)

	for _, d := range *deployments {
		// fmt.Printf("%s\n", d.Name)
		fmt.Printf("%v\n", d)
	}
}
