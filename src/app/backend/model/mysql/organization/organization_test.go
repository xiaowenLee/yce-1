package organization

import (
	"fmt"
	"testing"
	mysql "app/backend/common/util/mysql"
)

func Test_NewOrganization(*testing.T) {
	o := NewOrganization("dev", "1000000.00", "1000000.00", "add dev org", 1000, 2000, 1)
	fmt.Printf("%v\n", o)
}

func Test_QueryOrganizationById(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	org := new(Organization)
	org.QueryOrganizationById(1)

	fmt.Printf("%v\n", org)

}

func Test_InsertOrganization(*testing.T) {

	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	o := NewOrganization("dev", "1000000.00", "1000000.00", "add dev org", 1000, 2000, 1)
	o.InsertOrganization(2)

	org := new(Organization)
	org.QueryOrganizationById(2)

	fmt.Printf("%v\n", org)

}