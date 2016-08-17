package organization

import (
	mysql "app/backend/common/util/mysql"
	"app/backend/model/mysql/organization"
	"encoding/json"
	"log"
	"strconv"
)

type dcList struct {
	DcList []string `json:"dcList"`
}

func DcList(OrgId string) ([]string, error) {
	mysqlclient := mysql.MysqlInstance()
	mysqlclient.Open()

	myorganization := new(organization.Organization)
	oid, err := strconv.Atoi(OrgId)
	if err != nil {
		log.Printf("Get OrgId error: OrgId=%s, error=%s\n", OrgId, err)
	}

	err = myorganization.QueryOrganizationById(int32(oid))
	if err != nil {
		log.Printf("Get dcList error: OrgId=%s, error=%s\n", OrgId, err)
		return nil, err
	}

	dclist := new(dcList)
	err = json.Unmarshal([]byte(myorganization.DcList), &dclist)
	if err != nil {
		log.Printf("Decode dcList error: OrgId=%s, error=%s\n", OrgId, err)
		return nil, err
	}

	log.Printf("DcList: dclist=%s\n", dclist)

	return dclist.DcList, nil
}
