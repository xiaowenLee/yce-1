package organization

import (
	"strconv"
	"app/backend/model/mysql/organization"
	"app/backend/common/yce/datacenter"
	mylog "app/backend/common/util/log"
	mydatacenter "app/backend/model/mysql/datacenter"
	"encoding/json"
)


type DcIdListType struct {
	DcIdList []int32 `json:"dcIdList"`
}

func GetOrganizationById(orgId string) (*organization.Organization, error) {
	myorganization := new(organization.Organization)
	oid, err := strconv.Atoi(orgId)
	if err != nil {
		mylog.Log.Errorf("GetOrganizationByOrgID Error: orgId=%s, error=%s", orgId, err)
		return nil, err
	}

	err = myorganization.QueryOrganizationById(int32(oid))
	if err != nil {
		mylog.Log.Errorf("GetOrganizationByOrgID Error: orgId=%s, error=%s", orgId, err)
		return nil, err
	}


	return myorganization, nil

}

func GetDataCentersByOrganization(org *organization.Organization) ([]mydatacenter.DataCenter, error) {
	// Get datacenter-id-list for a organization(orgId)
	var dcIdList DcIdListType

	orgId := org.Id
	err := json.Unmarshal([]byte(org.DcIdList), &dcIdList)
	if err != nil {
		mylog.Log.Errorf("GetDataCentersByOrg Error: orgId=%s, error=%s", orgId, err)
		return nil, err
	}

	// Get datacenters by dcId which in dcIdList
	dataCenters := make([]mydatacenter.DataCenter, 0)

	for i := 0; i < len(dcIdList.DcIdList); i++ {
		dc, err := datacenter.GetDataCenterById(dcIdList.DcIdList[i])
		if err != nil {
			mylog.Log.Errorf("GetDataCentersByOrg Error: orgId=%s, error=%s", orgId, err)
			return nil, err
		}
		dataCenters = append(dataCenters, *dc)
	}
	mylog.Log.Infof("GetDataCentersByOrganization len(datacenters)=%d", len(dataCenters))
	return dataCenters, nil
}
