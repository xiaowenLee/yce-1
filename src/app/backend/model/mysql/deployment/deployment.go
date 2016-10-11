package deployment

import (
	mylog "app/backend/common/util/log"
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
)

var log = mylog.Log

const (
	DEPLOYMENT_SELECT = "SELECT id, name, actionType, actionVerb, actionUrl, " +
		"actionAt, actionOp, dcList, success, reason, json, comment, orgId " +
		"FROM deployment where id=?"

	DEPLOYMENT_BYNAME = "SELECT id, name, actionType, actionVerb, actionUrl, " +
		"actionAt, actionOp, dcList, success, reason, json, comment, orgId " +
		"FROM deployment where name=? ORDER BY id DESC LIMIT 30"

	DEPLOYMENT_INSERT = "INSERT INTO deployment(name, actionType, actionVerb, actionUrl, " +
		"actionAt, actionOp, dcList, success, reason, json, comment, orgId) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	DEPLOYMENT_ACTIONTYPE_STAT = "SELECT id, name, actionType, actionVerb, actionUrl, " +
		"actionAt, actionOp, dcList, success, reason, json, comment, orgId " +
		"FROM deployment where actionType=?"
	VALID   = 1
	INVALID = 0
>>>>>>> a7a204ff83126d54dc9fa22d5c035ecd20b87340
)

type Deployment struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	ActionType int32  `json:"actionType"`
	ActionVerb string `json:"actionVerb"`
	ActionUrl  string `json:"actionUrl"`
	ActionAt   string `json:"actionAt"`
	ActionOp   int32  `json:"actionOp"`
	DcList     string `json:"dcList"`
	Success    int32  `json:"success"`
	Reason     string `json:"reason",omitempty`
	Json       string `json:"json"`
	Comment    string `json:"comment,omitempty"`
	OrgId      int32  `json:"orgId"`
}

type Statistics map[int32]int32

type OperationStat map[string]Statistics // day-->op-->total

func NewDeployment(name, actionVerb, actionUrl, dcList, reason, json, comment string, actionType, actionOp, success int32, orgId int32) *Deployment {
	return &Deployment{
		Name:       name,
		ActionType: actionType,
		ActionVerb: actionVerb,
		ActionUrl:  actionUrl,
		ActionAt:   localtime.NewLocalTime().String(),
		ActionOp:   actionOp,
		DcList:     dcList,
		Success:    success,
		Reason:     reason,
		Json:       json,
		Comment:    comment,
		OrgId:      orgId,
	}
}

func (d *Deployment) QueryDeploymentById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	//Prepare select-statement
	stmt, err := db.Prepare(DEPLOYMENT_SELECT)
	if err != nil {
		log.Errorf("QueryDeploymentById Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Query user by id
	var jsonFile []byte
	var comment []byte

	err = stmt.QueryRow(id).Scan(&d.Id, &d.Name, &d.ActionType, &d.ActionVerb, &d.ActionUrl, &d.ActionAt, &d.ActionOp, &d.DcList, &d.Success, &d.Reason, &jsonFile, &comment, &d.OrgId)
	if err != nil {
		log.Errorf("QueryDeploymentById Error: err=%s", err)
		return err
	}

	d.Json = string(jsonFile)
	d.Comment = string(comment)

	log.Infof("QueryDeploymentById: id=%d, name=%s, actionType=%d, actionVerb=%s, actionUrl=%s, actionAt=%s, actionOp=%d, dcList=%s, success=%d, reason=%s, json=%s, comment=%s, orgId=%d",
		d.Id, d.Name, d.ActionType, d.ActionVerb, d.ActionUrl, d.ActionAt, d.ActionOp, d.DcList, d.Success, d.Reason, d.Json, d.Comment, d.OrgId)

	return nil
}

func (d *Deployment) InsertDeployment() error {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(DEPLOYMENT_INSERT)
	if err != nil {
		log.Errorf("InsertDeployment Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update ActionAt
	d.ActionAt = localtime.NewLocalTime().String()

	// Insert a deployment
	_, err = stmt.Exec(d.Name, d.ActionType, d.ActionVerb, d.ActionUrl, d.ActionAt, d.ActionOp, d.DcList, d.Success, d.Reason, d.Json, d.Comment, d.OrgId)
	if err != nil {
		log.Errorf("InsertDeployment Error: err=%s", err)
		return err
	}

	log.Infof("InsertDeploymentById: id=%d, name=%s, actionType=%d, actionVerb=%s, actionUrl=%s, actionAt=%s, actionOp=%d, dcList=%s, success=%d, reason=%s, json=%s, comment=%s, orgId=%d",
		d.Id, d.Name, d.ActionType, d.ActionVerb, d.ActionUrl, d.ActionAt, d.ActionOp, d.DcList, d.Success, d.Reason, d.Json, d.Comment, d.OrgId)

	return nil
}

func QueryDeploymentByAppName(name string) ([]Deployment, error) {
	// New deployment point array
	deployments := make([]Deployment, 0)

	db := mysql.MysqlInstance().Conn()

	// Prepare select-by-name-statement
	stmt, err := db.Prepare(DEPLOYMENT_BYNAME)
	if err != nil {
		log.Errorf("QueryDeploymentByAppName Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		log.Errorf("QueryDeploymentByAppName Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := new(Deployment)

		var comment []byte
		err = rows.Scan(&d.Id, &d.Name, &d.ActionType, &d.ActionVerb, &d.ActionUrl, &d.ActionAt,
			&d.ActionOp, &d.DcList, &d.Success, &d.Reason, &d.Json, &comment, &d.OrgId)
		d.Comment = string(comment)
		if err != nil {
			log.Errorf("QueryDeploymentByAppName Error: err=%s", err)
			return nil, err
		}
		deployments = append(deployments, *d)

		log.Infof("QueryDeploymentByAppName: id=%d, name=%s, actionType=%d, actionVerb=%s, actionUrl=%s, actionAt=%s, actionOp=%d, dcList=%s, success=%d, reason=%s, json=%s, comment=%s, orgId=%d",
			d.Id, d.Name, d.ActionType, d.ActionVerb, d.ActionUrl, d.ActionAt, d.ActionOp, d.DcList, d.Success, d.Reason, d.Json, d.Comment, d.OrgId)
	}

	return deployments, nil
}

func StatDeploymentByActionType(actionType int) (count int32, err error) {

	db := mysql.MysqlInstance().Conn()

	// Prepare select-by-actionType
	stmt, err := db.Prepare(DEPLOYMENT_ACTIONTYPE_STAT)
	if err != nil {
		log.Errorf("StatDeploymentByActionType Error: err=%s", err)
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(actionType)
	if err != nil {
		log.Errorf("StatDeploymentByActionType Error: err=%s", err)
		return 0, err
	}
	defer rows.Close()

	count = 0
	for rows.Next() {
		d := new(Deployment)
		var comment []byte

		err = rows.Scan(&d.Id, &d.Name, &d.ActionType, &d.ActionVerb, &d.ActionUrl, &d.ActionAt,
			&d.ActionOp, &d.DcList, &d.Success, &d.Reason, &d.Json, &comment, &d.OrgId)
		d.Comment = string(comment)
		if err != nil {
			log.Errorf("StatDeploymentByActionType Error: err=%s", err)
			return 0, err
		}

		count++

		log.Debugf("QueryDeploymentByAppName: id=%d, name=%s, actionType=%d, actionVerb=%s, actionUrl=%s, actionAt=%s, actionOp=%d, dcList=%s, success=%d, reason=%s, json=%s, comment=%s, orgId=%d", d.Id, d.Name, d.ActionType, d.ActionVerb, d.ActionUrl, d.ActionAt, d.ActionOp, d.DcList, d.Success, d.Reason, d.Json, d.Comment, d.OrgId)

	}

	return count, nil
}

func QueryOperationStat() (OperationStat, error) {

	ops := make(OperationStat)

	db := mysql.MysqlInstance().Conn()
	//Prepare select-statement
	stmt, err := db.Prepare(OPERATION_LOG)
	if err != nil {
		log.Errorf("QueryDeploymentById Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Errorf("QueryOperationStat Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var op, total int32
		err = rows.Scan(&date, &total, &op)
		if err != nil {
			log.Errorf("QueryOperationStat Scan Error: err=%s", err)
			return nil, err
		}
		log.Infof("QueryOperationStat Scan: date=%s, total=%d, op=%d", date, total, op)
		// date exist in OperationStat
		if statistics, ok := ops[date]; ok {
			statistics[op] = total

		} else {
			s := make(Statistics)
			s[op] = total
			ops[date] = s
		}
	}

	log.Infof("QueryOperationStat: len(OperationStat)=%d", len(ops))
	return ops, nil
}
