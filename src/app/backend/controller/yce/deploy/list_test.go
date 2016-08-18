package deploy

import (
	mysql "app/backend/common/util/mysql"
	"fmt"
	"testing"
)

func Test_DcHost(t *testing.T) {

	mysqlclient := mysql.NewMysqlClient(mysql.DB_HOST, mysql.DB_USER, mysql.DB_PASSWORD, mysql.DB_NAME, mysql.MAX_POOL_SIZE)
	fmt.Println(mysqlclient)

	dcHost, err := getDcHost("1")
	if err != nil {
		t.Errorf("getDcHost error: error=%s\n", err)
	} else {
		fmt.Println(dcHost)
	}

	podlist, err := getPodList(dcHost, "1")
	if err != nil {
		t.Errorf("getPodList error: error=%s\n", err)
	} else {
		fmt.Println(podlist)
	}

}
