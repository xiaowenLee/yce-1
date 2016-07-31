package quota

import (
	localtime "app/backend/common/util/time"
	mysql "app/backend/common/util/mysql"
	"fmt"
	"log"
	"encoding/json"
)

const (

	QUOTA_SELECT = "SELECT id, name, cpu, mem, rbd, price, " +
		"status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM quota WHERE id=?"

	QUOTA_INSERT = "INSERT INTO " +
		"quota(name, cpu, mem, rbd, price, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	QUOTA_UPDATE = "UPDATE quota " +
		"SET name=?, cpu=?, mem=?, rbd=?, price=?, status=?, modifiedAt=?, modifiedOp=?, comment=? " +
		"WHERE id=?"

	QUOTA_DELETE = "UPDATE quota " +
		"SET status=?, modifiedAt=?, modifiedOp=? " +
		"WHERE id=?"

	VALID = 1
	INVALID = 0
)

type Quota struct {
	Id         int32   `json:"id"`
	Name       string  `json:"name"`
	Cpu        int32   `json:"cpu"`
	Mem        int32   `json:"mem"`
	Rbd        int32   `json:"rbd"`
	Price      string  `json:"price"`
	Status     int32   `json:"status"`
	CreatedAt  string  `json:"createdAt"`
	ModifiedAt string  `json:"modifiedAt"`
	ModifiedOp int     `json:"modifiedOp"`
	Comment    string  `json:"comment"`
}

func NewQuota(name, price, comment string, cpu, mem, rbd, modifiedOp int32) *Quota {
	return &Quota{
		Name: name,
		Cpu: cpu,
		Mem: mem,
		Rbd: rbd,
		Price: price,
		Status: VALID,
		CreatedAt: localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment: comment,
	}
}

func (q *Quota) QueryQuotaById(id int32) {
	db := mysql. MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(QUOTA_SELECT)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()

	// Query quota by id
	err = stmt.QueryRow(id).Scan(&q.Name, &q.Cpu, &q.Mem, &q.Rbd,
		&q.Price, &q.Status, &q.CreatedAt, &q.ModifiedAt, &q.ModifiedOp, q.Comment)


}

func (q *Quota) InsertQuota(op int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepared insert-statement
	stmt, err := db.Prepare(QUOTA_INSERT)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()

	// Update createAt, modifiedAt, modifiedOp
	q.CreatedAt = localtime.NewLocalTime().String()
	q.ModifiedAt = localtime.NewLocalTime().String()
	q.ModifiedOp = op

	// Insert a user
	_, err = stmt.Exec(q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status,
		q.CreatedAt, q.ModifiedAt, q.ModifiedOp, q.Comment, q.Id)

	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

}

func (q *Quota) UpdateQuota(op int32) {

	db := mysql.MysqlInstance().Conn()

	// Prepared update-statement
	stmt, err := db.Prepare(QUOTA_UPDATE)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	q.ModifiedAt = localtime.NewLocalTime().String()
	q.ModifiedOp = op

	// Update a quota
	_, err = stmt.Exec(q.Name, q.Cpu, q.Mem, q.Rbd, q.Price, q.Status, q.ModifiedAt, q.ModifiedOp, q.Comment, q.Id)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

}

func (q *Quota) DeleteQuota(op int32) {

	db := mysql.MysqlInstance().Conn()

	// Prepared delete-statement
	stmt, err := db.Prepare(QUOTA_DELETE)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	q.ModifiedAt = localtime.NewLocalTime().String()
	q.ModifiedOp = op

	// Update a quota
	_, err = stmt.Exec(q.Status, q.ModifiedAt, q.ModifiedOp, q.Id)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

}

func (q *Quota) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), q)

	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}

func (q *Quota) EncodeJson() string {
	data, err := json.MarshalIndent(q, "", " ")
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	return string(data)
}
