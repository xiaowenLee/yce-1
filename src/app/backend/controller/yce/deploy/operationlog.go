package deploy

import (
	myerror "app/backend/common/yce/error"
	mydeploy "app/backend/model/mysql/deployment"
	myuser "app/backend/model/mysql/user"
	"app/backend/common/util/mysql"
	"app/backend/model/yce/deploy"
	"strconv"
	"encoding/json"
	yce "app/backend/controller/yce"
	yceutils "app/backend/controller/yce/utils"
)

type ListOperationLogController struct {
	yce.Controller
}

// Query Deployments according to orgId
func (loc *ListOperationLogController) queryOperationLogMySQL(orgId int32) (deployments []mydeploy.Deployment) {

	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(SELECT_DEPLOYMENT)
	if err != nil {
		log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(orgId)
	if err != nil {
		log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	for rows.Next() {
		dp := new(mydeploy.Deployment)
		var comment []byte
		var jsonFile []byte
		err := rows.Scan(&dp.Id, &dp.Name, &dp.ActionType, &dp.ActionVerb, &dp.ActionUrl, &dp.ActionAt, &dp.ActionOp, &dp.DcList, &dp.Success, &dp.Reason, &jsonFile, &comment);
		if err != nil {
			log.Errorf("queryOperationLogMySQL Error: error=%s", err)
			loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
			return
		}

		dp.Json = string(jsonFile)
		dp.Comment = string(comment)

		deployments = append(deployments, *dp)
		log.Debugf("query result: id=%d, name=%s", dp.Id, dp.Name)
	}
	log.Infof("queryOperationLogMySQL successfully, totally %d deployments", len(deployments))
	return deployments
}


// getOperationLogList
func (loc *ListOperationLogController) getOperationLog(deployment mydeploy.Deployment) *deploy.OperationLogType {
	opLog := new(deploy.OperationLogType)
	userName := myuser.QueryUserNameByUserId(deployment.ActionOp)

	dcIdListJSON := []byte(deployment.DcList)
	log.Debugf("dcIdListJSON=%s", deployment.DcList)
	dcIdList := new(deploy.DcIdListType)
	err := json.Unmarshal(dcIdListJSON, dcIdList)
	if err != nil {
		log.Errorf("getOperationLog Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return nil

	}

	dcNameList, ye := yceutils.GetDcNameListByDcIdList(dcIdList.DcIdList)
	if ye != nil {
		loc.Ye = ye
		return nil
	}
	opLog.UserName = userName
	opLog.DcName = dcNameList
	opLog.Record = &deployment

	log.Infof("getOperationLog userName=%s, dcName=%v deploymentName=%s", opLog.UserName, opLog.DcName, opLog.Record.Name)

	return opLog
}

func (loc *ListOperationLogController) getOperationLogList(deployments []mydeploy.Deployment) string {
	opLogList := new(deploy.OperationLogList)
	opLogList.OperationLog = make([]deploy.OperationLogType, 0)


	for _, dp := range deployments{
		opLog := loc.getOperationLog(dp)
		opLogList.OperationLog = append(opLogList.OperationLog, *opLog)

		log.Infof("ListOperationController getOperation: name=%s, userName=%s, len(dcName):%d", dp.Name, opLog.UserName, len(opLog.DcName))
	}

	opLogListJson, _ := json.Marshal(opLogList)
	opLogListString := string(opLogListJson)

	log.Infof("ListOperationController getOperationLogList over: len(deployment)=%d", len(opLogList.OperationLog))
	return opLogListString

}

func (loc ListOperationLogController) Get() {
	orgId := loc.Param("orgId")

	sessionIdFromClient := loc.RequestHeader("Authorization")
	log.Debugf("ListOperationLogController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)


	// Validate sessionId with orgId
	loc.ValidateSession(sessionIdFromClient, orgId)
	if loc.CheckError() {
		return
	}

	// Get OperationLogMySQL
	oId, _ := strconv.Atoi(orgId)
	dp := loc.queryOperationLogMySQL(int32(oId))

	// Get OperationLog
	opString := loc.getOperationLogList(dp)
	if loc.CheckError() {
		return
	}

	loc.WriteOk(opString)
	log.Infof("ListOperationLogController get over!")

	return
}
