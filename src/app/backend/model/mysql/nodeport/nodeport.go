package nodeport

import (
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	mydatacenter "app/backend/model/mysql/datacenter"

	//temporarily
	"errors"
)



func NewNodePort(port, dcId int32, svcName string, modifiedOp int32) *NodePort {
	return &NodePort{
		Port:       port,
		DcId:       dcId,
		SvcName:    svcName,
		Status:     VALID,
		CreatedAt:  localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment:    "",
	}
}

// 推荐一个未使用的或VALID的端口
// TODO: check the recommand algorithm
/*
func Recommand(datacenters []mydatacenter.DataCenter) (np *NodePort) {
	availablePorts := make(map[int32]int)
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
*/

func Recommand(datacenters []mydatacenter.DataCenter) (np *NodePort) {
	dcIdList := make([]int32, 0)

	for _, dc := range datacenters {
		dcIdList = append(dcIdList, dc.Id)
	}

	availableNodePorts := make(map[int32]int)

	for _, dcId := range dcIdList {
		npList, err := QueryNodePortByDcIdIfValid(dcId)
		if err != nil { 
			log.Errorf("Recommand QueryNodeByDcIdIfValid Error: err=%s", err)
			return nil
		}
	
		for _, np := range npList {
			availableNodePorts[np.Port] += 1
			if availableNodePorts[np.Port] == len(dcIdList) {
				log.Infof("Recommand nodePort=%d", np.Port)
				return &np
			}
		}
	}

	log.Infof("ddc nodePort Failed!")
	return nil
}

func QueryNodePortByDcIdIfValid(dcId int32)([]NodePort, error) {
	nodeports := make([]NodePort, 0)

	db := mysql.MysqlInstance().Conn()

	// Prepare select-all-statement
	stmt, err := db.Prepare(NP_SELECT_VALID_BY_DC)
	if err != nil {
		log.Errorf("QueryNodePortByDcIdIfValid Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(dcId)
	if err != nil {
		log.Errorf("QueryNodePortByDcIdIfValid Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		np := new(NodePort)
		err = rows.Scan(&np.Port, &np.DcId, &np.SvcName, &np.Status, &np.CreatedAt, &np.ModifiedAt, &np.ModifiedOp, &np.Comment)
		if err != nil {
			log.Errorf("QueryNodePortByDcIdIfValid Error: err=%s", err)
			return nil, err
		}

		nodeports = append(nodeports, *np)

		log.Infof("QueryNodePortByDcIdIfValid: port=%d, dcId=%d, svcName=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d, comment=%s",
		np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	}

	log.Infof("QueryNodePortByDcIdIfValid: len(nodeports)=%d", len(nodeports))
	return nodeports, nil

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
				log.Errorf("QueryNewNodePort Error: error=%s", err)
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

func QueryAllInvalidNodePort() ([]NodePort, error) {

	nodeports := make([]NodePort, 0)

	db := mysql.MysqlInstance().Conn()

	// Prepare select-all-statement
	stmt, err := db.Prepare(NP_SELECT_INVALID)
	if err != nil {
		log.Errorf("QueryAllInvalidNodePort Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Errorf("QueryAllInvalidNodePort Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		np := new(NodePort)
		err = rows.Scan(&np.Port, &np.DcId, &np.SvcName, &np.Status, &np.CreatedAt, &np.ModifiedAt, &np.ModifiedOp, &np.Comment)
		if err != nil {
			log.Errorf("QueryAllInvalidNodePort Error: err=%s", err)
			return nil, err
		}

		nodeports = append(nodeports, *np)

		log.Infof("QueryAllInvalidNodePort: port=%d, dcId=%d, svcName=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d, comment=%s",
		np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	}

	log.Infof("QueryAllInvalidNodePort: len(nodeports)=%d", len(nodeports))
	return nodeports, nil

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
	stmt, err := db.Prepare(NP_SELECT_BY_DC_AND_PORT)
	if err != nil {
		log.Errorf("QueryNodePortByPortAndDcId Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(port, dcId).Scan(&np.Port, &np.DcId, &np.SvcName, &np.Status, &np.CreatedAt, &np.ModifiedAt, &np.ModifiedOp, &np.Comment)
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

func (np *NodePort) QueryNodePortByPortAndDcIdIfValid(port, dcId int32) error {
	db := mysql.MysqlInstance().Conn()

	//"SELECT port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE status=1 AND port=? AND dcId=?"
	stmt, err := db.Prepare(NP_SELECT_VALID_BY_DC_AND_PORT)
	if err != nil {
		log.Errorf("QueryNodePortByPortAndDcIdIfValid Error: err=%s", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment, np.Port, np.DcId)
	if err != nil {
		log.Errorf("QueryNodePortByPortAndDcIdIfValid Error: err=%s", err)
		return err
	}

	return nil

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
	//np.Status = INVALID
	np.Status = VALID
	np.CreatedAt = localtime.NewLocalTime().String()
	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op

	//backup for InsertOnDuplicateKeyUpdateStatus
	svcName := np.SvcName
	status := np.Status

	//err = np.QueryNodePortByPortAndDcId(np.Port, np.DcId)
	err = np.QueryNodePortByPortAndDcIdIfValid(np.Port, np.DcId)

	// if port with dcId exist and it's INVALID
	if err == nil && np.Status == INVALID {
		log.Errorf("InserNodePort Error: Error=%s", "Port exists")
		return errors.New("Port Exists")
	} else {
		// if port with dcId exist and it's VALID || port with dcId doesn't exist
		err := np.InsertOnDuplicateKeyUpdateStatus(svcName, status, op)
		if err != nil {
			log.Errorf("InsertOnDuplicateKeyUpstateStatus Error: error=%s", err)
		}
		return nil
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

	log.Debugf("InsertOnDuplicateKeyUpdateStatus: Port=%d, DcId=%d, SvcName=%s, Status=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
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
	//np.Status = VALID
	np.Status = INVALID

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


func (np *NodePort) UseNodePort(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_UPDATE)
	if err != nil {
		log.Errorf("UseNodePort Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op
	np.Status = INVALID

	_, err = stmt.Exec(np.Port, np.DcId, np.SvcName, np.Status, np.ModifiedAt, np.ModifiedOp, np.Comment, np.Port, np.DcId)
	// update error
	if err != nil {
		log.Errorf("UseNodePort Error: error=%s", err)
		return err
	}

	// update ok or even no update(not exist)
	log.Debugf("UseNodePort: Port=%d, DcId=%d, SvcName=%s, Status=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}

func (np *NodePort) ReleaseNodePort(op int32) error {

	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(NP_UPDATE)
	if err != nil {
		log.Errorf("ReleaseNodePort Error: error=%s", err)
		return err
	}
	defer stmt.Close()

	np.ModifiedAt = localtime.NewLocalTime().String()
	np.ModifiedOp = op
	np.Status = VALID

	_, err = stmt.Exec(np.Port, np.DcId, np.SvcName, np.Status, np.ModifiedAt, np.ModifiedOp, np.Comment, np.Port, np.DcId)
	// update error
	if err != nil {
		log.Errorf("ReleaseNodePort Error: error=%s", err)
		return err
	}

	// update ok or even no update(not exist)
	log.Debugf("ReleaseNodePort: Port=%d, DcId=%d, SvcName=%s, Status=%d, CreatedAt=%s, ModifiedAt=%s, ModifiedOp=%d, Comment=%s", np.Port, np.DcId, np.SvcName, np.Status, np.CreatedAt, np.ModifiedAt, np.ModifiedOp, np.Comment)
	return nil
}
