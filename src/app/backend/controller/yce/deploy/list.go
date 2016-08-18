package deploy

import (
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	organization "app/backend/common/yce/organization"
	datacenter "app/backend/common/yce/datacenter"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
	deploy "app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"log"
	"strconv"
)

type ListDeployController struct {
	*iris.Context
	org *myorganization.Organization
	dclist []mydatacenter.DataCenter
}

func (ldc *ListDeployController)getDcHost() (server []string, err error)  {
	server = make([]string, len(ldc.dclist))
	for i := 0; i < len(server); i++ {
		server[i] = ldc.dclist[i].Host + ":" + strconv.Itoa(int(ldc.dclist[i].Port))
	}
	return server, nil
}

func (ldc *ListDeployController)getDcName() (name []string, err error) {
	name = make([]string, len(ldc.dclist))
	for i := 0; i < len(name); i++ {
		name[i] = ldc.dclist[i].Name
	}
	return name, nil
}

func (ldc *ListDeployController) getPodList(dcId []int32, dcName []string, dcHost []string, orgId string) (list string, err error) {

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

		tmpdata[i].DcName = dcName[i]
		tmpdata[i].DcId = dcId[i]
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

	sessionIdClient := ldc.RequestHeader("Authorization")
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

	var dclist deploy.DcList
	log.Printf("%s\n", ldc.org.DcList)
	err = json.Unmarshal([]byte(ldc.org.DcList), &dclist)
	if err != nil {
		log.Printf("dclist=%s error=%s\n", dclist, err)
	}
	log.Printf("%s\n", dclist.Dclist)

	ldc.dclist = make([]mydatacenter.DataCenter, len(dclist.Dclist))
	for i := 0; i < len(dclist.Dclist); i++ {
		tmpDc, err := datacenter.GetDataCenterById(dclist.Dclist[i])
		if err != nil {
			log.Printf("Get Organization By orgId error: orgId=%s, error=%s\n", orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			json, _ := ye.EncodeJson()
			ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ldc.Write(json)
			return
		}
		ldc.dclist[i] = *tmpDc
	}

	orgName := ldc.org.Name

	ss := session.SessionStoreInstance()

	if ok, err := ss.ValidateOrgId(sessionIdClient, orgId); ok {
		server, err := ldc.getDcHost()
		if err != nil {
			log.Printf("Get Datacenter Host error: sessionId=%s, orgId=%s, err=%s\n", sessionIdClient, orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			json, _ := ye.EncodeJson()
			ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ldc.Write(json)
			return
		}

		name, err := ldc.getDcName()
		if err != nil {
			log.Printf("Get Datacenter Name error: sessionId=%s, orgId=%s, err=%s\n", sessionIdClient, orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			json, _ := ye.EncodeJson()
			ldc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ldc.Write(json)
			return
		}

		id := make([]int32, len(ldc.dclist))
		for i := 0; i < len(ldc.dclist); i++ {
			id[i] = ldc.dclist[i].Id
		}


		podlist, err := ldc.getPodList(id, name, server, orgName)
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
