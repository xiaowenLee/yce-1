package organization

import (
	"log"
	"encoding/json"
	mysql "app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"github.com/shopspring/decimal"
)

const (
	ORG_SELECT = "SELECT id, name, cpuQuota, memQuota, budget, balance, status, " +
		"createdAt, modifiedAt, modifiedOp, comment " +
		"FROM organization WHERE id=?"

	ORG_INSERT = "INSERT INTO organization(name, cpuQuota, memQuota, budget, " +
		"balance, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	ORG_UPDATE = "UPDATE organization SET name=?, cpuQuota=?, memQuota=?, budget=?, " +
		"balance=?, status=?, modifiedAt=?, modifiedOp=?, comment=? " +
		"WHERE id=?"

	ORG_DELETE = "UPDATE organization SET status=?, modifiedAt=?, modifiedOp=? WHERE id=?"

	VALID   = 1
	INVALID = 0
)

type Organization struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	CpuQuota   int32  `json:"cpu_quota"`
	MemQuota   int32  `json:"mem_quota"`
	Budget     string `json:"buget"`
	Balance    string `json:"balance"`
	Status     int32  `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32  `json:"modifiedOp"`
	Comment    string `json:"comment,omitempty"`
}

func NewOrganization(name, budget, balance, comment string, cpuQuota, memQuota, modifiedOp int32) *Organization {

	return &Organization{
		Name:       name,
		CpuQuota:   cpuQuota,
		MemQuota:   memQuota,
		Budget:     budget,
		Balance:    balance,
		Status:     VALID,
		CreatedAt:  localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment:    comment,
	}
}

func (o *Organization) QueryOrganizationById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(ORG_SELECT)
	if err != nil {
		log.Printf("QueryOrganizationById Error: err=%s\n", err)
		return err
	}
	defer stmt.Close()

	// Query organization by id
	err = stmt.QueryRow(id).Scan(&o.Id, &o.Name, &o.CpuQuota, &o.MemQuota, &o.Budget, &o.Balance, &o.Status, &o.CreatedAt, &o.ModifiedAt, &o.ModifiedOp, &o.Comment)
	if err != nil {
		log.Printf("QureyOrganizationById Error: err=%s\n", err)
		return err
	}

	log.Printf("QueryOrganizationById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return nil
}

func (o *Organization) QueryBudgetById(id int32) (budget decimal.Decimal, err error) {
	o.QueryOrganizationById(id)
	budget, err = decimal.NewFromString(o.Budget)

	log.Printf("QueryBudgetById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return budget, err
}

func (o *Organization) QueryBalanceById(id int32) (balance decimal.Decimal, err error) {
	o.QueryOrganizationById(id)
	balance, err = decimal.NewFromString(o.Balance)

	log.Printf("QueryBalanceById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return balance, err
}

func (o *Organization) InsertOrganization(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(ORG_INSERT)
	if err != nil {
		log.Printf("InsertOrganization Error: err=%s\n", err)
		return err
	}
	defer stmt.Close()

	// Update createdAt, modifiedAt, modifiedOp
	o.CreatedAt = localtime.NewLocalTime().String()
	o.ModifiedAt = localtime.NewLocalTime().String()
	o.ModifiedOp = op

	// Insert a organization
	_, err = stmt.Exec(o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp, o.Comment)
	if err != nil {
		log.Printf("InsertOrganization Error: err=%s\n", err)
		return err
	}

	log.Printf("InsertOrganization: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return nil
}

func (o *Organization) UpdateOrganization(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare update-statement
	stmt, err := db.Prepare(ORG_UPDATE)
	if err != nil {
		log.Printf("UpdateOrganization Error: err=%s\n", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt and modifiedOp
	o.ModifiedAt = localtime.NewLocalTime().String()
	o.ModifiedOp = op

	// Update a org: name, cpuQuota, memQuota, budget, balance, status, modifiedAt, modifiedOp, comment
	_, err = stmt.Exec(o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.ModifiedAt, o.ModifiedOp, o.Comment, o.Id)

	if err != nil {
		log.Printf("UpdateOrganization Error: err=%s\n", err)
		return err
	}

	log.Printf("UpdateOrganization: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)

	return nil
}

func (o *Organization) UpdateBudgetById(budget string, op int32) {
	o.Budget = budget

	log.Printf("UpdateBudgetById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	o.UpdateOrganization(op)
}

func (o *Organization) UpdateBalanceById(balance string, op int32) {
	o.Balance = balance

	log.Printf("UpdateBudgetById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	o.UpdateOrganization(op)
}

func (o *Organization) DeleteOrganization(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepared delete-statement
	stmt, err := db.Prepare(ORG_DELETE)
	if err != nil {
		log.Printf("DeleteOrganization Error: err=%s\n", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt and modifiedOp
	o.ModifiedAt = localtime.NewLocalTime().String()
	o.ModifiedOp = op
	o.Status = INVALID

	// Delete a org
	_, err = stmt.Exec(o.Status, o.ModifiedAt, o.ModifiedOp, o.Id)
	if err != nil {
		log.Printf("DeleteOrganization Error: err=%s\n", err)
		return err
	}

	log.Printf("DeleteBudgetById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return nil
}

func (o *Organization) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), o)

	if err != nil {
		log.Printf("DecodeJson Error: err=%s\n", err)
	}
}

func (o *Organization) EncodeJson() string {
	data, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		log.Printf("DecodeJson Erro: err=%s\n", err)
		return ""
	}
	return string(data)
}
