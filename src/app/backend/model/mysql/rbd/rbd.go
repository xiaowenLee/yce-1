package rbd

import (
	mylog "app/backend/common/util/log"
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"encoding/json"
)

var log = mylog.Log

const (
	RBD_SELECT = "SELECT id, name, pool, size, filesystem, " +
		"orgId, dcId, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM rbd WHERE id=?"

	RBD_INSERT = "INSERT INTO rbd(name, pool, size, filesystem, " +
		"orgId, dcId, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	RBD_UPDATE = "UPDATE rbd SET name=?, pool=?, size=?, filesystem=?, " +
		"orgId=?, dcId=?, status=?, createdAt=?, modifiedAt=?, modifiedOp=?, comment=? " +
		"WHERE id=?"

	RBD_DELETE = "UPDATE rbd SET status=?, modifiedAt=?, modifiedOp=? WHERE id=?"

	DEFAULT_FILESYSTEM = "ext4"

	VALID   = 1
	INVALID = 0
)

type Rbd struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	Pool       string `json:"pool"`
	Size       int32  `json:"size"`
	FileSystem string `jsonN:"filesystem"`
	OrgId      int32  `json:"orgId"`
	DcID       int32  `json:"dcId"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

func NewRbd(name, pool, comment string, size, orgId, dcId, modifiedOp int32) *Rbd {
	return &Rbd{
		Name:       name,
		Pool:       pool,
		Size:       size,
		FileSystem: DEFAULT_FILESYSTEM,
		OrgId:      orgId,
		DcID:       dcId,
		Status:     VALID,
		CreatedAt:  localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment:    comment,
	}
}

func (r *Rbd) QueryRbdById(id int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(RBD_SELECT)
	if err != nil {
		log.Fatalf("Rbd QureyRbdById Error: err=%s", err)
	}
	defer stmt.Close()

	// Query quota by id
	err = stmt.QueryRow(id).Scan(&r.Id, &r.Name, &r.Pool, &r.Size, &r.FileSystem,
		&r.OrgId, &r.DcID, &r.Status, &r.CreatedAt, &r.ModifiedAt, &r.ModifiedOp, &r.Comment)

	if err != nil {
		log.Errorf("Rbd QureyRbdById Error: err=%s", err)
	}

}

func (r *Rbd) InsertRbd(op int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepared insert-statement
	stmt, err := db.Prepare(RBD_INSERT)
	if err != nil {
		log.Fatalf("Rbd InsertRbd Error: err=%s", err)
	}
	defer stmt.Close()

	// Update createAt, modifiedAt, modifiedOp
	r.CreatedAt = localtime.NewLocalTime().String()
	r.ModifiedAt = localtime.NewLocalTime().String()
	r.ModifiedOp = op

	// Insert a user
	_, err = stmt.Exec(r.Name, r.Pool, r.Size, r.FileSystem, r.OrgId,
		r.DcID, r.Status, r.CreatedAt, r.ModifiedAt, r.ModifiedOp, r.Comment)

	if err != nil {
		log.Errorf("Rbd InsertRbd Error: err=%s", err)
	}

}

func (r *Rbd) UpdateRbd(op int32) {

	db := mysql.MysqlInstance().Conn()

	// Prepared update-statement
	stmt, err := db.Prepare(RBD_UPDATE)
	if err != nil {
		log.Fatalf("Rbd UpdateRbd Error: err=%s", err)
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	r.ModifiedAt = localtime.NewLocalTime().String()
	r.ModifiedOp = op

	// Update a quota
	_, err = stmt.Exec(r.Name, r.Pool, r.Size, r.FileSystem, r.OrgId, r.DcID,
		r.Status, r.CreatedAt, r.ModifiedAt, r.ModifiedOp, r.Comment, r.Id)

	if err != nil {
		log.Errorf("Rbd UpdateRbd Error: err=%s", err)
	}

}

func (r *Rbd) DeleteRbd(op int32) {

	db := mysql.MysqlInstance().Conn()

	// Prepared delete-statement
	stmt, err := db.Prepare(RBD_DELETE)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	r.ModifiedAt = localtime.NewLocalTime().String()
	r.ModifiedOp = op
	r.Status = INVALID

	// Update a quota
	_, err = stmt.Exec(r.Status, r.ModifiedAt, r.ModifiedOp, r.Id)
	if err != nil {
		log.Fatal(err)
	}

}

func (r *Rbd) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), r)

	if err != nil {
		log.Fatal(err)
	}
}

func (r *Rbd) EncodeJson() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
