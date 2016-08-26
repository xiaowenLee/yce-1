package nodeport

import (
	mylog "app/backend/common/util/log"
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"

	//temporarily
	"errors"

)

var log = mylog.Log

//选择时需要同时满足dcId和port
const (
	NP_SELECT = "SELECT id, port, dcId, svcId, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE port=? AND dcId=?"

	NP_INSERT = "INSERT INTO " + "nodeport(port, dcId, svcId, status, createdAt, modifiedAt, modifiedOp, comment) " + "VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	NP_UPDATE = "UPDATE nodeport " + "SET port=?, dcId=?, svcId=?, status=?, modifiedAt=?, modifiedOp=?, comment=? " + "WHERE port=? AND dcId=?"

	NP_DELETE = "DELETE FROM nodeport " + "WHERE port=? AND dcId=?"

	VALID = 1
	INVALID = 0
)

type NodePort struct {
	Id         int32  `json:"id"`
	Port       int32  `json:"port"`
	DcId       int32  `json:"dcId"`
	SvcId      int32  `json:"svcId"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `jsno:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

func NewNodePort(port, dcId, svcId, modifiedOp int32) *NodePort {
	return &NodePort{
		Port: port,
		DcId: dcId,
		SvcId: svcId,
		Status: VALID,
		CreatedAt: localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment: "",
	}
}

// 查询是否存在该port和dcId组合, 如果存在,返回Nil, 如果不存在,返回err
func (np *NodePort) QueryNodePortByPortAndDcId(port, dcId int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(NP_SELECT)
	if err != nil {
		log.Errorf("QueryNodePortByPortAndDcId Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(port, dcId).Scan(&np.Id, &np.Port, &np.DcId, &np.SvcId, &np.Status, &np.CreatedAt, &np.ModifiedAt, &np.ModifiedOp, &np.Comment)
	// port with dcId not exist
	if err != nil {
		log.Errorf("QuertNodePortByPortAndDcId Error: error=%s", err)
		return err
	}

	// port with dcId exist
	log.Debugf("QueryNodePortByPortAndDcId: Id=%d, Port=%d, DcId=%d, SvcId=%d, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%s, Comment=%s", np.Id, np.Port, np.DcId, np.SvcId, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}


// 根据NodePort号和所属DcId号查找相应的serviceId, 存在返回ServiceId和Nil, 不存在返回-1和err
func (np *NodePort) QueryServiceIdByPortAndDcId(port, dcId int32) (int32, error) {
	err := np.QueryNodePortByPortAndDcId(port, dcId)
	if err != nil {
		log.Errorf("QueryServiceIdByPortAndDcId Error: error=%s", err)
		return -1, err
	}

	log.Debugf("QueryServiceIdByPortAndDcId: Id=%d, Port=%d, DcId=%d, SvcId=%d, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%s, Comment=%s", np.Id, np.Port, np.DcId, np.SvcId, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return np.SvcId, nil
}

// 插入port, 如果已存在,应该返回err, 如果不存在,插入成功返回Nil, 插入失败返回err
func (np *NodePort) InsertNodePort(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_INSERT)
	if err != nil {
		log.Errorf("InsertNodePort Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	//TODO: ensure the dc is owned by the user
	//TODO: check the bound of NodePort
	np.Status = INVALID
	np.CreatedAt = localtime.NewLocalTime().String()
	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op

	err = np.QueryNodePortByPortAndDcId(np.Port, np.DcId)
	// if port with dcId exist
	if err == nil {
		log.Errorf("InserNodePort Error: error=%s", "Port exists")
		// err is nil, need to new one("Port exists")
		return errors.New("Port Exists")
	}

	// if port with dcId not exist
	_, err = stmt.Exec(np.Port, np.DcId, np.SvcId, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)

	// if insert failed
	if err != nil {
		log.Errorf("InsertNodePort Error: error=%s", err)
		return err
	}

	log.Debugf("InsertNodePort: Id=%d, Port=%d, DcId=%d, SvcId=%d, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%s, Comment=%s", np.Id, np.Port, np.DcId, np.SvcId, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}

// 更新port对应的信息, 该记录不存在或该记录存在更新成功返回nil, 更新失败返回err
func (np *NodePort) UpdateNodePortByPortAndDcId(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_UPDATE)
	if err != nil {
		log.Errorf("UpdateNodePortByPortAndDcId Error: error=%s", err)
		return err
	}
	defer stmt.Close()


	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op

	_, err = stmt.Exec(np.Port, np.DcId, np.SvcId, np.Status, np.ModifiedAt, np.ModifiedOp, np.Comment, np.Port, np.DcId)
	// update error
	if err != nil {
		log.Errorf("UpdateNodePortByPortAndDcId Error: error=%s", err)
		return err
	}

	// update ok or even no updation(not exist)
	log.Debugf("UpdateNodePortByPortAndDcId: Id=%d, Port=%d, DcId=%d, SvcId=%d, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%s, Comment=%s", np.Id, np.Port, np.DcId, np.SvcId, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}

// 删除port对应的信息,  该记录不存在或该记录存在删除成功返回nil, 删除失败返回err
func (np *NodePort) DeleteNodePortByPortAndDcId(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_DELETE)
	if err != nil {
		log.Errorf("DeleteNodeByPortAndDcId Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op

	_, err = stmt.Exec(np.Port, np.DcId)
	// delete error
	if err != nil {
		log.Errorf("DeleteNodePortByPortAndDcId Error: error=%s", err)
		return err
	}

	// delete success or even no deletion
	log.Debugf("DeleteNodePortByPortAndDcId: Id=%d, Port=%d, DcId=%d, SvcId=%d, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%s, Comment=%s", np.Id, np.Port, np.DcId, np.SvcId, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}
