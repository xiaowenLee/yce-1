package organization

import (
	"strconv"
	"app/backend/model/mysql/organization"
	"app/backend/common/yce/datacenter"
	mylog "app/backend/common/util/log"
	mydatacenter "app/backend/model/mysql/datacenter"

	"encoding/json"
)

var log =  mylog.Log

type DcList struct {
	DataCenter []string `json:"dcList"`
}

func GetOrganizationById(orgId string) (*organization.Organization, error) {
	//mysqlclient := mysql.MysqlInstance()
	//mysqlclient.Open()

	myorganization := new(organization.Organization)
	oid, err := strconv.Atoi(orgId)
	if err != nil {
		log.Errorf("GetOrganizationById error: orgId=%s, error=%s", orgId, err)
		return nil, err
	}

	err = myorganization.QueryOrganizationById(int32(oid))
	if err != nil {
		log.Errorf("GetOrganizationById error: orgId=%s, error=%s", orgId, err)
		return nil, err
	}

	return myorganization, nil

}

func GetDataCentersByOrganization(org *organization.Organization) ([]mydatacenter.DataCenter, error){
	// Get datacenter-id-list for a organization(orgId)
	var dcList DcList

	err := json.Unmarshal([]byte(org.DcList), &dcList)
	if err != nil {
		log.Errorf("DecodeJSON error: dc=%s error=%s", dcList, err)
		return nil, err
	}

	orgId := org.Id
	// Get datacenters by dcId which in dcList
	dataCenters := make([]mydatacenter.DataCenter, len(dcList.DataCenter))

	for i := 0; i < len(dcList.DataCenter); i++ {
		dc, err := datacenter.GetDataCenterById(dcList.DataCenter[i])
		if err != nil {
			log.Errorf("Get Organization By orgId error: orgId=%s, error=%s", orgId, err)
			return nil, err
		}
		dataCenters[i] = *dc
	}

	return dataCenters, nil
}
