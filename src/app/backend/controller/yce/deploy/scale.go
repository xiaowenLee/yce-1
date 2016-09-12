package deploy

import (
	"app/backend/model/yce/deploy"
	"app/backend/common/util/session"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"github.com/kataras/iris"

	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	mydeployment "app/backend/model/mysql/deployment"
	myoption "app/backend/model/mysql/option"
	"app/backend/common/yce/organization"
	"strconv"
	"app/backend/common/util/Placeholder"
	"github.com/kubernetes/kubernetes/pkg/util/json"
)

const (
	SCALE_ACTION_TYPE                = myoption.SCALING
	SCALE_ACTION_VERBE               = "POST"
	SCALE_ACTION_URL                 = "/api/v1/organizations/<orgId>/datacenters/<dcId>/deployments/<name>/scale"
)

type ScaleDeploymentController struct{
	*iris.Context
	k8sClient *client.Client
	apiServer string
	Ye *myerror.YceError
	orgId string
	userId string
	dcId string
	name string
	s *deploy.ScaleDeployment
	deployment extensions.Deployment
}

func (sdc *ScaleDeploymentController) WriteBack() {
	sdc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("ScaleDeploymentController Response YceError: controller=%p, code=%d, note=%s", sdc, sdc.Ye.Code, myerror.Errors[sdc.Ye.Code].LogMsg)
	sdc.Write(sdc.Ye.String())
}

// validate sessionId from Client
func (sdc *ScaleDeploymentController) valideSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		sdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		sdc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	mylog.Log.Infof("ScaleDeployment sessionId successfully: sessionId=%s, orgId=%d", sessionId, orgId)
	return
}

// get ApiServer And K8sClient By DcId
func (sdc *ScaleDeploymentController) getApiServerAndK8sClientByDcId() {
	dc := new(mydatacenter.DataCenter)

	//TODO: find a better way
	dcId := sdc.s.DcIdList[0]
	sdc.dcId = strconv.Itoa(int(dcId))

	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("ScaleDeployment getApiServerById QueryDataCenterById error: dcId=%d, error=%s", dcId, err)
		sdc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	host := dc.Host
	port := strconv.Itoa(int(dc.Port))
	sdc.apiServer = host + ":" + port

	config := &restclient.Config{
		Host: sdc.apiServer,
	}

	c, err := client.New(config)
	if err != nil {
		mylog.Log.Errorf("ScaleDeployment create K8sClient error: apiServer=%s, error=%s", sdc.apiServer, err)
		sdc.Ye = myerror.NewYceError(myerror.EKUBE_CLIENT, "")
		return
	}

	sdc.k8sClient = c
	mylog.Log.Infof("ScaleDeployment GetApiServerAndK8sClientByDcID over: apiServer=%s, k8sClient=%p", sdc.apiServer, sdc.k8sClient)

}

// get Deployment By Name
func (sdc *ScaleDeploymentController) getDeploymentByName() {
	//get Organization by OrgId
	org, err := organization.GetOrganizationById(sdc.orgId)
	if err != nil {
		mylog.Log.Errorf("ScaleDeployment getDatacentersByOrgId Error: orgId=%s, error=%s", sdc.orgId, err)
		sdc.Ye = myerror.NewYceError(myerror.EYCE_ORGTODC, "")
		return

	}


	// get Deployments by deployment's name and namespace
	namespace := org.Name
	dp, err := sdc.k8sClient.Extensions().Deployments(namespace).Get(sdc.name)
	if err != nil {
		mylog.Log.Errorf("ScaleDeployment getDeployByName Error: apiServer=%s, namespace=%s, deployment-name=%s, err=%s\n",
			sdc.apiServer, namespace, sdc.name, err)
		sdc.Ye = myerror.NewYceError(myerror.EKUBE_GET_DEPLOYMENT, "")
		return
	}

	sdc.deployment = *dp

	mylog.Log.Infof("ScaleDeployment GetDeploymentByName over: apiServer=%s, namespace=%s, name=%s, deployment=%p\n",
		sdc.apiServer, namespace, sdc.name, dp)
}


// Scale directly
func (sdc *ScaleDeploymentController) ScaleSimple() {
	sdc.deployment.Spec.Replicas = sdc.s.NewSize
	_, err := sdc.k8sClient.Extensions().Deployments(sdc.deployment.Namespace).Update(&sdc.deployment)
	if err != nil {
		mylog.Log.Errorf("ScaleDeployment ScaleSimple Error: name=%s, namespace=%s, newsize=%d", sdc.deployment.Name, sdc.deployment.Namespace, sdc.s.NewSize)
	}

	mylog.Log.Infof("ScaleDeployment ScaleSimple Successfully")
}

// create a deployment record
func (sdc *ScaleDeploymentController) createMysqlDeployment(success bool, name, json, reason, dcList string, userId, orgId int32) error {
	//TODO: actionUrl not complete
	uph := placeholder.NewPlaceHolder(SCALE_ACTION_URL)
	orgIdString := strconv.Itoa(int(orgId))
	actionUrl := uph.Replace("<orgId>", orgIdString, "<dcId>", sdc.dcId, "<name>", name)
	actionOp := userId
	dp := mydeployment.NewDeployment(name, SCALE_ACTION_VERBE, actionUrl, dcList, reason, json, "Rolilng Update a Deployment", int32(SCALE_ACTION_TYPE), actionOp, int32(1), orgId)
	err := dp.InsertDeployment()
	if err != nil {
		mylog.Log.Errorf("CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		sdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	mylog.Log.Infof("ScaleDeployment CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

// encode DcIdList
func (sdc *ScaleDeploymentController) encodeDcIdList() string {
	dcIdList := &deploy.DcIdListType{
		DcIdList:sdc.s.DcIdList,
	}
	data, _ := json.Marshal(dcIdList)

	mylog.Log.Infof("ScaleDeployController encodeDcIdList: dcIdList=%s", string(data))
	return string(data)
}

func (sdc ScaleDeploymentController) Post() {
	sdc.orgId = sdc.Param("orgId")
	sdc.name = sdc.Param("deploymentName")

	//validate the session
	sessionIdFromClient := sdc.RequestHeader("Authorization")
	sdc.valideSession(sessionIdFromClient, sdc.orgId)
	if sdc.Ye != nil {
		sdc.WriteBack()
		return
	}

	// ScaleDeployment Params
	sdc.s = new(deploy.ScaleDeployment)
	sdc.ReadJSON(sdc.s)

	//get ApiServer and K8sClient
	sdc.getApiServerAndK8sClientByDcId()
	if sdc.Ye != nil {
		sdc.WriteBack()
		return
	}

	//get Deployment
	sdc.getDeploymentByName()
	if sdc.Ye != nil {
		sdc.WriteBack()
		return
	}

	//scale the deployment
	sdc.ScaleSimple()
	if sdc.Ye != nil {
		sdc.WriteBack()
		return
	}

	// prepare for create mysql deployment
	dd, _ := json.Marshal(sdc.deployment)
	dcIdList := sdc.encodeDcIdList()
	oId, _ := strconv.Atoi(sdc.orgId)

	// create mysql deployment
	sdc.createMysqlDeployment(true, sdc.name, string(dd), sdc.s.Comments, dcIdList, sdc.s.UserId, int32(oId))
	if sdc.Ye != nil {
		sdc.WriteBack()
		return
	}

	// success
	sdc.Ye = myerror.NewYceError(myerror.EOK, "")
	sdc.WriteBack()
	mylog.Log.Infoln("ScaleDeploymentController over!")
	return

}