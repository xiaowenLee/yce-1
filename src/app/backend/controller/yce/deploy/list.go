package deploy

import (
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	mylog "app/backend/common/util/log"
	"strconv"
)

//TODO: 添加错误码
type ListDeployController struct {
	*iris.Context
	org    *myorganization.Organization
	dcList []mydatacenter.DataCenter
	Ye *myerror.YceError
}

func (ldc *ListDeployController) WriteBack() {
	ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("ListDeployController Response YceError: controller=%p, code=%d, note=%s", ldc, ldc.Ye.Code, myerror.Errors[ldc.Ye.Code].LogMsg)
	ldc.Write(ldc.Ye.String())
}

// Return: [ dcHost1, dcHost2, ...] eg. ["172.21.1.11:8080", "10.149.149.3:8080", ...]
func (ldc *ListDeployController) getDcHost() []string {
	server := make([]string, len(ldc.dcList))

	for i := 0; i < len(server); i++ {
		server[i] = ldc.dcList[i].Host + ":" + strconv.Itoa(int(ldc.dcList[i].Port))
	}

	mylog.Log.Infof("ListDeployController getDcHost: ldc=%p, dcList=%p, server=%v", ldc, &ldc.dcList, server)
	return server
}

// Return: [ dcName1, dcName2, ...] eg. ["电信机房", "世纪互联", ...]
func (ldc *ListDeployController) getDcName() []string {
	name := make([]string, len(ldc.dcList))

	for i := 0; i < len(name); i++ {
		name[i] = ldc.dcList[i].Name
	}

	mylog.Log.Infof("ListDeployController getDcName: name=%v", name)
	return name
}

// Args: dcHostList <- ldc.getDcHost()
// Return: {"dcId": dcId, "dcName": dcName, "podList": json.Marshal(api.PodList)}
func (ldc *ListDeployController) getAppDisplayDeployment(dcHostList []string) (list string) {
	deployData := make([]deploy.AppDisplayDeployment, len(dcHostList))

	orgId := strconv.Itoa(int(ldc.org.Id))

	for i := 0; i < len(dcHostList); i++ {
		newconfig := &restclient.Config{
			Host: dcHostList[i],
		}
		newCli, err := client.New(newconfig)
		if err != nil {
			mylog.Log.Errorf("Get new restclient error: error=%s", err)
			ldc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
			return ""
		}

		podList, err := newCli.Pods(ldc.org.Name).List(api.ListOptions{})
		if err != nil {
			mylog.Log.Errorf("Get podlist error: server=%s, orgId=%s, error=%s", dcHostList[i], orgId, err)
			ldc.Ye = myerror.NewYceError(myerror.EKUBE_LIST_PODS, "")
			return ""
		}

		dcNameList := ldc.getDcName()

		deployData[i].DcName = dcNameList[i]
		deployData[i].DcId = ldc.dcList[i].Id
		deployData[i].PodList = *podList

		mylog.Log.Infof("ListDeployController getAppDisplayDeployment: dcId=%d, dcName=%s, podList=%p, len(podList.Items)=%d",
			ldc.dcList[i].Id, dcNameList[i], podList, len(podList.Items))
	}

	podListJson, err := json.Marshal(deployData)
	if err != nil {
		mylog.Log.Errorf("Get podListJson error: orgId=%s, error=%s", orgId, err)
		ldc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return ""
	}

	list = string(podListJson)
	return list
}

func (ldc *ListDeployController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		ldc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		// relogin
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		ldc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}
	return
}

// GET /api/v1/organizations/{orgId}/users/{userId}/deployments
func (ldc ListDeployController) Get() {
	sessionIdFromClient := ldc.RequestHeader("Authorization")
	orgId := ldc.Param("orgId")

	// Validate OrgId error
	ldc.validateSession(sessionIdFromClient, orgId)
	if ldc.Ye != nil {
		ldc.WriteBack()
		return
	}

	// Valid session
	org, err := organization.GetOrganizationById(orgId)
	ldc.org = org

	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ldc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		ldc.WriteBack()
		return
	}

	// Get Datacenters by a organization
	ldc.dcList, err = organization.GetDataCentersByOrganization(ldc.org)
	if err != nil {
		mylog.Log.Errorf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s", sessionIdFromClient, orgId, err)
		ldc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		ldc.WriteBack()
		return
	}


	// Get ApiServer for every datacenter
	server := ldc.getDcHost()

	// Get App DisplayDeployment
	appDisplayDeployment := ldc.getAppDisplayDeployment(server)
	if ldc.Ye != nil {
		ldc.WriteBack()
		return
	}

	ldc.Ye = myerror.NewYceError(myerror.EOK, appDisplayDeployment)
	ldc.WriteBack()

	mylog.Log.Infoln("ListDeploymentController Get over!")
	return
}
