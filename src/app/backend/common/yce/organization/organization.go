package organization

import (
	mysql "app/backend/common/util/mysql"
	"app/backend/model/mysql/datacenter"
	"app/backend/model/mysql/organization"
	"encoding/json"
	"log"
	"strconv"
)

type dcList struct {
	DcList []string `json:"dcList"`
}

func DcHost(orgId string) ([]string, error) {
	mysqlclient := mysql.MysqlInstance()
	mysqlclient.Open()

	myorganization := new(organization.Organization)
	oid, err := strconv.Atoi(orgId)
	if err != nil {
		log.Printf("Get orgId error: orgId=%s, error=%s\n", orgId, err)
	}

	err = myorganization.QueryOrganizationById(int32(oid))
	if err != nil {
		log.Printf("Get dcList error: orgId=%s, error=%s\n", orgId, err)
		return nil, err
	}

	dclist := new(dcList)
	err = json.Unmarshal([]byte(myorganization.DcList), &dclist)
	if err != nil {
		log.Printf("Decode dcList error: orgId=%s, error=%s\n", orgId, err)
		return nil, err
	}

	mydatacenter := new(datacenter.DataCenter)

	num := len(dclist.DcList)
	server := make([]string, num)

	for i := 0; i < num; i++ {
		id, err := strconv.Atoi(dclist.DcList[i])
		if err != nil {
			log.Printf("Strconv.Atoi dcList error: orgId=%s, error=%s\n", orgId, err)
		}
		err = mydatacenter.QueryDataCenterById(int32(id))
		if err != nil {
			log.Printf("Get dcHost error: OrgID=%s, error=%s\n", orgId, err)
		}
		server[i] = mydatacenter.Host + ":" + strconv.Itoa(int(mydatacenter.Port))
	}

	return server, nil
}
