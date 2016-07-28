package organization

import (
	localtime "app/backend/common/util/time"
	mysql "app/backend/common/util/mysql"
	"log"
	"fmt"
	"encoding/json"
)

const (
	ORG_SELECT = "SELECT id, name, cpuQuota, memQuota, budget, balance, createdAt, modifiedAt, modifiedOp, comment FROM organization where id=?"
	ORG_INSERT = ""
	ORG_UPDATE = ""
	ORG_DELETE = ""
	VALID = 1
	INVALID = 0
)

type Organization struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	CpuQuota   int32  `json:"cpu_quota"`
	MemQuota   int32  `json:"mem_quota"`
	Budget     string  `json:"buget"`
	Balance    string  `json:"balance"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
	ModifiedOp int32    `json:"modifiedOp"`
	Comment    string `json:"comment,omitempty"`
}


func NewOrganization(name, budget, balance, comment string, cpuQuota, memQuota, modifiedOp int32) *Organization {

	return &Organization{
		Name: name,
		CpuQuota: cpuQuota,
		MemQuota: memQuota,
		Budget: budget,
		Balance: balance,
		CreatedAt: localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment: comment,
	}
}

func (o *Organization) QueryOrganizationById(id int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(ORG_SELECT)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()


	// Query organization by id
	stmt.QueryRow(id).Scan(&o.Id, &o.Name, &o.CpuQuota, &o.MemQuota, &o.Budget, &o.Balance, &o.CreatedAt, &o.ModifiedAt, &o.ModifiedOp, &o.Comment)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	fmt.Printf("%v\n", o)
}

func (o *Organization) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), o)

	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}

func (o *Organization) EncodeJson() string {
	data, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	return string(data)
}
