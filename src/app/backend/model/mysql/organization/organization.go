package organization

import (
	localtime "app/backend/common/util/time"
	"encoding/json"
	"github.com/shopspring/decimal"
	mysql "app/backend/common/util/mysql"
)






func NewOrganization(name, budget, balance, comment, dcIdList string, cpuQuota, memQuota, modifiedOp int32) *Organization {

	return &Organization{
		Name:       name,
		CpuQuota:   cpuQuota,
		MemQuota:   memQuota,
		Budget:     budget,
		Balance:    balance,
		Status:     VALID,
		DcIdList:   dcIdList,
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
		log.Errorf("QueryOrganizationById Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	var comment []byte
	// Query organization by id
	err = stmt.QueryRow(id).Scan(&o.Id, &o.Name, &o.CpuQuota, &o.MemQuota, &o.Budget, &o.Balance, &o.Status, &o.DcIdList, &o.CreatedAt, &o.ModifiedAt, &o.ModifiedOp, &comment)
	if err != nil {
		log.Errorf("QureyOrganizationById Error: err=%s", err)
		return err
	}

	o.Comment = string(comment)

	log.Infof("QueryOrganizationById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return nil
}

func QueryAllOrganizations() ([]Organization, error) {
	organizations := make([]Organization, 0)

	db := mysql.MysqlInstance().Conn()

	// Prepare select-all-statement
	stmt, err := db.Prepare(ORG_SELECT_ALL)
	if err != nil {
		log.Errorf("QueryAllOrganizations Error: err=%s", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(VALID)
	if err != nil {
		log.Errorf("QueryAllOrganizations Error: err=%s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		o := new(Organization)

		var comment []byte
		err = rows.Scan(&o.Id, &o.Name, &o.CpuQuota, &o.MemQuota, &o.Budget, &o.Balance,
			&o.Status, &o.DcIdList, &o.CreatedAt, &o.ModifiedAt, &o.ModifiedOp, &comment)
		o.Comment = string(comment)
		if err != nil {
			log.Errorf("QueryAllOrganizations Error: err=%s", err)
			return nil, err
		}
		organizations = append(organizations, *o)

		log.Infof("QueryAllOrganizations: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
			o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)

	}

	log.Infof("QueryAllOrganizations: len(organizations)=%d", len(organizations))
	return organizations, nil
}

func (o *Organization) QueryOrganizationByName(name string) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(ORG_SELECT_NAME)
	if err != nil {
		log.Errorf("QueryOrganizationByName Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	var comment []byte
	// Query organization by name
	err = stmt.QueryRow(name, VALID).Scan(&o.Id, &o.Name, &o.CpuQuota, &o.MemQuota, &o.Budget, &o.Balance, &o.Status, &o.DcIdList, &o.CreatedAt, &o.ModifiedAt, &o.ModifiedOp, &comment)
	if err != nil {
		log.Errorf("QureyOrganizationByName Error: err=%s", err)
		return err
	}

	o.Comment = string(comment)

	log.Infof("QueryOrganizationByName: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return nil
}

func (o *Organization) QueryBudgetById(id int32) (budget decimal.Decimal, err error) {
	o.QueryOrganizationById(id)
	budget, err = decimal.NewFromString(o.Budget)

	log.Infof("QueryBudgetById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return budget, err
}

func (o *Organization) QueryBalanceById(id int32) (balance decimal.Decimal, err error) {
	o.QueryOrganizationById(id)
	balance, err = decimal.NewFromString(o.Balance)

	log.Infof("QueryBalanceById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return balance, err
}

func (o *Organization) InsertOrganization() error {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(ORG_INSERT_ON_DUPLICATE_KEY_UPDATE)
	if err != nil {
		log.Errorf("InsertOrganization Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update createdAt, modifiedAt, modifiedOp
	o.CreatedAt = localtime.NewLocalTime().String()
	o.ModifiedAt = localtime.NewLocalTime().String()

	// Insert a organization
	_, err = stmt.Exec(o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp, o.Comment, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, VALID, o.DcIdList)
	if err != nil {
		log.Errorf("InsertOrganization Error: err=%s", err)
		return err
	}

	log.Infof("InsertOrganization: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return nil
}

func (o *Organization) UpdateOrganization(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare update-statement
	stmt, err := db.Prepare(ORG_UPDATE)
	if err != nil {
		log.Errorf("UpdateOrganization Error: err=%s", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt and modifiedOp
	o.ModifiedAt = localtime.NewLocalTime().String()
	o.ModifiedOp = op

	// Update a org: name, cpuQuota, memQuota, budget, balance, status, modifiedAt, modifiedOp, comment
	_, err = stmt.Exec(o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.ModifiedAt, o.ModifiedOp, o.Comment, o.Id)

	if err != nil {
		log.Errorf("UpdateOrganization Error: err=%s", err)
		return err
	}

	log.Infof("UpdateOrganization: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)

	return nil
}

func (o *Organization) UpdateBudgetById(budget string, op int32) {
	o.Budget = budget

	log.Infof("UpdateBudgetById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	o.UpdateOrganization(op)
}

func (o *Organization) UpdateBalanceById(balance string, op int32) {
	o.Balance = balance

	log.Infof("UpdateBudgetById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	o.UpdateOrganization(op)
}

func (o *Organization) UpdateQuotaById(quota *QuotaType, op int32) {
	o.CpuQuota = quota.CpuQuota
	o.MemQuota = quota.MemQuota

	log.Infof("UpdateQuotaById: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	o.UpdateOrganization(op)
}

func (o *Organization) DeleteOrganization(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepared delete-statement
	stmt, err := db.Prepare(ORG_DELETE)
	if err != nil {
		log.Errorf("DeleteOrganization Error: err=%s", err)
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
		log.Errorf("DeleteOrganization Error: err=%s", err)
		return err
	}

	log.Infof("DeleteOrganization: id=%d, name=%s, cpuQuota=%d, memQuota=%d, budget=%s, balance=%s, status=%d, dcIdList=%s, createdAt=%s, modifiedAt=%s, modifiedOp=%d",
		o.Id, o.Name, o.CpuQuota, o.MemQuota, o.Budget, o.Balance, o.Status, o.DcIdList, o.CreatedAt, o.ModifiedAt, o.ModifiedOp)
	return nil
}

func (o *Organization) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), o)

	if err != nil {
		log.Errorf("DecodeJson Error: err=%s", err)
	}
}

func (o *Organization) EncodeJson() string {
	data, err := json.Marshal(o)
	if err != nil {
		log.Errorf("EncodeJson Erro: err=%s", err)
		return ""
	}
	return string(data)
}
