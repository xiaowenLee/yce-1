package dcquota

import (
	"app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"fmt"
	"log"
	"os/user"
	"encoding/json"
)

const (
	DCQUOTA_SELECT = "SELECT id, dcId, orgId, podNumList, podCpuMax, podMemMax, podCpuMin, podMemMin, rbdQuota, podRbdMax, podRbdMin, price, status, createdAt, modifiedAt, modifiedOp, comment FROM dcquota WHERE id=?"
	DCQUOTA_INSERT = "INSERT INTO dcquota(dcId, orgId, podNumList, podCpuMax, podMemMax, podCpuMin, podMemMin, rbdQuota, podRbdMax, podRbdMin, price, status, createdAt, modifiedAt, modifiedOp, comment) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	DCQUOTA_UPDATE = "UPDATE dcquota SET dcId=?, orgId=?, podNumList=?, podCpuMax=?, podMemMax=?, podCpuMin=?, podMemMin=?, rbdQuota=?, podRbdMax=?, podRbdMin=?, price=?, status=?, modifiedAt=?, modifiedOp=?, comment=? WHERE id=?"
	DCQUOTA_DELETE = "UPDATE dcquota SET status=?, modifiedAt=?, modifiedOp=? WHERE id=?"

	VALID   = 1
	INVALID = 0
)

type DcQuota struct {
	Id          int32  `json:"id"`
	DcId        int32  `json:"dcId"`
	OrgId       int32  `json:"orgId"`
	PodNumLimit int32  `json:"podNumLimit"`
	PodCpuMax   int32  `json:"podCpuMax"`
	PodMemMax   int32  `json:"podMemMax"`
	PodCpuMin   int32  `json:"podCpuMin"`
	PodMemMin   int32  `json:"podMemMin"`
	RbdQuota    int32  `json:"rbdQuota"`
	PodRbdMax   int32  `json:"podRbdMax"`
	PodRbdMin   int32  `json:"podRbdMin"`
	Price       string `json:"price"`
	Status      int32  `json:"status"`
	CreatedAt   string `json:"createdAt`
	ModifiedAt  string `json:"modifiedAt"`
	ModifiedOp  int32  `json:"modifiedOp"`
	Comment     string `json:"comment"`
}

func NewDcQuota(id, dcId, orgId, podNumLimit, podCpuMax, podMemMax, podCpuMin, podMemMin, rbdQuota, PodRbdMax, podRbdMin, modifiedOp int32, price, comment string) *DcQuota {
	return &DcQuota{
		Id:          id,
		DcId:        dcId,
		OrgId:       orgId,
		PodNumLimit: podNumLimit,
		PodCpuMax:   podCpuMax,
		PodMemMax:   podMemMax,
		PodCpuMin:   podCpuMin,
		PodMemMin:   podMemMin,
		RbdQuota:    rbdQuota,
		PodRbdMax:   PodRbdMax,
		PodRbdMin:   podRbdMin,
		Price:       price,
		Status:      VALID,
		CreatedAt:   localtime.NewLocalTime().String(),
		ModifiedAt:  localtime.NewLocalTime().String(),
		ModifiedOp:  modifiedOp,
		Comment:     comment,
	}
}

func (dc *DcQuota) QueryDcQuotaById(id int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(DCQUOTA_SELECT)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()

	// Query DcQuota by id
	err = stmt.QueryRow(id).Scan(&dc.Id, &dc.DcId, &dc.OrgId, &dc.PodNumLimit, &dc.PodCpuMax, &dc.PodMemMax, &dc.PodCpuMin, &dc.PodMemMin, &dc.RbdQuota, &dc.PodRbdMax, &dc.PodRbdMin, &dc.Price, &dc.Status, &dc.CreatedAt, &dc.ModifiedAt, &dc.ModifiedOp, &dc.Comment)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	fmt.Printf("%v\n", dc)
}

func (dc *DcQuota) InsertDcQuota(op int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(DCQUOTA_INSERT)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()

	// Update createdAt, modifiedAt, modifiedOp
	dc.CreatedAt = localtime.NewLocalTime().String()
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op

	// Insert a dcQuota
	_, err = stmt.Exec(dc.DcId, dc.OrgId, dc.PodNumLimit, dc.PodCpuMax, dc.PodMemMax, dc.PodCpuMin, dc.PodMemMin, dc.RbdQuota, dc.PodRbdMax, dc.PodRbdMin, dc.Price, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp, dc.Comment)

	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}

func (dc *DcQuota) UpdateDcQuota(op int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepared update-statement
	stmt, err := db.Prepare(DCQUOTA_UPDATE)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op

	// Update a dcQuota
	_, err = stmt.Exec(dc.DcId, dc.OrgId, dc.PodNumLimit, dc.PodCpuMax, dc.PodMemMax, dc.PodCpuMin, dc.PodMemMin, dc.RbdQuota, dc.PodRbdMax, dc.PodRbdMin, dc.Price, dc.Status, dc.ModifiedAt, dc.ModifiedOp, dc.Comment, dc.Id)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}

func (dc *DcQuota) DeleteDcQuota(op int32) {
	db := mysql.MysqlInstance().Conn()

	// Prepared delet-statement
	stmt, err := db.Prepare(DCQUOTA_DELETE)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	defer stmt.Close()

	// Update modifiedAt and modifiedOp
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op

	// Set user status INVALID
	dc.Status = INVALID
	_, err = stmt.Exec(dc.Status, dc.ModifiedAt, dc.ModifiedOp, dc.Id)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}


func (dc *DcQuota) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), dc)

	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
}

func (dc *DcQuota) EncodeJson() string {
	data, err := json.MarshalIndent(dc, "", " ")
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}
	return string(data)
}
