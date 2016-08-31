package nodeport

import (
	mylog "app/backend/common/util/log"
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	mydatacenter "app/backend/model/mysql/datacenter"

	//temporarily
	"errors"

)

var log = mylog.Log

//选择时需要同时满足dcId和port
const (
	NP_SELECT = "SELECT svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE port=? AND dcId=?"

	NP_INSERT = "INSERT INTO " + "nodeport(port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment) " + "VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	NP_UPDATE = "UPDATE nodeport " + "SET port=?, dcId=?, svcName=?, status=?, modifiedAt=?, modifiedOp=?, comment=? " + "WHERE port=? AND dcId=?"

	NP_DELETE = "UPDATE nodeport " + "SET status=?, modifiedAt=?, modifiedOp=?, comment=? " + "WHERE port=? AND dcId=?"

	NP_INSERT_ON_DUPLICATE_KEY_UPDATE = "INSERT INTO nodeport(port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment) " + "VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE svcName=?, status=?, modifiedAt=?, modifiedOp=? "

	NP_SELECT_VALID = "SELECT port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE status=? AND dcId=?"

	NP_SELECT_NEW = "SELECT port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE port=? AND dcId=? ORDER BY port DESC"

	VALID = 1
	INVALID = 0
	PORT_START = 30000
	PORT_LIMIT = 32767
)

type NodePort struct {
	Port       int32  `json:"port"`
	DcId       int32  `json:"dcId"`
	SvcName    string `json:"svcName"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `jsno:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

func NewNodePort(port, dcId int32, svcName string, modifiedOp int32) *NodePort {
	return &NodePort{
		Port: port,
		DcId: dcId,
		SvcName: svcName,
		Status: VALID,
		CreatedAt: localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment: "",
	}
}

// 推荐一个未使用的或VALID的端口
// TODO: check the recommand algorithm
func Recommand(datacenters []mydatacenter.DataCenter) (np *NodePort) {
	availablePorts := make(map[int32] int, 1)
	np = new(NodePort)

	for _, v := range datacenters {
		err := np.QueryValidNodePort(v.Id)
		if err != nil {
			mylog.Log.Errorf("Recommand From Valid NodePort failed ")
			break
		}
		availablePorts[np.Port] += 1
	}

	if availablePorts[np.Port] == len(datacenters) {
		mylog.Log.Infof("Recommand NodePort: nodePort=%d", np.Port)
		return np
	} else {
		np.QueryNewNodePort(datacenters)

		return np
	}

}

func (np *NodePort) QueryNewNodePort(datacenters []mydatacenter.DataCenter) {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(NP_SELECT_NEW)
	if err != nil {
		log.Errorf("QueryValidNodePort Error: error=%s", err)
		return
	}

	for i := PORT_START; i <= PORT_LIMIT; i++ {
		var comment []byte
		var j int
		for j = 0; j < len(datacenters); j++ {
			err := stmt.QueryRow(i, datacenters[j].Id).Scan(&np.Port, &np.DcId, &np.SvcName, &np.Status, &np.CreatedAt, &np.ModifiedAt, &np.ModifiedOp, &comment)
			if err != nil {
				mylog.Log.Errorf("QueryNewNodePort Error: error=%s", err)
				continue
			} else {
				break
			}
		}

		if j >= len(datacenters) {
			np.Port = int32(i)
			log.Infof("QueryNewNodePort: Port=%d, DcId=%d, SvcName=%s, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
			return
		}
	}

	log.Debugf("QueryNewNodePort failed")

}


func (np *NodePort) QueryValidNodePort(dcId int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(NP_SELECT_VALID)
	if err != nil {
		log.Errorf("QueryValidNodePort Error: error=%s", err)
		return err
	}
	var comment []byte
	err = stmt.QueryRow(VALID, dcId).Scan(&np.Port, &np.DcId, &np.SvcName, &np.Status, &np.CreatedAt, &np.ModifiedAt, &np.ModifiedOp, &comment)
	np.Comment = string(comment)

	if err != nil {
		log.Errorf("QueryValidNodePort Error: error=%s", err)
		return err
	}
	log.Debugf("QueryValidNodePort: Port=%d, DcId=%d, SvcName=%s, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
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

	var comment []byte
	err = stmt.QueryRow(port, dcId).Scan(&np.SvcName, &np.Status, &np.CreatedAt, &np.ModifiedAt, &np.ModifiedOp, &comment)
	np.Comment = string(comment)

	// port with dcId not exist
	if err != nil {
		log.Errorf("QuertNodePortByPortAndDcId Error: error=%s", err)
		return err
	}

	// port with dcId exist
	np.Port = port
	np.DcId = dcId
	log.Debugf("QueryNodePortByPortAndDcId: Port=%d, DcId=%d, SvcName=%s, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}


// 根据NodePort号和所属DcId号查找相应的serviceId, 存在返回ServiceId和Nil, 不存在返回""和err
func (np *NodePort) QueryServiceNameByPortAndDcId(port, dcId int32) (string, error) {
	err := np.QueryNodePortByPortAndDcId(port, dcId)
	if err != nil {
		log.Errorf("QueryServiceNameByPortAndDcId Error: error=%s", err)
		return "", err
	}

	log.Debugf("QueryServiceNameByPortAndDcId: Port=%d, DcId=%d, SvcName=%s, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return np.SvcName, nil
}

// 插入port, 如果已存在且INVALID,返回err; 如果存在且为VALID, 更新记录的svcName, status, modifiedAt, modifiedOp, 并返回Nil。如果不存在,插入新记录并返回Nil, 插入失败返回err
func (np *NodePort) InsertNodePort(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_INSERT)
	if err != nil {
		log.Errorf("InsertNodePort Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	//NOTE: 确保用户拥有这个数据中心
	//NOTE: 检查NodePort是否超出了边界
	np.Status = INVALID
	np.CreatedAt = localtime.NewLocalTime().String()
	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op

	//backup for InsertOnDuplicateKeyUpdateStatus
	svcName := np.SvcName
	status := np.Status


	err = np.QueryNodePortByPortAndDcId(np.Port, np.DcId)

	// if port with dcId exist and it's INVALID
	if err == nil && np.Status == 0{
		log.Errorf("InserNodePort Error: error=%s", "Port exists")
		return errors.New("Port Exists")
	} else {
		// if port with dcId exist and it's VALID || port with dcId doesn't exist
		err := np.InsertOnDuplicateKeyUpdateStatus(svcName, status, op)
		if err != nil {
			log.Errorf("InsertOnDuplicateKeyUpstateStatus Error: error=%s", err)
		}
	}

	log.Debugf("InsertNodePort: Port=%d, DcId=%d, SvcName=%s, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}

//插入时,如果记录存在,就更新里面的一些字段,如果不存在则插入新记录。作用在具有唯一索引或主键。
func (np *NodePort) InsertOnDuplicateKeyUpdateStatus(svcName string, status int32, op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_INSERT_ON_DUPLICATE_KEY_UPDATE)
	if err != nil {
		log.Errorf("InsertOnDuplicateUpdateStatus Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	// update existed record.
	modifiedAt := localtime.NewLocalTime().String()
	modifiedOp := op

	_, err = stmt.Exec(np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment, svcName, status, modifiedAt, modifiedOp)
	if err != nil {
		log.Errorf("InsertOnDuplicateKeyUpdateStatus Error: error=%s", err)
		return err
	}

	log.Debugf("DeleteNodePorgOnDuplicateKeyUpdateStatus: Port=%d, DcId=%d, SvcName=%s, Status=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}

// 更新port对应的信息, 该记录不存在或该记录存在更新成功返回nil, 更新失败返回err
func (np *NodePort) UpdateNodePort(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_UPDATE)
	if err != nil {
		log.Errorf("UpdateNodePort Error: error=%s", err)
		return err
	}
	defer stmt.Close()


	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op

	_, err = stmt.Exec(np.Port, np.DcId, np.SvcName, np.Status, np.ModifiedAt, np.ModifiedOp, np.Comment, np.Port, np.DcId)
	// update error
	if err != nil {
		log.Errorf("UpdateNodePort Error: error=%s", err)
		return err
	}

	// update ok or even no update(not exist)
	log.Debugf("UpdateNodePort: Port=%d, DcId=%d, SvcName=%s, Status=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}


// 删除port对应的信息(修改status为VALID),  该记录不存在或该记录存在删除成功返回nil, 删除失败返回err
func (np *NodePort) DeleteNodePort(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_DELETE)
	if err != nil {
		log.Errorf("DeleteNodePort Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op
	np.Status = VALID

	_, err = stmt.Exec(np.Status, np.ModifiedAt, np.ModifiedOp, np.Comment, np.Port, np.DcId)
	// delete error
	if err != nil {
		log.Errorf("DeleteNodePort Error: error=%s", err)
		return err
	}

	// delete success or even no deletion
	log.Debugf("DeleteNodePort: Port=%d, DcId=%d, SvcName=%s, Statu=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}
