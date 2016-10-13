package deploymentstat

import (
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	client "k8s.io/kubernetes/pkg/client/unversioned"

	myerror "app/backend/common/yce/error"
	"encoding/json"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type StatDeploymentController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
	orgId      string
	orgName    string
}

// get deployment statistics including new ReplicaSet's name and Pods' name.
func (sdc *StatDeploymentController) statDeployment(dp *extensions.Deployment, cli *client.Client) *DeploymentsType {
	dps := new(DeploymentsType)

	// get rsName
	rsNew, ye := yceutils.GetNewReplicaSetByDeployment(cli, dp)
	if ye != nil {
		sdc.Ye = ye
		return nil
	}
	log.Debugf("StatDeploymentController statDeployment: rsName=%s", rsNew.Name)

	// get podName
	podList, ye := yceutils.GetPodListByReplicaSet(cli, rsNew)
	if ye != nil {
		sdc.Ye = ye
		return nil
	}
	log.Debugf("StatDeploymentController statDeployment: len(podList)=%d", len(podList.Items))

	dps.DeploymentName = dp.Name
	dps.RsName = rsNew.Name
	dps.PodName = make([]string, 0)
	for _, pod := range podList.Items {
		dps.PodName = append(dps.PodName, pod.Name)
	}
	return dps
}

// get datacenters statistics especially datacenter's Id and name
func (sdc *StatDeploymentController) getDatacenterStat(index int, dpList []extensions.Deployment, cli *client.Client) []DeploymentsType {
	dcs := make([]DeploymentsType, 0)

	dcList, ye := yceutils.GetDatacentersByOrgId(sdc.orgId)
	if ye != nil {
		sdc.Ye = ye
		return nil
	}
	log.Debugf("StatDeploymentController getDatacenterStat: len(dcList)=%d", len(dcList))

	for _, dp := range dpList {

		dps := sdc.statDeployment(&dp, cli)
		if sdc.CheckError() {
			return nil
		}

		dps.DcId = dcList[index].Id
		dps.DcName = dcList[index].Name
		dcs = append(dcs, *dps)
	}

	return dcs
}

// get DeploymentStat. It's the main purpose.
func (sdc *StatDeploymentController) getDeploymentStat() string {

	dpStat := new(DeploymentStatType)
	dpStat.DeploymentStat = make([]DeploymentsType, 0)

	for index, cli := range sdc.k8sClients {
		// get deployments
		dpList, ye := yceutils.GetDeploymentByNamespace(cli, sdc.orgName)
		if ye != nil {
			sdc.Ye = ye
			return ""
		}
		log.Debugf("StatDeploymentController getDeploymentStat: len(dpList)=%d", len(dpList))

		dcs := sdc.getDatacenterStat(index, dpList, cli)
		if sdc.CheckError() {
			return ""
		}

		for _, dc := range dcs {
			dpStat.DeploymentStat = append(dpStat.DeploymentStat, dc)
		}
	}

	// encode to json then convert to string
	dpStatJSON, err := json.Marshal(dpStat.DeploymentStat)
	if err != nil {
		log.Errorf("StatDeploymentController getDeploymentStat Error: error=%s", err)
		sdc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}
	log.Infoln("StatDeploymentController getDeploymentStat Success")
	log.Debugln(dpStatJSON)
	dpStatString := string(dpStatJSON)

	return dpStatString
}

func (sdc StatDeploymentController) Get() {
	sessionIdFromClient := sdc.RequestHeader("Authorization")
	orgId := sdc.Param("orgId")
	sdc.orgId = orgId

	log.Debugf("StatDeploymentController get Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)

	// validate session by orgId
	sdc.ValidateSession(sessionIdFromClient, orgId)
	if sdc.CheckError() {
		return
	}

	// get datacenters by orgId
	dcIdList, ye := yceutils.GetDcIdListByOrgId(orgId)
	if ye != nil {
		sdc.Ye = ye
	}
	if sdc.CheckError() {
		return
	}

	// get apiServers
	sdc.apiServers, sdc.Ye = yceutils.GetApiServerList(dcIdList)
	if sdc.CheckError() {
		return
	}

	// get k8sClients
	sdc.k8sClients, sdc.Ye = yceutils.CreateK8sClientList(sdc.apiServers)
	if sdc.CheckError() {
		return
	}

	// get namespace
	sdc.orgName, sdc.Ye = yceutils.GetOrgNameByOrgId(orgId)
	if sdc.CheckError() {
		return
	}

	// get deployment statistics
	dpStatString := sdc.getDeploymentStat()
	if sdc.CheckError() {
		return
	}

	// write back
	sdc.WriteOk(dpStatString)
	log.Infoln("StatDeploymentController Get over !")

	return
}
