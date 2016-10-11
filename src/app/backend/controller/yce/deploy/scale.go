package deploy

import (
	"app/backend/model/yce/deploy"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"

	"app/backend/common/util/Placeholder"
	myerror "app/backend/common/yce/error"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
	mydeployment "app/backend/model/mysql/deployment"
	"github.com/kubernetes/kubernetes/pkg/util/json"
	"strconv"
)

type ScaleDeploymentController struct {
	yce.Controller
	k8sClient  *client.Client
	apiServer  string
	orgId      string
	userId     string
	dcId       string
	name       string
	s          *deploy.ScaleDeployment
	deployment *extensions.Deployment
}

// Scale directly
func (sdc *ScaleDeploymentController) scaleSimple() {
	sdc.deployment.Spec.Replicas = sdc.s.NewSize
	_, err := sdc.k8sClient.Extensions().Deployments(sdc.deployment.Namespace).Update(sdc.deployment)
	if err != nil {
		log.Errorf("ScaleDeployment ScaleSimple Error: name=%s, namespace=%s, newsize=%d", sdc.deployment.Name, sdc.deployment.Namespace, sdc.s.NewSize)
		sdc.Ye = myerror.NewYceError(myerror.EKUBE_SCALE_DEPLOYMENT, "")
		return
	}

	log.Infof("ScaleDeployment ScaleSimple Successfully")
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
		log.Errorf("CreateMysqlDeployment Error: actionUrl=%s, actionOp=%d, dcList=%s, err=%s",
			actionUrl, actionOp, dcList, err)
		sdc.Ye = myerror.NewYceError(myerror.EMYSQL_INSERT, "")
		return err
	}

	log.Infof("ScaleDeployment CreateMysqlDeployment successfully: actionUrl=%s, actionOp=%d, dcList=%s",
		actionUrl, actionOp, dcList)
	return nil
}

func (sdc ScaleDeploymentController) Post() {
	sdc.orgId = sdc.Param("orgId")
	sdc.name = sdc.Param("deploymentName")

	sessionIdFromClient := sdc.RequestHeader("Authorization")

	log.Debugf("ScaleDeploymentController Params: sessionId=%s, orgId=%s, name=%s", sessionIdFromClient, sdc.orgId, sdc.name)

	// Validate the session
	sdc.ValidateSession(sessionIdFromClient, sdc.orgId)
	if sdc.CheckError() {
		return
	}

	// ScaleDeployment Params
	sdc.s = new(deploy.ScaleDeployment)
	err := sdc.ReadJSON(sdc.s)
	if err != nil {
		log.Errorf("ScaleDeployController ReadJSON Error: error=%s", err)
		sdc.Ye = myerror.NewYceError(myerror.EJSON, "")
	}

	if sdc.CheckError() {
		return
	}

	// Get DcIdList
	if len(sdc.s.DcIdList) == 0 {
		log.Errorln("Empty DcIdList!")
		sdc.Ye = myerror.NewYceError(myerror.EINVALID_PARAM, "")
	}

	if sdc.CheckError() {
		return
	}
	dcId := sdc.s.DcIdList[0]
	sdc.dcId = strconv.Itoa(int(dcId))

	// Get ApiServer
	sdc.apiServer, sdc.Ye = yceutils.GetApiServerByDcId(dcId)
	if sdc.CheckError() {
		return
	}

	// Create K8sClient
	sdc.k8sClient, sdc.Ye = yceutils.CreateK8sClient(sdc.apiServer)
	if sdc.CheckError() {
		return
	}

	// Get Namespace
	namespace, ye := yceutils.GetOrgNameByOrgId(sdc.orgId)
	if ye != nil {
		sdc.Ye = ye
	}
	if sdc.CheckError() {
		return
	}

	// Get Deployment
	sdc.deployment, sdc.Ye = yceutils.GetDeploymentByNameAndNamespace(sdc.k8sClient, sdc.name, namespace)
	if sdc.CheckError() {
		return
	}

	//scale the deployment
	sdc.scaleSimple()
	if sdc.CheckError() {
		return
	}

	// prepare for create mysql deployment
	dd, _ := json.Marshal(sdc.deployment)
	dcIdList, ye := yceutils.EncodeDcIdList(sdc.s.DcIdList)
	if ye != nil {
		sdc.Ye = ye
	}
	if sdc.CheckError() {
		return
	}

	oId, _ := strconv.Atoi(sdc.orgId)

	// create mysql deployment
	sdc.createMysqlDeployment(true, sdc.name, string(dd), sdc.s.Comments, dcIdList, sdc.s.UserId, int32(oId))
	if sdc.CheckError() {
		return
	}

	// success
	sdc.WriteOk("")
	log.Infoln("ScaleDeploymentController over!")
	return
}
