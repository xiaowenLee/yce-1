package deploy

import (
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	organization "app/backend/common/yce/organization"
	myorganization "app/backend/model/mysql/organization"
	deploy "app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
)

type ListDeployController struct {
	*iris.Context
	org *myorganization.Organization
}

func (ldc *ListDeployController) getDcHost(orgId string) ([]string, error) {

	dcHost, err := organization.DcHost(orgId)
	if err != nil {
		log.Printf("Get dcHost error: orgId=%s, error=%s\n", orgId, err)
		return nil, err
	}

	return dcHost, nil
}

func (ldc *ListDeployController) getDcName(orgId string) ([]string, error) {
	dcName, err := organization.DcName(orgId)
	if err != nil {
		log.Printf("Get dcName error: orgId=%s, error=%s\n", orgId, err)
		return nil, err
	}
	return dcName, nil
}

func (ldc *ListDeployController) getPodList(dcName []string, dcHost []string, orgId string) (list string, err error) {

	tmpdata := make([]deploy.Data, len(dcHost))

	for i := 0; i < len(dcHost); i++ {
		newconfig := &restclient.Config{
			Host: dcHost[i],
		}
		newCli, err := client.New(newconfig)
		if err != nil {
			log.Printf("Get new restclient error: error=%s\n", err)
			return "", err
		}

		podlist, err := newCli.Pods(orgId).List(api.ListOptions{})
		if err != nil {
			log.Printf("Get podlist error: server=%s, orgId=%s, error=%s\n", dcHost[i], orgId, err)
			return "", err
		}

		tmpdata[i].DataCenter = dcName[i]
		tmpdata[i].PodList = *podlist
	}

	podListJson, err := json.Marshal(tmpdata)
	if err != nil {
		log.Printf("Get podListJson error: orgId=%s, error=%s\n", orgId, err)
	}

	list = string(podListJson)
	return list, nil
}
func (ldc ListDeployController) Get() {

	sessionIdClient := ldc.RequestHeader("sessionId")
	orgId := ldc.Param("orgId")
	//userId := ldc.Param("userId")

	tmpOrg, err := organization.GetOrganizationById(orgId)
	if err != nil {
		log.Printf("Get Organization By orgId error: orgId=%s, error=%s\n", orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		json, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(json)
		return
	}

	ldc.org = tmpOrg
	orgName := ldc.org.Name

	ss := session.SessionStoreInstance()

	if ok, err := ss.ValidateOrgId(sessionIdClient, orgId); ok {
		server, err := ldc.getDcHost(orgId)
		if err != nil {
			log.Printf("Get Datacenter Host error: sessionId=%s, orgId=%s, err=%s\n", sessionIdClient, orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			json, _ := ye.EncodeJson()
			ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ldc.Write(json)
			return
		}

		name, err := ldc.getDcName(orgId)
		if err != nil {
			log.Printf("Get Datacenter Name error: sessionId=%s, orgId=%s, err=%s\n", sessionIdClient, orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			json, _ := ye.EncodeJson()
			ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ldc.Write(json)
			return
		}

		podlist, err := ldc.getPodList(name, server, orgName)
		if err != nil {
			log.Printf("Get Podlist error: sessionId=%s, orgId=%s, error=%s\n", sessionIdClient, orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			json, _ := ye.EncodeJson()
			ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ldc.Write(json)
			return
		}

		ye := myerror.NewYceError(0, "OK", podlist)
		json, _ := ye.EncodeJson()

		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(json)
		return
	} else {
		log.Printf("Validate Session error: sessionId=%s, error=%s\n", sessionIdClient, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		json, _ := ye.EncodeJson()
		ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ldc.Write(json)
		return
	}

}
