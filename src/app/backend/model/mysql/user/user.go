package user

import (
	localtime "app/backend/common/util/time"
	mysql "app/backend/common/util/mysql"
	"encoding/json"
)


func NewUser(name, password, comment string, orgId, modifiedOp int32) *User {

	return &User{
		Name:       name,
		Password:   password,
		OrgId:      orgId,
		Status:     VALID,
		Comment:    comment,
		NavList:    USER_NAVLIST,
		ModifiedAt: localtime.NewLocalTime().String(),
		CreatedAt:  localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
	}
}

func QueryAllUsers() ([]User, error) {
	users := make([]User, 0)

	db := mysql.MysqlInstance().Conn()

	// Prepare select-all-statement
	stmt, err := db.Prepare(USER_SELECT_ALL)
	if err != nil {
		log.Errorf("QueryAllUsers Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Errorf("QueryAllUsers Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(User)

		var comment []byte
		var password string
		//TODO: didn't pass password back to frontend now
		err = rows.Scan(&u.Id, &u.Name, &password, &u.OrgId, &u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp, &comment)
		u.Comment = string(comment)

		if err != nil {
			log.Errorf("QueryAllUsers Error: err=%s", err)
			return nil, err
		}
		users = append(users, *u)

		log.Infof("QueryAllUsers: id=%d, name=%s, orgId=%d, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
			u.Id, u.Name, u.OrgId, u.Status, u.CreatedAt, u.ModifiedAt, u.ModifiedOp)

	}

	log.Infof("QueryAllUsers: len(users)=%d", len(users))
	return users, nil
}

func QueryUsersByOrgId(orgId int32) ([]User, error) {
	users := make([]User, 0)

	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(USER_SELECT_BY_ORGID)
	if err != nil {
		log.Errorf("QueryUsersByOrgId Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(orgId)
	if err != nil {
		log.Errorf("QueryUsersByOrgId Error: err=%s", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		u := new(User)
		var comment []byte
		err = rows.Scan(&u.Id, &u.Name, &u.Password, &u.OrgId, &u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp, &comment)
		if err != nil {
			log.Errorf("QueryUsersByOrgId Error: err=%s", err)
			return nil, err
		}
		u.Comment = string(comment)

		users = append(users, *u)
		log.Infof("QueryUsersByOrgId: id=%d, name=%s, orgId=%d, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
			u.Id, u.Name, u.OrgId, u.Status, u.CreatedAt, u.ModifiedAt, u.ModifiedOp)
	}
	log.Infof("QueryUsersByOrgId: len(users)=%d", len(users))
	return users, nil

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
		&u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp, &u.NavList)

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
		&u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp, &u.NavList)
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
		&u.CreatedAt, &u.ModifiedAt, &u.ModifiedOp, &u.NavList)
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
		u.CreatedAt, u.ModifiedAt, u.ModifiedOp, u.Comment, u.NavList)

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
	_, err = stmt.Exec(u.Password, u.OrgId, u.ModifiedAt, u.ModifiedOp, u.NavList, u.Id)

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
	u.NavList = USER_NAVLIST

	// Set user status  INVALID
	u.Status = INVALID
	_, err = stmt.Exec(u.Status, u.ModifiedAt, u.ModifiedOp, u.NavList, u.Id)
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
	log.Infof("queryUserNameByUserId successfully")
	return u.Name
}

/*
// Query NavList
func QueryNavListById(userId int32) (string, error) {
	db := mysql.MysqlInstance().Conn()

	// Preaper user-paswrod statement
	stmt, err := db.Prepare(USER_NAVLIST)
	if err != nil {
		log.Fatalf("QueryUserNavList Error: err=%s", err)
		return nil
	}
	defer stmt.Close()

	navList := ""

	// Query Id by name and password
	err = stmt.QueryRow(userId).Scan(&navList)

	if err != nil {
		log.Errorf("QueryUserNavListById Error: err=%s", err)
		return "", err
	}

	return navList, nil
}
*/