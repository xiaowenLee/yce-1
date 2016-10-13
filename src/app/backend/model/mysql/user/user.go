package user

import (
	mylog "app/backend/common/util/log"
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"encoding/json"
)

var log = mylog.Log

const (
	USER_PASSWORD = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp FROM user WHERE name=? and password=?"

	USER_SELECT = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp FROM user WHERE id=? "
	USER_CHECK_DUPLICATE_NAME = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp FROM user WHERE name=? "
	USER_CHECK_DUPLICATE = "SELECT id, name, password, orgId, createdAt, modifiedAt, modifiedOp FROM user WHERE name=? AND orgId=?"

	USER_INSERT = "INSERT INTO " +
		"user(name, password, orgId, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	USER_UPDATE = "UPDATE user SET password=?, orgId=?, modifiedAt=?, modifiedOp=? WHERE id=?"
	USER_DELETE = "UPDATE user SET status=?, modifiedAt=?, modifiedOp=? WHERE id=?"
	VALID       = 1
	INVALID     = 0
)

type User struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	OrgId      int32  `json:"orgId"`
	Password   string `json:"password"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment"`
}

func NewUser(name, password, comment string, orgId, status, modifiedOp int32) *User {

	return &User{
		Name:       name,
		Password:   password,
		OrgId:      orgId,
		Status:     status,
		Comment:    comment,
		ModifiedAt: localtime.NewLocalTime().String(),
		CreatedAt:  localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
	}
}

func (u *User) QueryUserByNameAndPassword(name, password string) error {
	db := mysql.MysqlInstance().Conn()

	// Preaper user-paswrod statement
	stmt, err := db.Prepare(USER_PASSWORD)
	if err != nil {
		log.Fatalf("QueryUserByNameAndPassword Error: err=%s", err)
		return nil
	}
	defer stmt.Close()

	// Query Id by name and password
	err = stmt.QueryRow(name, password).Scan(&u.Id, &u.Name, &u.Password, &u.OrgId,
		&u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp)

	if err != nil {
		log.Errorf("QueryUserByNameAndPassword Error: err=%s", err)
		return err
	}

	return nil
}

func (u *User) QueryUserById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(USER_SELECT)
	if err != nil {
		log.Fatalf("QueryUserById Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Query user by id
	err = stmt.QueryRow(id).Scan(&u.Id, &u.Name, &u.Password, &u.OrgId,
		&u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp)
	if err != nil {
		log.Errorf("QueryUserById Error: err=%s", err)
		return err
	}

	return nil
}

func (u *User) QueryUserByUserName(name string) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(USER_CHECK_DUPLICATE_NAME)
	if err != nil {
		log.Fatalf("QueryUserByUserName Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Query user by name
	err = stmt.QueryRow(name).Scan(&u.Id, &u.Name, &u.Password, &u.OrgId,
		&u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp)
	if err != nil {
		log.Errorf("QueryUserByName Error: err=%s", err)
		return err
	}

	return nil
}

func (u *User) QueryUserByNameAndOrgId(name string, orgId int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(USER_CHECK_DUPLICATE)
	if err != nil {
		log.Fatalf("QueryUserByNameAndOrgId Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Query user by name and orgId
	err = stmt.QueryRow(name, orgId).Scan(&u.Id, &u.Name, &u.Password, &u.OrgId,
		&u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp)
	if err != nil {
		log.Errorf("QueryUserByNameAndOrgId Error: err=%s", err)
		return err
	}

	return nil
}

func (u *User) InsertUser(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(USER_INSERT)
	if err != nil {
		log.Fatalf("InsertUser Error: err=%s", err)
	}
	defer stmt.Close()

	// Update createdAt, modifiedAt, modifiedOp
	u.CreatedAt = localtime.NewLocalTime().String()
	u.ModifiedAt = localtime.NewLocalTime().String()
	u.ModifiedOp = op

	// Insert a user
	_, err = stmt.Exec(u.Name, u.Password, u.OrgId, u.Status,
		u.CreatedAt, u.ModifiedAt, u.ModifiedOp, u.Comment)

	if err != nil {
		log.Errorf("InsertUser Error: err=%s", err)
		return nil
	}

	return nil
}

func (u *User) UpdateUser(op int32) error {

	db := mysql.MysqlInstance().Conn()

	// Prepare update-statement
	stmt, err := db.Prepare(USER_UPDATE)
	if err != nil {
		log.Errorf("UpdateUser Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	u.ModifiedAt = localtime.NewLocalTime().String()
	u.ModifiedOp = op

	// Update a user: password or orgId
	_, err = stmt.Exec(u.Password, u.OrgId, u.ModifiedAt, u.ModifiedOp, u.Id)
	if err != nil {
		log.Errorf("UpdateUser Error: err=%s", err)
		return err
	}

	return nil
}

func (u *User) DeleteUser(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare delete-statement
	stmt, err := db.Prepare(USER_DELETE)
	if err != nil {
		log.Errorf("DeleteUser Error: err=%s", err)
		return err
	}

	defer stmt.Close()

	// Update modifiedAt and modifiedOp
	u.ModifiedAt = localtime.NewLocalTime().String()
	u.ModifiedOp = op

	// Set user status  INVALID
	u.Status = INVALID
	_, err = stmt.Exec(u.Status, u.ModifiedAt, u.ModifiedOp, u.Id)
	if err != nil {
		log.Errorf("DeleteUser Error: err=%s", err)
		return err
	}

	return nil
}

func (u *User) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), u)

	if err != nil {
		log.Errorf("DecodeJson Error: err=%s", err)
		return err
	}

	return nil
}

func (u *User) EncodeJson() (string, error) {
	data, err := json.Marshal(u)
	if err != nil {
		log.Errorf("EncodeJson Error: err=%s", err)
		return "", err
	}
	return string(data), nil
}

// Query UserName by UserId
func QueryUserNameByUserId(userId int32) (name string) {
	u := new(User)
	u.QueryUserById(userId)
	mylog.Log.Infof("queryUserNameByUserId successfully")
	return u.Name
}
