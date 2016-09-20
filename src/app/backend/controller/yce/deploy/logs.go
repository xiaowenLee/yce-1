package deploy

import (
	myerror "app/backend/common/yce/error"
	"k8s.io/kubernetes/pkg/api"
	client "k8s.io/kubernetes/pkg/client/unversioned"
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


// logs all
func (lpc *LogsPodController) logs() string {

	logs, ye := yceutils.GetPodLogsByPodName(lpc.k8sClient, lpc.params.LogOption, lpc.podName, lpc.orgId)
	if ye != nil {
		lpc.Ye = ye
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

	// Get ApiServer
	dcId := lpc.params.DcIdList[0]
	lpc.apiServer, lpc.Ye = yceutils.GetApiServerByDcId(dcId)
	if lpc.CheckError() {
		return
	}

	// Get K8sClient
	lpc.k8sClient, lpc.Ye = yceutils.CreateK8sClient(lpc.apiServer)
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
