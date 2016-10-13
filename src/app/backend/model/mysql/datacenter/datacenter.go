package datacenter

import (
	mylog "app/backend/common/util/log"
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"encoding/json"
)

var log = mylog.Log

const (
	DC_SELECT = "SELECT id, name, host, port, secret, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM datacenter WHERE id=?"

	DC_SELECT_ALL = "SELECT id, name, host, port, secret, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM datacenter"

	DC_INSERT = "INSERT INTO " +
		"datacenter(name, host, port, secret, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"

	DC_UPDATE = "UPDATE datacenter " +
		"SET name=?, host=?, port=?, secret=?, status=?, modifiedAt=?, modifiedOp=?, comment=? " +
		"WHERE id=?"

	DC_DELETE = "UPDATE datacenter " +
		"SET status=?, modifiedAt=?, modifiedOp=? " +
		"WHERE id=?"

	VALID   = 1
	INVALID = 0
)

type DataCenter struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int32  `json:"port"`
	Secret     string `json:"secret"` // maybe error
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

func NewDataCenter(name, host, secret, comment string, port, modifiedOp int32) *DataCenter {
	return &DataCenter{
		Name:       name,
		Host:       host,
		Port:       port,
		Secret:     secret,
		Status:     VALID,
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

	rows, err := stmt.Query()
	if err != nil {
		log.Errorf("QueryAllDatacenters Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		dc := new(DataCenter)
		var secret []byte
		var comment []byte
		err = rows.Scan(&dc.Id, &dc.Name, &dc.Host, &dc.Port, &secret, &dc.Status,
			&dc.CreatedAt, &dc.ModifiedAt, &dc.ModifiedOp, &comment)
		dc.Comment = string(comment)
		dc.Secret = string(secret)
		if err != nil {
			log.Errorf("QueryAllDatacenters Error: err=%s", err)
			return nil, err
		}
		datacenters = append(datacenters, *dc)

		log.Infof("QueryAllDatacenters: id=%d, name=%s, host=%d, port=%d, secret=%s, status=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
			dc.Id,  dc.Name,  dc.Host,  dc.Port,  dc.Secret,  dc.Status, dc.CreatedAt,  dc.ModifiedAt,  dc.ModifiedOp,  comment)

	}

	log.Infof("QueryAllDatacenters: len(datacenters)=%d", len(datacenters))
	return datacenters, nil
}

func (dc *DataCenter) QueryDataCenterById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(DC_SELECT)
	if err != nil {
		log.Errorf("QueryDataCenterById Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	var secret []byte
	var comment []byte
	// Query idc by id
	err = stmt.QueryRow(id).Scan(&dc.Id, &dc.Name, &dc.Host, &dc.Port, &secret, &dc.Status,
		&dc.CreatedAt, &dc.ModifiedAt, &dc.ModifiedOp, &comment)

	dc.Secret = string(secret)
	dc.Comment = string(comment)

	if err != nil {
		log.Errorf("QueryDataCenterById Error: err=%s", err)
		return err
	}

	log.Infof("QureyDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)

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

	// Insert a idc
	_, err = stmt.Exec(dc.Name, dc.Host, dc.Port, dc.Secret, dc.Status,
		dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp, dc.Comment)

	if err != nil {
		log.Errorf("InsertDataCenter Error: err=%s", err)
		return err
	}

	log.Infof("InsertDataCenterById: id=%d, name=%s, host=%d, port=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
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
	_, err = stmt.Exec(dc.Name, dc.Host, dc.Port, dc.Secret, dc.Status,
		dc.ModifiedAt, dc.ModifiedOp, dc.Comment, dc.Id)

	if err != nil {
		log.Errorf("DataCenter UpdateDataCenter Exec Error: err=%s", err)
		return err
	}

	log.Infof("UpdateDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
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

	log.Infof("DeleteDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
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

