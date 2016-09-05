package datacenter

import (
	"strconv"
	mylog "app/backend/common/util/log"
	"app/backend/model/mysql/datacenter"
)

var log =  mylog.Log

func GetDataCenterById(dcId string) (*datacenter.DataCenter, error) {
	//mysqlclient := mysql.MysqlInstance()
	//mysqlclient.Open()

	mydatacenter := new(datacenter.DataCenter)
	dcid, err := strconv.Atoi(dcId)
	if err != nil {
		log.Errorf("GetDataCenterById error: dcId=%s, error=%s", dcId, err)
		return nil, err
	}

	err = mydatacenter.QueryDataCenterById(int32(dcid))
	if err != nil {
		log.Errorf("GetDataCenterById error: dcId=%s, error=%s", dcId, err)
		return nil, err
	}

	return mydatacenter, nil
}