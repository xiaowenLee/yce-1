package organization

import (
	// mysql "app/backend/common/util/mysql"
	// "app/backend/model/mysql/datacenter"
	"app/backend/model/mysql/organization"
	"log"
	"strconv"
)

type dcList struct {
	DcList []string `json:"dcList"`
}

func GetOrganizationById(orgId string) (*organization.Organization, error) {
	//mysqlclient := mysql.MysqlInstance()
	//mysqlclient.Open()

	myorganization := new(organization.Organization)
	oid, err := strconv.Atoi(orgId)
	if err != nil {
		log.Printf("GetOrganizationById error: orgId=%s, error=%s\n", orgId, err)
		return nil, err
	}

	err = myorganization.QueryOrganizationById(int32(oid))
	if err != nil {
		log.Printf("GetOrganizationById error: orgId=%s, error=%s\n", orgId, err)
		return nil, err
	}

	return myorganization, nil

}
