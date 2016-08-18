package deploy

import (
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	organization "app/backend/common/yce/organization"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
)

type ListDeployController struct {
	cli *client.Client
	*iris.Context
}

func NewListDeployController(server string) *ListDeployController {
	config := &restclient.Config{
		Host: server,
	}
	cli, err := client.New(config)
	if err != nil {
		log.Printf("Get ListDeployController error: error=%s\n", err)
	}

	instance := &ListDeployController{cli: cli}
	return instance
}

func (ldc *ListDeployController) getDcHost(orgId string) ([]string, error) {
	dcHost, err := organization.DcHost(orgId)
	if err != nil {
		log.Printf("Get dcList error: orgId=%s, error=%s\n", orgId, err)
		return nil, err
	}

	return dcHost, nil
}

func (ldc *ListDeployController) getPodList(dcHost []string, orgId string) (list string, err error) {

	for _, v := range dcHost {
		newconfig := &restclient.Config{
			Host: v,
		}
		newCli, err := client.New(newconfig)
		if err != nil {
			log.Printf("Get new restclient error: error=%s\n", err)
			return "", err
		}

		podlist, err := newCli.Pods(orgId).List(api.ListOptions{})
		if err != nil {
			log.Printf("Get podlist error: server=%s, orgId=%s, error=%s\n", v, orgId, err)
			return "", err
		}

		podListJson, err := json.Marshal(podlist)
		if err != nil {
			log.Printf("Get podListJson error: server=%s, orgId=%s, error=%s\n", v, orgId, err)
		}

		list += string(podListJson)
	}
	return list, nil
}

func (ldc ListDeployController) Get() {

	sessionIdClient := ldc.RequestHeader("sessionId")
	orgId := ldc.RequestHeader("orgId")
	orgName := ldc.Param("orgName")
	userId := ldc.RequestHeader("uid")

	log.Printf("sessionIdClient=%s, orgId=%s, orgName=%s, userId=%s\n", sessionIdClient, orgId, orgName, userId)

	ss := session.SessionStoreInstance()

	if ok, err := ss.ValidateOrgId(sessionIdClient, orgId); ok {
		server, err := ldc.getDcHost(orgId)
		if err != nil {
			log.Printf("Get Datacenter Host error: sessionId=%s, orgId=%s, err=%s\n", sessionIdClient, orgId, err)
		}

		podlist, err := ldc.getPodList(server, orgName)
		if err != nil {
			log.Printf("Get Podlist error: sessionId=%s, orgId=%s, error=%s\n", sessionIdClient, orgId, err)
		}
		ye := myerror.NewYceError(0, "OK", podlist)
		json, _ := ye.EncodeJson()

		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(json)
	} else {
		log.Printf("Validate Session error: sessionId=%s, error=%s\n", sessionIdClient, err)
	}

}
