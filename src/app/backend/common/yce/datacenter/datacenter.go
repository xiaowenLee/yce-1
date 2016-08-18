package datacenter

import (
	"log"
	"strconv"
	"app/backend/model/mysql/datacenter"
)

func GetDataCenterById(dcId string) (*datacenter.DataCenter, error) {
	//mysqlclient := mysql.MysqlInstance()
	//mysqlclient.Open()

	mydatacenter := new(datacenter.DataCenter)
	dcid, err := strconv.Atoi(dcId)
	if err != nil {
		log.Printf("GetDataCenterById error: dcId=%s, error=%s\n", dcId, err)
		return nil, err
	}

	err = mydatacenter.QueryDataCenterById(int32(dcid))
	if err != nil {
		log.Printf("GetDataCenterById error: dcId=%s, error=%s\n", dcId, err)
		return nil, err
	}

	return mydatacenter, nil
}