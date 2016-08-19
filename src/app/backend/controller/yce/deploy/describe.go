package deploy

/*
import (
	"app/backend/common/util/session"
	"app/backend/common/yce/datacenter"
	myerror "app/backend/common/yce/error"
	"app/backend/common/yce/organization"
	mydatacenter "app/backend/model/mysql/datacenter"
	"app/backend/model/yce/deploy"
	"encoding/json"
	"github.com/kataras/iris"
	"log"
)

type DescribeDeployController struct {
	*iris.Context
	list *ListDeployController
	//TODO: 可能会有ReplicaSet、Service等的成员
}

func (ddc DescribeDeployController) Get() {

	ddc.list = new(ListDeployController)

	sessionIdClient := ddc.RequestHeader("sessionId")
	orgId := ddc.Param("orgId")
	//userId := ldc.Param("userId")

	tmpOrg, err := organization.GetOrganizationById(orgId)
	if err != nil {
		log.Printf("Get Organization By orgId error: orgId=%s, error=%s\n", orgId, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		ddc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ddc.Write(errJson)
		return
	}
	ddc.list.org = tmpOrg

	var dc deploy.DcList
	err = json.Unmarshal([]byte(ddc.list.org.DcList), &dc)
	if err != nil {
		log.Printf("DecodeJSON error: dc=%s error=%s\n", dc, err)
	}

	ddc.list.dclist = make([]mydatacenter.DataCenter, len(dc.DataCenter))
	for i := 0; i < len(dc.DataCenter); i++ {
		tmpDc, err := datacenter.GetDataCenterById(dc.DataCenter[i])
		if err != nil {
			log.Printf("Get Organization By orgId error: orgId=%s, error=%s\n", orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			errJson, _ := ye.EncodeJson()
			ddc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ddc.Write(errJson)
			return
		}
		ddc.list.dclist[i] = *tmpDc
	}

	orgName := ddc.list.org.Name

	ss := session.SessionStoreInstance()

	if ok, err := ss.ValidateOrgId(sessionIdClient, orgId); ok {
		server, err := ddc.list.getDcHost()
		if err != nil {
			log.Printf("Get Datacenter Host error: sessionId=%s, orgId=%s, err=%s\n", sessionIdClient, orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			errJson, _ := ye.EncodeJson()
			ddc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ddc.Write(errJson)
			return
		}

		name, err := ddc.list.getDcName()
		if err != nil {
			log.Printf("Get Datacenter Name error: sessionId=%s, orgId=%s, err=%s\n", sessionIdClient, orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			errJson, _ := ye.EncodeJson()
			ddc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ddc.Write(errJson)
			return
		}

		id := make([]int32, len(ddc.list.dclist))
		for i := 0; i < len(ddc.list.dclist); i++ {
			id[i] = ddc.list.dclist[i].Id
		}

		podlist, err := ddc.list.getPodList(id, name, server, orgName)
		if err != nil {
			log.Printf("Get Podlist error: sessionId=%s, orgId=%s, error=%s\n", sessionIdClient, orgId, err)
			ye := myerror.NewYceError(1, "ERR", "请求失败")
			errJson, _ := ye.EncodeJson()
			ddc.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ddc.Write(errJson)
			return
		}

		ye := myerror.NewYceError(0, "OK", podlist)
		errJson, _ := ye.EncodeJson()

		ddc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ddc.Write(errJson)
		return
	} else {
		log.Printf("Validate Session error: sessionId=%s, error=%s\n", sessionIdClient, err)
		ye := myerror.NewYceError(1, "ERR", "请求失败")
		errJson, _ := ye.EncodeJson()
		ddc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ddc.Write(errJson)
		return
	}

}
*/
