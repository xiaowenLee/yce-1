package quota

import (
	mylog "app/backend/common/util/log"
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"encoding/json"
)

var log = mylog.Log

const (
	QUOTA_SELECT = "SELECT id, name, cpu, mem, rbd, price, " +
		"status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM quota WHERE id=?"

	QUOTA_SELECT_ALL = "SELECT id, name, cpu, mem, rbd, price, " +
		"status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM quota"

	QUOTA_SELECT_ALL_ORDER_BY_CPU = "SELECT id, name, cpu, mem, rbd, price, " +
		"status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM quota " +
		"ORDER BY cpu ASC"

	QUOTA_INSERT = "INSERT INTO " +
		"quota(name, cpu, mem, rbd, price, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	QUOTA_UPDATE = "UPDATE quota " +
		"SET name=?, cpu=?, mem=?, rbd=?, price=?, status=?, modifiedAt=?, modifiedOp=?, comment=? " +
		"WHERE id=?"

	QUOTA_DELETE = "UPDATE quota " +
		"SET status=?, modifiedAt=?, modifiedOp=? " +
		"WHERE id=?"

	VALID   = 1
	INVALID = 0
)

type Quota struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Cpu        int32  `json:"cpu"`
	Mem        int32  `json:"mem"`
	Rbd        int32  `json:"rbd"`
	Price      string `json:"price"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

func NewQuota(name, price, comment string, cpu, mem, rbd, modifiedOp int32) *Quota {
	return &Quota{
		Name:       name,
		Cpu:        cpu,
		Mem:        mem,
		Rbd:        rbd,
		Price:      price,
		Status:     VALID,
		CreatedAt:  localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment:    comment,
	}
}

func (q *Quota) QueryQuotaById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(QUOTA_SELECT)
	if err != nil {
		log.Errorf("QueryQuotaById Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Query quota by id
	var comment []byte
	err = stmt.QueryRow(id).Scan(&q.Id, &q.Name, &q.Cpu, &q.Mem, &q.Rbd,
		&q.Price, &q.Status, &q.CreatedAt, &q.ModifiedAt, &q.ModifiedOp, &comment)

	q.Comment = string(comment)

	if err != nil {
		log.Errorf("QueryQuotaById Error: err=%s", err)
		return err
	}

	log.Infof("QueryQuotaById: id=%d, name=%s, cpu=%d, mem=%d, rbd=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		q.Id, q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status, q.CreatedAt, q.ModifiedAt, q.ModifiedOp)
	return nil
}

func (q *Quota) InsertQuota(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepared insert-statement
	stmt, err := db.Prepare(QUOTA_INSERT)
	if err != nil {
		log.Errorf("InsertQuota Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update createAt, modifiedAt, modifiedOp
	q.CreatedAt = localtime.NewLocalTime().String()
	q.ModifiedAt = localtime.NewLocalTime().String()
	q.ModifiedOp = op

	// Insert a user
	_, err = stmt.Exec(q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status,
		q.CreatedAt, q.ModifiedAt, q.ModifiedOp, q.Comment)

	if err != nil {
		log.Errorf("InsertQuota Error: err=%s", err)
		return err
	}

	log.Infof("QueryQuotaById: id=%d, name=%s, cpu=%d, mem=%d, rbd=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		q.Id, q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status, q.CreatedAt, q.ModifiedAt, q.ModifiedOp)
	return nil
}

func (q *Quota) UpdateQuota(op int32) error {

	db := mysql.MysqlInstance().Conn()

	// Prepared update-statement
	stmt, err := db.Prepare(QUOTA_UPDATE)
	if err != nil {
		log.Errorf("UpdateQuota Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	q.ModifiedAt = localtime.NewLocalTime().String()
	q.ModifiedOp = op

	// Update a quota
	_, err = stmt.Exec(q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status, q.ModifiedAt, q.ModifiedOp, q.Comment, q.Id)
	if err != nil {
		log.Errorf("UpdateQuota Error: err=%s", err)
		return nil
	}

	log.Infof("UpdateQuota: id=%d, name=%s, cpu=%d, mem=%d, rbd=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		q.Id, q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status, q.CreatedAt, q.ModifiedAt, q.ModifiedOp)
	return nil
}

func (q *Quota) DeleteQuota(op int32) error {

	db := mysql.MysqlInstance().Conn()

	// Prepared delete-statement
	stmt, err := db.Prepare(QUOTA_DELETE)
	if err != nil {
		log.Errorf("UpdateQuota Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	q.ModifiedAt = localtime.NewLocalTime().String()
	q.ModifiedOp = op
	q.Status = INVALID

	// Update a quota
	_, err = stmt.Exec(q.Status, q.ModifiedAt, q.ModifiedOp, q.Id)
	if err != nil {
		log.Errorf("UpdateQuota Error: err=%s", err)
		return err
	}

	log.Infof("UpdateQuota: id=%d, name=%s, cpu=%d, mem=%d, rbd=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		q.Id, q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status, q.CreatedAt, q.ModifiedAt, q.ModifiedOp)
	return nil
}

func QueryAllQuotas() ([]Quota, error) {
	// New quotas pint array
	quotas := make([]Quota, 0)

	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(QUOTA_SELECT_ALL)
	if err != nil {
		log.Errorf("QueryQuota Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Errorf("QueryAllQuotas Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		q := new(Quota)

		var comment []byte
		err = rows.Scan(&q.Id, &q.Name, &q.Cpu, &q.Mem, &q.Rbd,
			&q.Price, &q.Status, &q.CreatedAt, &q.ModifiedAt, &q.ModifiedOp, &comment)

		q.Comment = string(comment)

		if err != nil {
			log.Errorf("QueryAllQuotas rows.Next() Error: err=%s", err)
			return nil, err
		}

		quotas = append(quotas, *q)
		log.Infof("QueryAllQuotas row.Next(): id=%d, name=%s, cpu=%d, mem=%d, rbd=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
			q.Id, q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status, q.CreatedAt, q.ModifiedAt, q.ModifiedOp)
	}

	return quotas, nil

}


func QueryAllQuotasOrderByCpu() ([]Quota, error) {
	// New quotas pint array
	quotas := make([]Quota, 0)

	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(QUOTA_SELECT_ALL_ORDER_BY_CPU)
	if err != nil {
		log.Errorf("QueryQuota Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Errorf("QueryAllQuotas Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		q := new(Quota)

		var comment []byte
		err = rows.Scan(&q.Id, &q.Name, &q.Cpu, &q.Mem, &q.Rbd,
			&q.Price, &q.Status, &q.CreatedAt, &q.ModifiedAt, &q.ModifiedOp, &comment)

		q.Comment = string(comment)

		if err != nil {
			log.Errorf("QueryAllQuotas rows.Next() Error: err=%s", err)
			return nil, err
		}

		quotas = append(quotas, *q)
		log.Infof("QueryAllQuotas row.Next(): id=%d, name=%s, cpu=%d, mem=%d, rbd=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
			q.Id, q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status, q.CreatedAt, q.ModifiedAt, q.ModifiedOp)
	}

	return quotas, nil

}

func (q *Quota) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), q)

	if err != nil {
		log.Errorf("DecodeJson Error: err=%s", err)
	}
}

func (q *Quota) EncodeJson() string {
	data, err := json.Marshal(q)
	if err != nil {
		log.Errorf("DecodeJson Error: err=%s", err)
		return ""
	}
	return string(data)
}
