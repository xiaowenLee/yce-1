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
	"log"
	"strconv"
)

//TODO: 添加错误码

type ListDeployController struct {
	*iris.Context
	org    *myorganization.Organization
	dcList []mydatacenter.DataCenter
}

// Return: [ dcHost1, dcHost2, ...] eg. ["172.21.1.11:8080", "10.149.149.3:8080", ...]
func (ldc *ListDeployController) getDcHost() ([]string, error) {
	server := make([]string, len(ldc.dcList))

	for i := 0; i < len(server); i++ {
		server[i] = ldc.dcList[i].Host + ":" + strconv.Itoa(int(ldc.dcList[i].Port))
	}

	log.Printf("ListDeployController getDcHost: server=%v\n", server)
	return server, nil
}

// Return: [ dcName1, dcName2, ...] eg. ["电信机房", "世纪互联", ...]
func (ldc *ListDeployController) getDcName() ([]string, error) {
	name := make([]string, len(ldc.dcList))

	for i := 0; i < len(name); i++ {
		name[i] = ldc.dcList[i].Name
	}

	log.Printf("ListDeployController getDcName: name=%v\n", name)
	return name, nil
}

// Args: dcHostList <- ldc.getDcHost()
// Return: {"dcId": dcId, "dcName": dcName, "podList": json.Marshal(api.PodList)}
func (ldc *ListDeployController) getAppDisplayDeployment(dcHostList []string) (list string, err error) {
	deployData := make([]deploy.AppDisplayDeployment, len(dcHostList))

	orgId := strconv.Itoa(int(ldc.org.Id))

	for i := 0; i < len(dcHostList); i++ {
		newconfig := &restclient.Config{
			Host: dcHostList[i],
		}
		newCli, err := client.New(newconfig)
		if err != nil {
			log.Printf("Get new restclient error: error=%s\n", err)
			return "", err
		}

		podList, err := newCli.Pods(ldc.org.Name).List(api.ListOptions{})
		if err != nil {
			log.Printf("Get podlist error: server=%s, orgId=%s, error=%s\n", dcHostList[i], orgId, err)
			return "", err
		}

		dcNameList, err := ldc.getDcName()
		if err != nil {
			log.Printf("Get DcName error: server=%s, orgId=%s, error=%s\n", dcHostList[i], orgId, err)
			return "", err
		}

		deployData[i].DcName = dcNameList[i]
		deployData[i].DcId = ldc.dcList[i].Id
		deployData[i].PodList = *podList

		log.Printf("ListDeployController getAppDisplayDeployment: dcId=%d, dcName=%s, podList=%p, len(podList.Items)=%d\n",
			ldc.dcList[i].Id, dcNameList[i], podList, len(podList.Items))
	}

	podListJson, err := json.Marshal(deployData)
	if err != nil {
		log.Printf("Get podListJson error: orgId=%s, error=%s\n", orgId, err)
		return "", err
	}

	list = string(podListJson)

	// log.Printf("ListDeployController appDisplayDeployment Marshal over!")
	return list, nil
}

func (ldc *ListDeployController) valideSession(sessionId, orgId string) (*myerror.YceError, error){
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		log.Printf("Validate Session error: sessionId=%s, error=%s\n", sessionId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(errJson)
		return ye, err
	}

	// Session invalide
	if !ok {
		// relogin
		log.Printf("Validate Session failed: sessionId=%s, error=%s\n", sessionId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(errJson)
		return ye, err
	}

	return nil, nil
}

// GET /api/v1/organizations/{orgId}/users/{userId}/deployments
func (ldc ListDeployController) Get() {
	sessionIdFromClient := ldc.RequestHeader("Authorization")
	orgId := ldc.Param("orgId")

	// Validate OrgId error
	ye, err := ldc.valideSession(sessionIdFromClient, orgId)

	if ye != nil || err != nil {
		log.Printf("ListDeployController validateSession: sessionId=%s, orgId=%s, error=%s\n", sessionIdFromClient, orgId, err)
		errJson, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(errJson)
		return
	}

	// Valid session
	ldc.org, err = organization.GetOrganizationById(orgId)

	if err != nil {
		log.Printf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s\n", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(errJson)
		return
	}

	// Get Datacenters by a organization
	ldc.dcList, err = organization.GetDataCentersByOrganization(ldc.org)
	if err != nil {
		log.Printf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s\n", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(errJson)
		return
	}


	// Get ApiServer for every datacenter
	server, err := ldc.getDcHost()
	if err != nil {
		log.Printf("Get Datacenter Host error: sessionId=%s, orgId=%s, err=%s\n", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(errJson)
		return
	}

	// Get App DisplayDeployment
	appDisplayDeployment, err := ldc.getAppDisplayDeployment(server)
	if err != nil {
		log.Printf("Get Podlist error: sessionId=%s, orgId=%s, error=%s\n", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(errJson)
		return
	}

	ye = myerror.NewYceError(0, "OK", appDisplayDeployment)
	errJson, _ := ye.EncodeJson()
	ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ldc.Write(errJson)

	log.Printf("ListDeploymentController Get over!")
	return
}
