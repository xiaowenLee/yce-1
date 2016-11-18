package datacenter

import (
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"encoding/json"
)


func NewDataCenter(name, host, secret, nodePort, comment string, port, modifiedOp int32) *DataCenter {
	return &DataCenter{
		Name:       name,
		Host:       host,
		Port:       port,
		Secret:     secret,
		Status:     VALID,
		NodePort:   nodePort,
		CreatedAt:  localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment:    comment,
	}
}

func QueryAllDatacenters() ([]DataCenter, error) {
	datacenters := make([]DataCenter, 0)

	db := mysql.MysqlInstance().Conn()

	// Prepare select-all-statement
	stmt, err := db.Prepare(DC_SELECT_ALL)
	if err != nil {
		log.Errorf("QueryAllDatacenters Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(VALID)
	if err != nil {
		log.Errorf("QueryAllDatacenters Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dc := new(DataCenter)
		var secret []byte
		var comment []byte
		var nodePort []byte
		err = rows.Scan(&dc.Id, &dc.Name, &dc.Host, &dc.Port, &secret, &dc.Status, &nodePort,
			&dc.CreatedAt, &dc.ModifiedAt, &dc.ModifiedOp, &comment)
		dc.Comment = string(comment)
		dc.Secret = string(secret)
		dc.NodePort = string(nodePort)
		if err != nil {
			log.Errorf("QueryAllDatacenters Error: err=%s", err)
			return nil, err
		}
		datacenters = append(datacenters, *dc)

		log.Infof("QueryAllDatacenters: id=%d, name=%s, host=%d, port=%d, secret=%s, status=%s, nodePort=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
			dc.Id,  dc.Name,  dc.Host,  dc.Port,  dc.Secret,  dc.Status, dc.NodePort, dc.CreatedAt,  dc.ModifiedAt,  dc.ModifiedOp,  comment)

	}

	log.Infof("QueryAllDatacenters: len(datacenters)=%d", len(datacenters))
	return datacenters, nil
}

func (dc *DataCenter) QueryDataCenterById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(DC_SELECT_BY_ID)
	if err != nil {
		log.Errorf("QueryDataCenterById Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	var secret []byte
	var comment []byte
	var nodePort []byte
	// Query idc by id
	err = stmt.QueryRow(id,VALID).Scan(&dc.Id, &dc.Name, &dc.Host, &dc.Port, &secret, &dc.Status, &nodePort,
		&dc.CreatedAt, &dc.ModifiedAt, &dc.ModifiedOp, &comment)

	dc.Secret = string(secret)
	dc.Comment = string(comment)
	dc.NodePort = string(nodePort)

	if err != nil {
		log.Errorf("QueryDataCenterById Error: err=%s", err)
		return err
	}

	log.Infof("QureyDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, nodePort=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.NodePort, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)

	return nil
}

func (dc *DataCenter) QueryDataCenterByName(name string) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(DC_SELECT_BY_NAME)
	if err != nil {
		log.Errorf("QueryDataCenterByName Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	var secret []byte
	var comment []byte
	var nodePort []byte
	// Query idc by id
	err = stmt.QueryRow(name, VALID).Scan(&dc.Id, &dc.Name, &dc.Host, &dc.Port, &secret, &dc.Status, &nodePort,
		&dc.CreatedAt, &dc.ModifiedAt, &dc.ModifiedOp, &comment)

	dc.Secret = string(secret)
	dc.Comment = string(comment)
	dc.NodePort = string(nodePort)

	if err != nil {
		log.Errorf("QueryDataCenterByName Error: err=%s", err)
		return err
	}

	log.Infof("QueryDataCenterByName: id=%d, name=%s, host=%s, port=%d, status=%d, nodePort=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.NodePort, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)

	return nil
}

func (dc *DataCenter) QueryDuplicatedName(name string) error {

	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(DC_QUERY_DUPLICATED_NAME)
	if err != nil {
		log.Errorf("QueryDuplicatedName Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	var secret []byte
	var comment []byte
	var nodePort []byte
	// Query idc by id
	err = stmt.QueryRow(name, VALID).Scan(&dc.Id, &dc.Name, &dc.Host, &dc.Port, &secret, &dc.Status, &nodePort,
		&dc.CreatedAt, &dc.ModifiedAt, &dc.ModifiedOp, &comment)

	dc.Secret = string(secret)
	dc.Comment = string(comment)

	if err != nil {
		log.Errorf("QueryDuplicatedName Error: err=%s", err)
		return err
	}

	log.Infof("QureyDuplicatedName: id=%d, name=%s, host=%s, port=%d, status=%d, nodePort=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.NodePort, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)

	return nil
}


func (dc *DataCenter) InsertDataCenter(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(DC_INSERT)
	if err != nil {
		log.Errorf("InsertDataCenter Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update createdAt, modifiedAt, modifiedOp
	dc.CreatedAt = localtime.NewLocalTime().String()
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op
	dc.Status = VALID

	// Insert a idc
	_, err = stmt.Exec(dc.Name, dc.Host, dc.Port, dc.Secret, dc.Status, dc.NodePort,
		dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp, dc.Comment)

	if err != nil {
		log.Errorf("InsertDataCenter Error: err=%s", err)
		return err
	}

	log.Infof("InsertDataCenterById: id=%d, name=%s, host=%d, port=%d, status=%d, nodePort=%s, createdAt=%s, modifiedAt=%s, modifiedOp",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.NodePort, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil
}

func (dc *DataCenter) UpdateDataCenter(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare update-statement
	stmt, err := db.Prepare(DC_UPDATE)
	if err != nil {
		log.Errorf("DataCenter UpdateDataCenter Prepare Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op

	// Update a idc
	_, err = stmt.Exec(dc.Name, dc.Host, dc.Port, dc.Secret, dc.Status, dc.NodePort,
		dc.ModifiedAt, dc.ModifiedOp, dc.Comment, dc.Id)

	if err != nil {
		log.Errorf("DataCenter UpdateDataCenter Exec Error: err=%s", err)
		return err
	}

	log.Infof("UpdateDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, nodePort=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.NodePort, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil

}

func (dc *DataCenter) DeleteDataCenter(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare delete-statement
	stmt, err := db.Prepare(DC_DELETE)
	if err != nil {
		log.Errorf("DeleteDatCenter Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt and modifiedOp
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op

	// Set idc status to INVALID
	dc.Status = INVALID
	_, err = stmt.Exec(dc.Status, dc.ModifiedAt, dc.ModifiedOp, dc.Id)
	if err != nil {
		log.Errorf("DeleteDataCenter Error: err=%s", err)
		return err
	}

	log.Infof("DeleteDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, nodePort=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.NodePort, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil
}

func (dc *DataCenter) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), dc)

	if err != nil {
		log.Errorf("DecodeJson Error: err=%s", err)
	}
}

func (dc *DataCenter) EncodeJson() string {
	data, err := json.Marshal(dc)
	if err != nil {
		log.Errorf("EncdoeJson Error: err=%s", err)
		return ""
	}
	return string(data)
}

