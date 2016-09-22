package datacenter

import (
	mylog "app/backend/common/util/log"
	"app/backend/model/mysql/datacenter"
)

var log =  mylog.Log

func GetDataCenterById(dcId int32) (*datacenter.DataCenter, error) {

	mydatacenter := new(datacenter.DataCenter)
	err := mydatacenter.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("GetDataCenterById error: dcId=%s, error=%s", dcId, err)
		return nil, err
	}

	return mydatacenter, nil
}