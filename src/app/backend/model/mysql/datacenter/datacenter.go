package datacenter

import (
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"encoding/json"
	"log"
)

const (
	DC_SELECT = "SELECT id, name, host, port, secret, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM datacenter WHERE id=?"

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

func (dc *DataCenter) QueryDataCenterById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(DC_SELECT)
	if err != nil {
		log.Printf("QueryDataCenterById Error: err=%s\n", err)
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
		log.Printf("QueryDataCenterById Error: err=%s\n", err)
		return err
	}

	log.Printf("QureyDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)

	return nil
}

func (dc *DataCenter) InsertDataCenter(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(DC_INSERT)
	if err != nil {
		log.Printf("InsertDataCenter Error: err=%s\n", err)
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
		log.Printf("InsertDataCenter Error: err=%s\n", err)
		return err
	}

	log.Printf("InsertDataCenterById: id=%d, name=%s, host=%d, port=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp\n",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil
}

func (dc *DataCenter) UpdateDataCenter(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare update-statement
	stmt, err := db.Prepare(DC_UPDATE)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return err
	}

	log.Printf("UpdateDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil

}

func (dc *DataCenter) DeleteDataCenter(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare delete-statement
	stmt, err := db.Prepare(DC_DELETE)
	if err != nil {
		log.Printf("DeleteDatCenter Error: err=%s\n", err)
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
		log.Printf("DeleteDataCenter Error: err=%s\n", err)
		return err
	}

	log.Printf("DeleteDataCenterById: id=%d, name=%s, host=%s, port=%d, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		dc.Id, dc.Name, dc.Host, dc.Port, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil
}

func (dc *DataCenter) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), dc)

	if err != nil {
		log.Printf("DecodeJson Error: err=%s\n", err)
	}
}

func (dc *DataCenter) EncodeJson() string {
	data, err := json.Marshal(dc)
	if err != nil {
		log.Printf("EncdoeJson Error: err=%s\n", err)
		return ""
	}
	return string(data)
}
