package deploy

import (
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type LogsPodController struct {
	// must
	yce.Controller
	apiServer string
	k8sClient *client.Client

	// url param
	orgId string

	// json param
	params *LogsPodParam

	// pod
	pods    api.Pod
	podName string
}

// json from client
type LogsPodParam struct {
	UserId       string             `json:"userId"`
	DcIdList     []int32            `json:"dcIdList"`
	LogOption    *yceutils.LogOptionType 	`json:"logOption,omitempty"`
}


// parse params from json
func (lpc *LogsPodController) getParams() {
	log.Debugf("LogsPodController getParams: lpc.params=%p", lpc.params)
	log.Debugf("LogsPodController getParams: lpc.params.DcIdList=%p", &lpc.params.DcIdList)
	log.Debugf("LogsPodController getParams: lpc.params.logOption=%p", lpc.params.LogOption)
	err := lpc.ReadJSON(lpc.params)
	if err != nil {
		log.Errorf("LogsPodController getParams Error: error=%s", err)
		lpc.Ye = myerror.NewYceError(myerror.EYCE_LOGS_POD, "")
		return
	}
	log.Debugf("LogsPodController getParams successfully: dcId=%d, userId=%s", lpc.params.DcIdList[0], lpc.params.UserId)
	log.Debugf("LogsPodController getParams: LogOption=%v", lpc.params.LogOption)
}

// get dcId int32
func (lpc *LogsPodController) getDcId() int32 {
	//dcId, _ := strconv.Itoi(lpc.params.DcId)
	//return int32(dcId)
	if len(lpc.params.DcIdList) > 0 {
		return lpc.params.DcIdList[0]
	}else {
		log.Errorf("LogsPodController getDcId Error: len(DcIdList)=%d, err=no value in DcIdList, Index out of range", len(lpc.params.DcIdList))
		lpc.Ye = myerror.NewYceError(myerror.EOOM, "")
		return 0
	}
}

// getDatacenter by DcId
func (lpc *LogsPodController) getDatacenterByDcId(dcId int32) *mydatacenter.DataCenter {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		log.Errorf("LogsPodController getDatacenter QueryDataCenterById Error: dcId=%d, error=%s", dcId, err)
		lpc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return nil
	}

	log.Infof("LogsPodController getDatacenterByDcId successfully: name=%s, id=%d", dc.Name, dc.Id)
	return dc
}

// getApiServer
func (lpc *LogsPodController) getApiServer() {
	dcId := lpc.getDcId()
	if lpc.Ye != nil {
		return
	}

	datacenter := lpc.getDatacenterByDcId(dcId)
	if lpc.Ye != nil {
		return
	}

	host := datacenter.Host
	port := strconv.Itoa(int(datacenter.Port))
	apiServer := host + ":" + port

	lpc.apiServer = apiServer

	log.Infof("LogsPodController getApiServer successfully: apiServer=%s, dcId=%d", apiServer, dcId)

}

// create single k8s client by dcId
func (lpc *LogsPodController) createK8sClient() *client.Client {
	config := &restclient.Config{
		Host: lpc.apiServer,
	}

	c, err := client.New(config)
	if err != nil {
		log.Errorf("LogsPodController createK8sClient Error: apiServer=%s, error=%s", lpc.apiServer, err)
	}

	log.Debugf("LogsPodController createK8sClient successfully: apiServer=%s, k8sClient=%p", lpc.apiServer, c)
	return c
}

// get k8sclient
func (lpc *LogsPodController) getK8sClient() {
	lpc.k8sClient = lpc.createK8sClient()
	log.Infof("LogsPodController getK8sClient successfully: k8sClient=%p, apiServer=%s", lpc.k8sClient, lpc.apiServer)
}

// get OrgNameByOrgId
func (lpc *LogsPodController) getOrgNameByOrgId() string {
	organization := new(myorganization.Organization)

	orgId, _ := strconv.Atoi(lpc.orgId)
	organization.QueryOrganizationById(int32(orgId))
	log.Infof("LogsPodController getOrgNameByOrgId successfully: orgName=%s, orgId=%d", organization.Name, orgId)
	return organization.Name
}

// get Pod By podName
func (lpc *LogsPodController) getPodLogsByPodName() string {

	options := &api.PodLogOptions{
		Container:lpc.params.LogOption.Container,
		Follow: lpc.params.LogOption.Follow,
		Previous: lpc.params.LogOption.Previous,
		SinceSeconds:lpc.params.LogOption.SinceSeconds,
		SinceTime:lpc.params.LogOption.SinceTime,
		Timestamps:lpc.params.LogOption.Timestamps,
		TailLines:lpc.params.LogOption.TailLines,
		LimitBytes:lpc.params.LogOption.LimitBytes,
	}

	namespace := lpc.getOrgNameByOrgId()
	reader, err := lpc.k8sClient.Pods(namespace).GetLogs(lpc.podName, options).Stream()
	if err != nil {
		log.Errorf("LogsPodController GetLogs Error: podName=%s, error=%s", lpc.podName, err)
		lpc.Ye = myerror.NewYceError(myerror.EKUBE_LOGS_POD, "")
		return ""
	}
	defer reader.Close()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorf("LogsPodController ReadLogs Error: podName=%s, error=%s", lpc.podName, err)
		lpc.Ye = myerror.NewYceError(myerror.EKUBE_LOGS_POD, "")
		return ""
	}

	logs := string(b)
	log.Infof("LogsPodController getPodLosgByPodName successfully: len(bytes)=%d", len(b))

	return logs

}

// logs all
func (lpc *LogsPodController) logs() string {

	logs := lpc.getPodLogsByPodName()
	if lpc.Ye != nil {
		return ""
	}

	log.Infof("LogsPodController logs pod successfully")

	return logs
}

// main
func (lpc LogsPodController) Post() {

	lpc.params = new(LogsPodParam)
	lpc.params.LogOption = new(yceutils.LogOptionType)

	sessionIdFromClient := lpc.RequestHeader("Authorization")
	lpc.orgId = lpc.Param("orgId")
	lpc.podName = lpc.Param("podName")

	log.Debugf("LogsPodController Params: sessionId=%s, orgId=%s, podName=%s", sessionIdFromClient, lpc.orgId, lpc.podName)


	// Validate sessionId
	lpc.ValidateSession(sessionIdFromClient, lpc.orgId)
	if lpc.CheckError() {
		return
	}

	//Parse Param
	lpc.getParams()
	if lpc.CheckError() {
		return
	}

	// getApiServer
	lpc.getApiServer()
	if lpc.CheckError() {
		return
	}

	// getK8sClient
	lpc.getK8sClient()
	if lpc.CheckError() {
		return
	}

	// logs Pod
	logs := lpc.logs()
	if lpc.CheckError() {
		return
	}

	lpc.WriteOk(logs)
	log.Infoln("Logs Pod over!")
	return

}
