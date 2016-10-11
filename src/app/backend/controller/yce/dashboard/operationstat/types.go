package operationstat

import (
	mydeployment "app/backend/model/mysql/deployment"
	mylog "app/backend/common/util/log"
	myday "app/backend/common/util/day"
	"encoding/json"
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

type OperationStatistics struct {
	Date DateType `json:"date"`
	Statistics StatisticsType `json:"statistics"`
	mapping map[string]int32
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

	oss.mapping = make(map[string]int32)

	return oss
}

func (ops *OperationStatistics) bindDateArray() string {
	str := ""
	for _, day := range ops.Date {
		str += day + ","
	}
	return str
}

func (ops *OperationStatistics) InitOperationStatistics() {
	d := myday.NewToday()
	ds := d.GetLastDaysString(LAST_N_DAYS)

	for index, v := range ds {
		ops.Date = append(ops.Date, v)
		ops.mapping[v] = int32(index)
	}

	for i := 0; i < LAST_N_DAYS; i++ {
		ops.Statistics.Online = append(ops.Statistics.Online, 0)
		ops.Statistics.Rollingupgrade = append(ops.Statistics.Rollingupgrade, 0)
		ops.Statistics.Rollback = append(ops.Statistics.Rollback, 0)
		ops.Statistics.Scale = append(ops.Statistics.Scale, 0)
		ops.Statistics.Delete = append(ops.Statistics.Delete, 0)
	}
}

func  (ops *OperationStatistics) Transform(orgId int32) (string, error) {
	ost, err := mydeployment.QueryOperationStat(orgId)

	if err != nil {
		return "", err
	}

	oss := NewOperationStatistics()
	oss.InitOperationStatistics()

	for date, v := range ost {
		if index, ok := oss.mapping[date]; ok {
			for op, total := range v {
				switch op {
				case 2:
					oss.Statistics.Online[index] = total
					break
				case 3:
					oss.Statistics.Rollback[index] = total
					break
				case 4:
					oss.Statistics.Rollingupgrade[index] = total
					break
				case 8:
					oss.Statistics.Scale[index] = total
					break
				case 9:
					oss.Statistics.Delete[index] = total
					break

				}
			}
		}
	}

	data, err := json.Marshal(oss)

	if err != nil {
		log.Errorf("OperationStatistics Transform Json Error: err=%s", err)
		return "", err
	}

	return string(data), nil
}