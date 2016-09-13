package deploy

import (
	"github.com/kataras/iris"
	myerror "app/backend/common/yce/error"
	mylog "app/backend/common/util/log"
	mydeploy "app/backend/model/mysql/deployment"
	"app/backend/common/util/session"
	"app/backend/common/util/mysql"
	"app/backend/model/yce/deploy"
	"strconv"
	"encoding/json"
)

type ListOperationLogController struct {
	*iris.Context
	Ye *myerror.YceError
}

const (
	SELECT_DEPLOYMENT = "SELECT id, name, actionType, actionVerb, actionUrl, actionAt, actionOp, dcList, success, reason, json, comment FROM deployment WHERE orgId=? ORDER BY id DESC LIMIT 30"
	SELECT_USER = "SELECT name FROM user WHERE id=?"
	SELECT_DATACENTER = "SELECT name FROM datacenter WHERE id=?"
)


func (loc *ListOperationLogController) WriteBack() {
	loc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("ListOperationLogController Response YceError: controller=%p, code=%d, note=%s", loc, loc.Ye.Code, myerror.Errors[loc.Ye.Code].LogMsg)
	loc.Write(loc.Ye.String())
}

// Validate Session
func (loc *ListOperationLogController) validateSession(sessionId, orgId string) {
	// Validate the session
	ss := session.SessionStoreInstance()

	ok, err := ss.ValidateOrgId(sessionId, orgId)
	if err != nil {
		mylog.Log.Errorf("Validate Session error: sessionId=%s, error=%s", sessionId, err)
		loc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	// Session invalide
	if !ok {
		mylog.Log.Errorf("Validate Session failed: sessionId=%s, error=%s", sessionId, err)
		loc.Ye = myerror.NewYceError(myerror.EYCE_SESSION, "")
		return
	}

	mylog.Log.Infof("ListOperationLogController ValidateSession success")
	return
}
// Query Deployments according to orgId
func (loc *ListOperationLogController) queryOperationLogMySQL(orgId int32) (deployments []mydeploy.Deployment) {

	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(SELECT_DEPLOYMENT)
	if err != nil {
		mylog.Log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(orgId)
	if err != nil {
		mylog.Log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}

	for rows.Next() {
		dp := new(mydeploy.Deployment)
		var comment []byte
		var jsonFile []byte
		err := rows.Scan(&dp.Id, &dp.Name, &dp.ActionType, &dp.ActionVerb, &dp.ActionUrl, &dp.ActionAt, &dp.ActionOp, &dp.DcList, &dp.Success, &dp.Reason, &jsonFile, &comment);
		if err != nil {
			mylog.Log.Errorf("queryOperationLogMySQL Error: error=%s", err)
			loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
			return
		}

		dp.Json = string(jsonFile)
		dp.Comment = string(comment)

		deployments = append(deployments, *dp)
		mylog.Log.Debugf("query result: id=%d, name=%s", dp.Id, dp.Name)
	}
	mylog.Log.Infof("queryOperationLogMySQL successfully, totally %d deployments", len(deployments))
	return deployments
}

// Query UserName by UserId
func (loc *ListOperationLogController) queryUserNameByUserId(userId int32) (name string) {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(SELECT_USER)
	if err != nil {
		mylog.Log.Errorf("queryUserNameByUserId Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId).Scan(&name)
	if err != nil {
		mylog.Log.Errorf("queryUserNameByUserId Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	mylog.Log.Infof("queryUserNameByUserId successfully")
	return name
}

// Query DcName By DcId
func (loc *ListOperationLogController) queryDcNameByDcId(dcIdList []int32) (dcNameList []string) {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(SELECT_DATACENTER)
	if err != nil {
		mylog.Log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	defer stmt.Close()
	dcNameList = make([]string, len(dcIdList))
	for index, dcId := range dcIdList {
		err = stmt.QueryRow(dcId).Scan(&dcNameList[index])
		if err != nil {
			mylog.Log.Errorf("queryOperationLogMySQL Error: error=%s", err)
			loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
			return
		}
	}

	mylog.Log.Infof("queryDcNameByDcId successfully")
	return dcNameList
}

// getOperationLogList

func (loc *ListOperationLogController) getOperationLog(deployment mydeploy.Deployment) *deploy.OperationLogType {
	opLog := new(deploy.OperationLogType)
	userName := loc.queryUserNameByUserId(deployment.ActionOp)

	dcIdListJSON := []byte(deployment.DcList)
	mylog.Log.Debugf("dcIdListJSON=%s", deployment.DcList)
	dcIdList := new(deploy.DcIdListType)
	err := json.Unmarshal(dcIdListJSON, dcIdList)
	if err != nil {
		mylog.Log.Errorf("getOperationLog Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return nil

	}

	dcNameList := loc.queryDcNameByDcId(dcIdList.DcIdList)
	opLog.UserName = userName
	opLog.DcName = dcNameList
	opLog.Record = &deployment

	mylog.Log.Infof("getOperationLog userName=%s, dcName=%v deploymentName=%s", opLog.UserName, opLog.DcName, opLog.Record.Name)

	return opLog
}

func (loc *ListOperationLogController) getOperationLogList(deployments []mydeploy.Deployment) string {
	opLogList := new(deploy.OperationLogList)
	opLogList.OperationLog = make([]deploy.OperationLogType, 0)


	for _, dp := range deployments{
		opLog := loc.getOperationLog(dp)
		opLogList.OperationLog = append(opLogList.OperationLog, *opLog)

		mylog.Log.Infof("ListOperationController getOperation: name=%s, userName=%s, len(dcName):%d", dp.Name, opLog.UserName, len(opLog.DcName))
	}

	opLogListJson, _ := json.Marshal(opLogList)
	opLogListString := string(opLogListJson)

	mylog.Log.Infof("ListOperationController getOperationLogList over: len(deployment)=%d", len(opLogList.OperationLog))
	return opLogListString

}

// getOperationLog list
/*
func (loc *ListOperationLogController) getOperationLogList(deployments []mydeploy.Deployment) (opString string){

	opLog := loc.getOperationLog(deployments)
	opJson, err := json.Marshal(opLog)
	if err != nil {
		mylog.Log.Errorf("getOperationLog Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	opString = string(opJson)
	mylog.Log.Infof("ListOperationLogController getOperationLog successfully")
	return opString
}
*/
func (loc ListOperationLogController) Get() {
	orgId := loc.Param("orgId")

	sessionIdFromClient := loc.RequestHeader("Authorization")
	mylog.Log.Debugf("ListOperationLogController Params: sessionId=%s, orgId=%s", sessionIdFromClient, orgId)


	// Validate sessionId with orgId
	loc.validateSession(sessionIdFromClient, orgId)
	if loc.Ye != nil {
		loc.WriteBack()
		return
	}

	// Get OperationLogMySQL
	oId, _ := strconv.Atoi(orgId)
	dp := loc.queryOperationLogMySQL(int32(oId))



	// Get OperationLog
	opString := loc.getOperationLogList(dp)
	if loc.Ye != nil {
		loc.WriteBack()
		return
	}

	loc.Ye = myerror.NewYceError(myerror.EOK, opString)
	loc.WriteBack()

	mylog.Log.Infof("ListOperationLogController get over!")

	return

}
