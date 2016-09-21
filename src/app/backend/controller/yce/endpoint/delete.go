package endpoint
import (
	myerror "app/backend/common/yce/error"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type DeleteEndpointsController struct {
	yce.Controller
	k8sClients []*client.Client
	apiServers []string
}


// Publish k8s.Endpoint to every datacenter which in dcIdList
func (dec *DeleteEndpointsController) deleteEndpoints(namespace, epName string) {
	// Foreach every K8sClient to create service

	for index, cli := range dec.k8sClients {
		err := cli.Endpoints(namespace).Delete(epName)
		if err != nil {
			log.Errorf("deleteEndpoint Error: apiServer=%s, namespace=%s, error=%s", dec.apiServers[index], namespace, err)
			//dec.Ye = myerror.NewYceError(myerror.EKUBE_CREATE_SERVICE, "")
			dec.Ye = myerror.NewYceError(myerror.EKUBE_DELETE_ENDPOINT, "")
			return
		}

		log.Infof("Delete Endpoint successfully: namespace=%s, apiServer=%s", namespace, dec.apiServers[index])
	}
	log.Infof("DeleteEndpointsController Delete Endpoints success")
	return
}

func (dec DeleteEndpointsController) Delete() {
	sessionIdFromClient := dec.RequestHeader("Authorization")
	orgId := dec.Param("orgId")
	dcId := dec.Param("dcId")
	epName := dec.Param("epName")

	log.Debugf("DeleteEndpontsController Params: sessionId=%s, orgId=%s, dcId=%s, epName=%s", sessionIdFromClient, orgId, dcId, epName)


	// Validate OrgId error
	dec.ValidateSession(sessionIdFromClient, orgId)
	if dec.CheckError() {
		return
	}

	// Get DcIdList
	dcIdList := make([]int32, 0)
	datacenterId, _ := strconv.Atoi(dcId)
	dcIdList = append(dcIdList, int32(datacenterId))
	log.Debugf("DeleteEndpointController len(DcIdList)=%d", len(dcIdList))

	dec.apiServers, dec.Ye = yceutils.GetApiServerList(dcIdList)
	if dec.CheckError() {
		return
	}

	// Create K8sClient List
	dec.k8sClients, dec.Ye = yceutils.CreateK8sClientList(dec.apiServers)
	if dec.CheckError() {
		return
	}

	// Publish server to every datacenter
	orgName,  ye := yceutils.GetOrgNameByOrgId(orgId)
	if ye != nil {
		dec.Ye = ye
	}
	if dec.CheckError() {
		return
	}

	dec.deleteEndpoints(orgName, epName)
	if dec.CheckError() {
		return
	}

	dec.WriteOk("")
	log.Infoln("DeleteEndpointsController over!")
	return
}
