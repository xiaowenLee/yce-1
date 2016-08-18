package organization

import (
	"log"
	"strconv"
	"app/backend/model/mysql/organization"
)

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
