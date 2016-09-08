package deploy

import (
	mylog "app/backend/common/util/log"
	"app/backend/common/util/session"
	myerror "app/backend/common/yce/error"
	mydatacenter "app/backend/model/mysql/datacenter"
	myorganization "app/backend/model/mysql/organization"
	"github.com/kataras/iris"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api"
	unver "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"strconv"
)

type LogsPodController struct {
	// must
	*iris.Context
	apiServer string
	k8sClient *client.Client
	Ye        *myerror.YceError

	// url param
	orgId string

	// json param
	params *LogsPodParam

	// pod
	pods    api.Pod
	podName string
}

type LogOptionType struct {
	Container    string      `json:"container,omitempty"`    //暂时不做
	Follow       bool        `json:"follow,omitempty"`       //false 暂时不做, 页面开关,默认为关闭
	Previous     bool        `json:"previous,omitempty"`     //暂时不做
	SinceSeconds *int64      `json:"sinceSeconds,omitempty"` //暂时不做
	SinceTime    *unver.Time `json:"sinceTime,omitempty"`    //暂时不做
	Timestamps   bool        `json:"timeStamps,omitempty"`    //true, 时间戳,默认打开
	TailLines    *int64      `json:"tailLines,omitempty"`    //用户设定
	LimitBytes   *int64      `json:"limitBytes,omitempty"`   //暂时不做
}

// json from client
type LogsPodParam struct {
	UserId       string             `json:"userId"`
	DcIdList     []int32            `json:"dcIdList"`
	LogOption    *LogOptionType 	`json:"logOption,omitempty"`
}

// write response
func (lpc *LogsPodController) WriteBack() {
	lpc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("LogsPodController Response YceError: controller=%p, code=%d, note=%s", lpc, lpc.Ye.Code, myerror.Errors[lpc.Ye.Code].LogMsg)
	lpc.Write(lpc.Ye.String())
}

// validate Session with orgId
func (lpc *LogsPodController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("LogsPodController Validate Session error: sessionId=%s, error=%s", sessionId, err)
		lpc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("LogsPodController Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		lpc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	mylog.Log.Infof("LogsPodController validate sessionId successfully: sessionId=%s, orgId=%s", sessionId, orgId)
	return
}

// parse params from json
func (lpc *LogsPodController) getParams() {
	mylog.Log.Debugf("LogsPodController getParams: lpc.params=%p", lpc.params)
	mylog.Log.Debugf("LogsPodController getParams: lpc.params.DcIdList=%p", &lpc.params.DcIdList)
	mylog.Log.Debugf("LogsPodController getParams: lpc.params.logOption=%p", lpc.params.LogOption)
	err := lpc.ReadJSON(lpc.params)
	if err != nil {
		mylog.Log.Errorf("LogsPodController getParams Error: error=%s", err)
		lpc.Ye = myerror.NewYceError(myerror.EYCE_LOGS_POD, "")
		return
	}
	mylog.Log.Debugf("LogsPodController getParams successfully: dcId=%d, userId=%s", lpc.params.DcIdList[0], lpc.params.UserId)
	mylog.Log.Debugf("LogsPodController getParams: LogOption=%v", lpc.params.LogOption)
}

// get dcId int32
func (lpc *LogsPodController) getDcId() int32 {
	//dcId, _ := strconv.Itoi(lpc.params.DcId)
	//return int32(dcId)
	return lpc.params.DcIdList[0]
}

// getDatacenter by DcId
func (lpc *LogsPodController) getDatacenterByDcId(dcId int32) *mydatacenter.DataCenter {
	dc := new(mydatacenter.DataCenter)
	err := dc.QueryDataCenterById(dcId)
	if err != nil {
		mylog.Log.Errorf("LogsPodController getDatacenter QueryDataCenterById Error: dcId=%d, error=%s", dcId, err)
		lpc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return nil
	}

	mylog.Log.Infof("LogsPodController getDatacenterByDcId successfully: name=%s, id=%d", dc.Name, dc.Id)
	return dc
}

// getApiServer
func (lpc *LogsPodController) getApiServer() {
	dcId := lpc.getDcId()

	datacenter := lpc.getDatacenterByDcId(dcId)

	host := datacenter.Host
	port := strconv.Itoa(int(datacenter.Port))
	apiServer := host + ":" + port

	lpc.apiServer = apiServer

	mylog.Log.Infof("LogsPodController getApiServer successfully: apiServer=%s, dcId=%d", apiServer, dcId)

}

// create single k8s client by dcId
func (lpc *LogsPodController) createK8sClient() *client.Client {
	config := &restclient.Config{
		Host: lpc.apiServer,
	}

	c, err := client.New(config)
	if err != nil {
		mylog.Log.Errorf("LogsPodController createK8sClient Error: apiServer=%s, error=%s", lpc.apiServer, err)
	}

	mylog.Log.Debugf("LogsPodController createK8sClient successfully: apiServer=%s, k8sClient=%p", lpc.apiServer, c)
	return c
}

// get k8sclient
func (lpc *LogsPodController) getK8sClient() {
	lpc.k8sClient = lpc.createK8sClient()
	mylog.Log.Infof("LogsPodController getK8sClient successfully: k8sClient=%p, apiServer=%s", lpc.k8sClient, lpc.apiServer)
}

// get OrgNameByOrgId
func (lpc *LogsPodController) getOrgNameByOrgId() string {
	organization := new(myorganization.Organization)

	orgId, _ := strconv.Atoi(lpc.orgId)
	organization.QueryOrganizationById(int32(orgId))
	mylog.Log.Infof("LogsPodController getOrgNameByOrgId successfully: orgName=%s, orgId=%d", organization.Name, orgId)
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
		mylog.Log.Errorf("LogsPodController GetLogs Error: podName=%s, error=%s", lpc.podName, err)
		lpc.Ye = myerror.NewYceError(myerror.EKUBE_LOGS_POD, "")
		return ""
	}
	defer reader.Close()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		mylog.Log.Errorf("LogsPodController ReadLogs Error: podName=%s, error=%s", lpc.podName, err)
		lpc.Ye = myerror.NewYceError(myerror.EKUBE_LOGS_POD, "")
		return ""
	}

	logs := string(b)
	mylog.Log.Infof("LogsPodController getPodLosgByPodName successfully: len(bytes)=%d", len(b))

	return logs

}

// logs all
func (lpc *LogsPodController) logs() string {

	logs := lpc.getPodLogsByPodName()

	mylog.Log.Infof("LogsPodController logs pod successfully")

	return logs
}

// main
func (lpc LogsPodController) Get() {

	lpc.params = new(LogsPodParam)
	lpc.params.LogOption = new(LogOptionType)

	sessionIdFromClient := lpc.RequestHeader("Authorization")
	lpc.orgId = lpc.Param("orgId")
	lpc.podName = lpc.Param("podName")

	// validate sessionId
	lpc.validateSession(sessionIdFromClient, lpc.orgId)
	if lpc.Ye != nil {
		lpc.WriteBack()
		return
	}

	//Parse Param
	lpc.getParams()
	if lpc.Ye != nil {
		lpc.WriteBack()
		return
	}

	// getApiServer
	lpc.getApiServer()
	if lpc.Ye != nil {
		lpc.WriteBack()
		return
	}

	// getK8sClient
	lpc.getK8sClient()

	// logs Pod
	logs := lpc.logs()
	if lpc.Ye != nil {
		lpc.WriteBack()
		return
	}

	lpc.Ye = myerror.NewYceError(myerror.EOK, logs)
	lpc.WriteBack()

	// TODO: 成功写回
	mylog.Log.Infoln("Logs Pod over!")
	return

}
