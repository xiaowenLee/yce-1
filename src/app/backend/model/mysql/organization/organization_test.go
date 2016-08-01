package organization

import (
	"fmt"
	"testing"
	mysql "app/backend/common/util/mysql"
	"github.com/shopspring/decimal"
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

func Test_DeleteOrganization(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	o := new(Organization)
	o.QueryOrganizationById(3)
	fmt.Printf("%v\n", o)
	o.DeleteOrganization(2)
	fmt.Printf("%v\n", o)
}

func Test_QueryBudgetById_UpdateBudgetById(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	o := new(Organization)
	budget, err := o.QueryBudgetById(1)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err.Error())
	}

	fmt.Printf("Budget before: %v\n", budget)

	fee, _ := decimal.NewFromString("1000")
	budget = budget.Sub(fee)

	str := budget.String()
	fmt.Printf("Budget after: %s\n", str)

	o.UpdateBudgetById(str, 2)
}

func Test_QueryBalanceById_UpdateBalanceById(*testing.T) {
	mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	mysql.MysqlInstance().Open()

	o := new(Organization)
	balance, err := o.QueryBalanceById(2)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		panic(err.Error())
	}

	fmt.Printf("Balance before: %v\n", balance)

	fee, _ := decimal.NewFromString("1000")
	balance = balance.Sub(fee)

	str := balance.String()
	fmt.Printf("Balance after: %s\n", str)

	o.UpdateBalanceById(str, 2)
}

func Test_EncodeJson_DecodeJson(*testing.T) {
	o := NewOrganization("dev", "1000000.00", "1000000.00", "add dev org", 1000, 2000, 1)
	fmt.Println(o.EncodeJson())

	org := new(Organization)
	org.DecodeJson(o.EncodeJson())
	fmt.Println(org.Name)
}
