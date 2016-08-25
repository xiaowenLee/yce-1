package nodeport

import (
	mylog "app/backend/common/util/log"
	mysql "app/backend/common/util/mysql"
	"strconv"

	//temporarily
	"errors"
)

var log = mylog.Log


//选择时需要同时满足dcId和port
const (
	NP_SELECT = "SELECT id, port, dcId, status " + "FROM nodeport WHERE port=? AND dcId=?"

	NP_INSERT = "INSERT INTO " + "nodeport(port, dcId, status) " + "VALUES(?, ?, ?)"

	NP_UPDATE = "UPDATE nodeport " + "SET status=? " + "WHERE port=? AND dcID=?"

	NP_DELETE = "DELETE FROM nodeport " + "WHERE port=? AND dcId=?"
)

type NodePort struct {
	Id int32 `json:"id"`
	Port int32 `json:"port"`
	DcID int32 `json:"dcID"`
	Status string `json:"status"`
}

func NewNodePort(port, dcId int32) *NodePort{
	return &NodePort{
		Port: port,
		DcID: dcId,
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

	err = stmt.QueryRow(port, dcId).Scan(&np.Id, &np.Port, &np.DcID, &np.Status)
	// port with dcId not exist
	if err != nil {
		log.Errorf("QuertNodePortByPortAndDcId Error: error=%s", err)
		return err
	}

	// port with dcId exist
	log.Infof("QueryNodePortByPortAndDcId: Id=%d, Port=%d, DcID=%d, Statu=%s", np.Id, np.Port, np.DcID, np.Status)
	return nil
}

// 插入port, 如果已存在,应该返回err(但目前没有), 如果不存在,插入成功返回Nil, 插入失败返回err
func (np *NodePort) InsertNodePort(port, dcId int32, status string) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_INSERT)
	if err != nil {
		log.Errorf("InsertNodePort Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	//TODO: ensure the dc is owned by the user
	//TODO: check the bound of NodePort
	np.Port = port
	np.DcID = dcId
	np.Status = status

	err = np.QueryNodePortByPortAndDcId(np.Port, np.DcID)
	// if port with dcId exist
	if err == nil {
		log.Errorf("InserNodePort Error: error=%s", "Port exists")
		// err is nil, need to new one("Port exists")
		return errors.New("Port Exists")
	}

	// if port with dcId not exist
	_, err = stmt.Exec(np.Port, np.DcID, np.Status)

	// if insert failed
	if err != nil {
		log.Errorf("InsertNodePort Error: error=%s", err)
		return err
	}
	log.Infof("InsertNodePort: port=%d, dcId=%d, status=%s", np.Port, np.DcID, np.Status)
	return nil
}

// 更新port对应的信息, 该记录存在更新成功返回nil, 该记录不存在或更新失败返回err
func (np *NodePort) UpdateNodePortByPortAndDcId(port, dcId int32, status string) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_UPDATE)
	if err != nil {
		log.Errorf("UpdateNodePortByPortAndDcId Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	np.Port = port
	np.DcID = dcId
	np.Status = status
	_, err = stmt.Exec(np.Status, np.Port, np.DcID)
	// update error
	if err != nil {
		log.Errorf("UpdateNodePortByPortAndDcId Error: error=%s", err)
		return err
	}

	// update ok or even no updation
	log.Infof("UpdateNodePortByPortAndDcId: Port=%d, DcID=%d, Status=%s", np.Port, np.DcID, np.Status)
	return nil
}


// 删除port对应的信息, 该记录存在删除成功返回nil, 该记录不存在或删除失败返回err
func (np *NodePort) DeleteNodePortByPortAndDcId(port, dcId int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_DELETE)
	if err != nil {
		log.Errorf("DeleteNodeByPortAndDcId Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	np.Port = port
	np.DcID = dcId
	_, err = stmt.Exec(strconv.Itoa(int(np.Port)), strconv.Itoa(int(np.DcID)))
	// delete error
	if err != nil {
		log.Errorf("DeleteNodePortByPortAndDcId Error: error=%s", err)
		return err
	}

	// delete success or even no deletion
	log.Infof("DeleteNodePortByPortAndDcId: Port=%d, dcId=%d", np.Port, np.DcID)
	return nil

}