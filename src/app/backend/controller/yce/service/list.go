package service

import (
	"log"
	"github.com/kataras/iris"
	"app/backend/common/util/session"
	"app/backend/common/yce/organization"
	"app/backend/model/yce/service"
	myerror "app/backend/common/yce/error"
	myorganization "app/backend/model/mysql/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	"encoding/json"
	"strconv"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

type ListServiceController struct {
	*iris.Context
	org *myorganization.Organization
	dcList []mydatacenter.DataCenter
}



func (lsc *ListServiceController) validateSession(sessionId, orgId string) (*myerror.YceError, error) {
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		log.Printf("Validate Session error: sessionId=%s, error=%s\n", sessionId, err)
		ye := myerror.NewYceError(1, "请求失败")
		errJson, _ := ye.EncodeJson()
		lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		lsc.Write(errJson)
		return ye, err
	}

	// Session invalide
	if !ok {
		// relogin
		log.Printf("Validate Session failed: sessionId=%s, error=%s\n", sessionId, err)
		ye := myerror.NewYceError(1, "请求失败")
		errJson, _ := ye.EncodeJson()
		lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		lsc.Write(errJson)
		return ye, err
	}

	log.Printf("ListServiceController validate sessionId with orgId ok: sessionId=%s, orgId=%s\n", sessionId, orgId)
	return nil, nil
}

func (lsc *ListServiceController) getDcHost() ([]string, error) {
	server := make([]string, len(lsc.dcList))

	for i := 0; i < len(lsc.dcList); i++ {
		server[i] = lsc.dcList[i].Host + ":" + strconv.Itoa(int(lsc.dcList[i].Port))
	}

	log.Printf("ListServiceController getDcHost: server=%v\n", server)
	return server, nil
}

func (lsc *ListServiceController) getDisplayServices(dcHostList []string) (list string, err error) {
	serviceData := make([]service.Service, len(dcHostList))

	orgId := strconv.Itoa(int(lsc.org.Id))

	for i := 0; i < len(dcHostList); i++ {
		newconfig := &restclient.Config{
			Host: dcHostList[i],
		}

		newCli, err := client.New(newconfig)
		if err != nil  {
			log.Printf("Get new restclient error: error=%s\n", err)
			return "", err
		}

		svcList, err := newCli.Services(lsc.org.Name).List(api.ListOptions{})
		if err != nil {
			log.Printf("Get serviceList error: error=%s\n", err)
			return "", err
		}

		serviceData[i].DcId = lsc.dcList[i].Id
		serviceData[i].DcName = lsc.dcList[i].Name
		serviceData[i].ServiceList= *svcList

		log.Printf("ListServiceController getDisplayService: dcId=%d, dcName=%s, serviceList=%p, len(serviceList)=%d\n", serviceData[i].DcId, serviceData[i].DcName, &serviceData, len(serviceData[i].ServiceList.Items))
	}

	serviceListJson, err := json.Marshal(serviceData)

	if err != nil {
		log.Printf("Get serviceListJson error: orgId=%s, error=%s\n", orgId, err)
		return "", err
	}

	log.Printf("ListServiceController Get serviceList ok: orgId=%s\n", orgId)
	list = string(serviceListJson)
	return list, nil
}


//GET /api/v1/organizations/{orgId}/users/{userId}/services
func (lsc ListServiceController) Get() {
	sessionIdFromClient := lsc.RequestHeader("Authorization")
	orgId := lsc.Param("orgId")

	// Validate orgId error
	ye, err := lsc.validateSession(sessionIdFromClient, orgId)
	if ye != nil || err != nil {
		log.Printf("ListSerivceController validateSession error: sessionId=%s, error=%s\n", sessionIdFromClient, err)
		errJson, _ := ye.EncodeJson()
		lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		lsc.Write(errJson)
		return
	}


	// Valid session
	lsc.org, err = organization.GetOrganizationById(orgId)

	if err != nil {
		log.Printf("Get Organization By orgId error: sessionId=%s, orgId=%s, error=%s\n", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "请求失败")
		errJson, _ := ye.EncodeJson()
		lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		lsc.Write(errJson)
		return
	}

	// Get Datacenters by a organization
	lsc.dcList, err = organization.GetDataCentersByOrganization(lsc.org)
	if err != nil {
		log.Printf("Get Datacenters By Organization error: sessionId=%s, orgId=%s, error=%s\n", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "请求失败")
		errJson, _ := ye.EncodeJson()
		lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		lsc.Write(errJson)
		return
	}

	// Get ApiServer for every datacenter
	server, err := lsc.getDcHost()
	if err != nil {
		log.Printf("Get Datacenter Host error: sessionId=%s, orgId=%s\n", sessionIdFromClient, orgId)
		ye := myerror.NewYceError(1, "请求失败")
		errJson, _ := ye.EncodeJson()
		lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		lsc.Write(errJson)
		return
	}

	// Get DisplayServices
	displayServices, err := lsc.getDisplayServices(server)
	if err != nil {
		log.Printf("Get ServiceList error: sessionId=%s, orgId=%s, error=%s\n", sessionIdFromClient, orgId, err)
		ye := myerror.NewYceError(1, "请求失败")
		errJson, _ := ye.EncodeJson()
		lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		lsc.Write(errJson)
		return
	}

	ye = myerror.NewYceError(0, displayServices)
	errJson, _ :=ye.EncodeJson()
	lsc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	lsc.Write(errJson)

	log.Printf("ListServiceController Get over!")
	return
}