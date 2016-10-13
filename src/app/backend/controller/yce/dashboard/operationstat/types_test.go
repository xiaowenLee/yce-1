package operationstat

import (
	"testing"
	"fmt"
	mysql "app/backend/common/util/mysql"
	config "app/backend/common/yce/config"
)

func TestOperationStatistics_Transform(t *testing.T) {
	ops := NewOperationStatistics()

	config.Instance().Load()
	mysql.MysqlInstance().Open()

	str, _ :=  ops.Transform(1)
	fmt.Println(str)
}