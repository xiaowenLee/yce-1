package resourcestat

import (
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mydatacenter "app/backend/model/mysql/datacenter"
	"encoding/json"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	"strings"
)

type StatResourceController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
	orgId      string
	orgName    string
}

// get one Datacenter resource statistics
func (src *StatResourceController) statDatacenter(datacenter *mydatacenter.DataCenter, cli *client.Client) *DatacentersType {
	dc := new(DatacentersType)

	dc.DcId = datacenter.Id
	dc.DcName = datacenter.Name


	// get deployments in this datacenters
	dpList, ye := yceutils.GetDeploymentByNamespace(cli, src.orgName)
	if ye != nil {
		src.Ye = ye
		return nil
	}

	for _, dp := range dpList {
		// stat CPU
		c := dp.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()
		log.Debugf("StatResourceController statDatacenter cpu used: cpu=%s", c)
		cu, _ := strconv.Atoi(c)
		n := dp.Spec.Replicas
		dc.Cpu.Used += int32(cu) * n

		// stat MEM
		m := dp.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()
		log.Debugf("StatResourceController statDatacenter mem used: mem=%s", m)
		str := strings.Split(m, "G")
		mu, _ := strconv.Atoi(str[0])
		dc.Mem.Used += int32(mu) * n

	}

	return dc
}

// get datacenters' resource statistics
func (src *StatResourceController) statDatacenters() *ResourceStat {
	resStat := new(ResourceStat)
	resStat.Datacenters = make([]DatacentersType, 0)

	datacenterList, ye := yceutils.GetDatacentersByOrgId(src.orgId)
	if ye != nil {
		src.Ye = ye
		return nil
	}

	overall := new(DatacentersType)
	overall.DcId = -1
	overall.DcName = "总览"

	org, err := organization.GetOrganizationById(src.orgId)
	src.orgName = org.Name
	if err != nil {
		src.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return nil
	}

	// stat overall CPU
	overall.Cpu.Total = org.CpuQuota
	// stat overall Mem
	overall.Mem.Total = org.MemQuota

	tmpDcStat := make([]DatacentersType, 0)
	for index, datacenter := range datacenterList {
		dcStat := src.statDatacenter(&datacenter, src.k8sClients[index])
		if src.CheckError() {
			return nil
		}

		dcStat.Cpu.Total = overall.Cpu.Total / int32(len(datacenterList))
		dcStat.Mem.Total = overall.Mem.Total / int32(len(datacenterList))

		overall.Cpu.Used += dcStat.Cpu.Used
		overall.Mem.Used += dcStat.Mem.Used

		tmpDcStat = append(tmpDcStat, *dcStat)
	}

	resStat.Datacenters = append(resStat.Datacenters, *overall)
	for _, dcStat := range tmpDcStat {
		resStat.Datacenters = append(resStat.Datacenters, dcStat)
	}

	return resStat

}

// get Resource statistics. It's the main purpose.
func (src *StatResourceController) getResourceStat() string {

	resStat := src.statDatacenters()
	resStatJSON, err := json.Marshal(resStat.Datacenters)
	if err != nil {
		log.Errorf("StatResourceController getResourceStat Error: error=%s", err)
		src.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	resStatString := string(resStatJSON)

	return resStatString

}

func (src StatResourceController) Get() {
	sessionIdFromClient := src.RequestHeader("Authorization")
	orgId := src.Param("orgId")

	src.ValidateSession(sessionIdFromClient, orgId)
	if src.CheckError() {
		return
	}

	src.orgId = orgId

	// get dcIdList
	dcIdList, ye := yceutils.GetDcIdListByOrgId(src.orgId)
	if ye != nil {
		src.Ye = ye
	}
	if src.CheckError() {
		return
	}

	// get apiServer
	src.apiServers, src.Ye = yceutils.GetApiServerList(dcIdList)
	if src.CheckError() {
		return
	}

	// get k8sClient
	src.k8sClients, src.Ye = yceutils.CreateK8sClientList(src.apiServers)
	if src.CheckError() {
		return
	}

	// get ResourceStat
	resStatString := src.getResourceStat()
	if src.CheckError() {
		return
	}

	src.WriteOk(resStatString)
	log.Infoln("StatResourceController get over")

	return
}
