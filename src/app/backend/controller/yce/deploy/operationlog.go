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
	SELECT_DEPLOYMENT = "SELECT id, name, actionType, actionVerb, actionUrl, actionAt, actionOp, dcList, success, reason, json, comment FROM deployment WHERE orgId=?"
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

	return
}

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

	//SELECT_DEPLOYMENT = "SELECT id, name, actionTypes, actionVerb, actionUrl, actionAt, actionOp, dcList, success, reason, json, comment WHERE orgId=?"

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

	mylog.Log.Infof("queryOperationLogMySQL successfully")
	return deployments
}

func (loc *ListOperationLogController) queryUserNameByUserId(userId int32) (name string) {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(SELECT_USER)
	if err != nil {
		mylog.Log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId).Scan(&name)
	if err != nil {
		mylog.Log.Errorf("queryOperationLogMySQL Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EMYSQL_QUERY, "")
		return
	}
	mylog.Log.Infof("queryUserNameByUserId successfully")
	return name
}

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

func (loc *ListOperationLogController) getOperationLog(deployments []mydeploy.Deployment) (opString string){
	opLog := make([]deploy.OperationLog, len(deployments))

	for index, dp := range deployments{
		userName := loc.queryUserNameByUserId(dp.ActionOp)

		dcIdListJSON := []byte(dp.DcList)
		mylog.Log.Debugf("dcIdListJSON=%s", dp.DcList)
		dcIdList := new(deploy.DcIdListType)
		err := json.Unmarshal(dcIdListJSON, dcIdList)
		if err != nil {
			mylog.Log.Errorf("getOperationLog Error: error=%s", err)
			loc.Ye = myerror.NewYceError(myerror.EJSON, "")
			return
		}



		mylog.Log.Debugf("dcIdList:%s", dcIdList)
		dcNameList := loc.queryDcNameByDcId(dcIdList.DcIdList)
		opLog[index].UserName = userName
		opLog[index].DcName = dcNameList
		opLog[index].Record = &dp
	}

	opJson, err := json.Marshal(opLog)
	if err != nil {
		mylog.Log.Errorf("getOperationLog Error: error=%s", err)
		loc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return
	}

	opString = string(opJson)
	mylog.Log.Infof("getOperationLog successfully")
	return opString
}

func (loc ListOperationLogController) Get() {
	orgId := loc.Param("orgId")

	sessionIdFromClient := loc.RequestHeader("Authorization")

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
	opString := loc.getOperationLog(dp)
	if loc.Ye != nil {
		loc.WriteBack()
		return
	}

	loc.Ye = myerror.NewYceError(myerror.EOK, opString)
	loc.WriteBack()

	mylog.Log.Infof("ListOperationLogController get over!")

	return

}
