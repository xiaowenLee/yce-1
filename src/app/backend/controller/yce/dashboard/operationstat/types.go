package operationstat

import (
	mydeployment "app/backend/model/mysql/deployment"
	mylog "app/backend/common/util/log"
	"sort"
	"time"
)

var log =  mylog.Log

type DateType []string

type StatType []int32

type StatisticsType struct {
	Online StatType `json:"online"`
	Scale StatType `json:"scale"`
	Rollingupgrade StatType `json:"rollingupgrade"`
	Rollback StatType `json:"rollback"`
	Delete StatType `json:"delete"`
}

func (d DateType) Len() int {
	return len(d)
}

func (d DateType) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DateType) Less(i, j int) bool {
	ti, _ := time.Parse(d[i], "2011-01-19")
	tj, _ := time.Parse(d[j], "2011-01-19")

	return ti <= tj
}

type OperationStatistics struct {
	Date DateType `json:"date"`
	Statistics StatisticsType `json:"statistics"`
}

func getLastDay(d DateType) DateType {

}

func NewOperationStatistics() *OperationStatistics {

	oss := new(OperationStatistics)
	oss.Date = make(DateType, 0)
	oss.Statistics = *new(StatisticsType)
	oss.Statistics.Online = make(StatType, 0)
	oss.Statistics.Rollingupgrade = make(StatType, 0)
	oss.Statistics.Rollback = make(StatType, 0)
	oss.Statistics.Scale = make(StatType, 0)
	oss.Statistics.Delete = make(StatType, 0)

	return oss
}

func  Transform() (string, error) {

	ops, err := mydeployment.QueryOperationStat()

	if err != nil {
		return "", err
	}

	oss := NewOperationStatistics()

	return "", nil
}