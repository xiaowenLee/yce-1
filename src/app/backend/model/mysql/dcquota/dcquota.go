package dcquota

import (
	"app/backend/common/util/mysql"
	localtime "app/backend/common/util/time"
	"encoding/json"
	"log"
)

const (
	DCQUOTA_SELECT = "SELECT id, dcId, orgId, podNumLimit, podCpuMax, podMemMax, podCpuMin, " +
		"podMemMin, rbdQuota, podRbdMax, podRbdMin, " +
		"price, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM dcquota WHERE id=?"

	DCQUOTA_INSERT = "INSERT INTO dcquota(dcId, orgId, podNumLimit, podCpuMax, podMemMax, " +
		"podCpuMin, podMemMin, rbdQuota, podRbdMax, podRbdMin, " +
		"price, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	DCQUOTA_UPDATE = "UPDATE dcquota SET dcId=?, orgId=?, podNumLimit=?, podCpuMax=?, " +
		"podMemMax=?, podCpuMin=?, podMemMin=?, rbdQuota=?, podRbdMax=?, podRbdMin=?, " +
		"price=?, status=?, modifiedAt=?, modifiedOp=?, comment=? WHERE id=?"

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

func NewDcQuota(dcId, orgId, podNumLimit, podCpuMax, podMemMax, podCpuMin, podMemMin, rbdQuota, PodRbdMax, podRbdMin, modifiedOp int32, price, comment string) *DcQuota {
	return &DcQuota{
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

func (dc *DcQuota) QueryDcQuotaById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare select-statement
	stmt, err := db.Prepare(DCQUOTA_SELECT)
	if err != nil {
		log.Fatal("QueryDcQuotaById Error: err=%s\n", err)
		return err
	}
	defer stmt.Close()

	// Query DcQuota by id
	var comment []byte
	err = stmt.QueryRow(id).Scan(&dc.Id, &dc.DcId, &dc.OrgId, &dc.PodNumLimit, &dc.PodCpuMax,
		&dc.PodMemMax, &dc.PodCpuMin, &dc.PodMemMin, &dc.RbdQuota, &dc.PodRbdMax,
		&dc.PodRbdMin, &dc.Price, &dc.Status, &dc.CreatedAt, &dc.ModifiedAt, &dc.ModifiedOp, &comment)

	dc.Comment = string(comment)

	if err != nil {
		log.Printf("QueryDcQuotaById Error: err=%s\n", err)
		return err
	}

	log.Printf("QueryDcQuotaById: id=%d, dcId=%d, orgId=%d, podNumLimit=%d, podCpuMax=%d, podMemMax=%d, podCpuMin=%d, podMemMin=%d, rbdQuota=%d, podRbdMax=%d, podRbdMin=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		dc.Id, dc.DcId, dc.OrgId, dc.PodNumLimit, dc.PodCpuMax, dc.PodMemMax, dc.PodCpuMin, dc.PodMemMin, dc.PodRbdMax, dc.RbdQuota, dc.PodRbdMin, dc.Price, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil
}

func (dc *DcQuota) InsertDcQuota(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statement
	stmt, err := db.Prepare(DCQUOTA_INSERT)
	if err != nil {
		log.Printf("InsertDcQuota Error: err=%s\n", err)
		return err
	}
	defer stmt.Close()

	// Update createdAt, modifiedAt, modifiedOp
	dc.CreatedAt = localtime.NewLocalTime().String()
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op

	// Insert a dcQuota
	_, err = stmt.Exec(dc.DcId, dc.OrgId, dc.PodNumLimit, dc.PodCpuMax, dc.PodMemMax, dc.PodCpuMin, dc.PodMemMin, dc.RbdQuota, dc.PodRbdMax, dc.PodRbdMin, dc.Price, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp, dc.Comment)

	if err != nil {
		log.Printf("InsertDcQuota Error: err=%s\n", err)
		return err
	}

	log.Printf("InsertDcQuotaById: id=%d, dcId=%d, orgId=%d, podNumLimit=%d, podCpuMax=%d, podMemMax=%d, podCpuMin=%d, podMemMin=%d, rbdQuota=%d, podRbdMax=%d, podRbdMin=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		dc.Id, dc.DcId, dc.OrgId, dc.PodNumLimit, dc.PodCpuMax, dc.PodMemMax, dc.PodCpuMin, dc.PodMemMin, dc.PodRbdMax, dc.RbdQuota, dc.PodRbdMin, dc.Price, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil
}

func (dc *DcQuota) UpdateDcQuota(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepared update-statement
	stmt, err := db.Prepare(DCQUOTA_UPDATE)
	if err != nil {
		log.Printf("UpdateDcQuota Error: err=%s\n", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt, modifiedOp
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op

	// Update a dcQuota
	_, err = stmt.Exec(dc.DcId, dc.OrgId, dc.PodNumLimit, dc.PodCpuMax, dc.PodMemMax, dc.PodCpuMin, dc.PodMemMin, dc.RbdQuota, dc.PodRbdMax, dc.PodRbdMin, dc.Price, dc.Status, dc.ModifiedAt, dc.ModifiedOp, dc.Comment, dc.Id)
	if err != nil {
		log.Printf("UpdateDcQuota Error: err=%s\n", err)
		return err
	}

	log.Printf("UpdateDcQuotaById: id=%d, dcId=%d, orgId=%d, podNumLimit=%d, podCpuMax=%d, podMemMax=%d, podCpuMin=%d, podMemMin=%d, rbdQuota=%d, podRbdMax=%d, podRbdMin=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		dc.Id, dc.DcId, dc.OrgId, dc.PodNumLimit, dc.PodCpuMax, dc.PodMemMax, dc.PodCpuMin, dc.PodMemMin, dc.PodRbdMax, dc.RbdQuota, dc.PodRbdMin, dc.Price, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil
}

func (dc *DcQuota) DeleteDcQuota(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepared delet-statement
	stmt, err := db.Prepare(DCQUOTA_DELETE)
	if err != nil {
		log.Printf("DeleteDcQuota Error: err=%s\n", err)
		return err
	}
	defer stmt.Close()

	// Update modifiedAt and modifiedOp
	dc.ModifiedAt = localtime.NewLocalTime().String()
	dc.ModifiedOp = op

	// Set user status INVALID
	dc.Status = INVALID
	_, err = stmt.Exec(dc.Status, dc.ModifiedAt, dc.ModifiedOp, dc.Id)
	if err != nil {
		log.Printf("DeleteDcQuota Error: err=%s\n", err)
		return err
	}

	log.Printf("DeleteDcQuotaById: id=%d, dcId=%d, orgId=%d, podNumLimit=%d, podCpuMax=%d, podMemMax=%d, podCpuMin=%d, podMemMin=%d, rbdQuota=%d, podRbdMax=%d, podRbdMin=%d, price=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d\n",
		dc.Id, dc.DcId, dc.OrgId, dc.PodNumLimit, dc.PodCpuMax, dc.PodMemMax, dc.PodCpuMin, dc.PodMemMin, dc.PodRbdMax, dc.RbdQuota, dc.PodRbdMin, dc.Price, dc.Status, dc.CreatedAt, dc.ModifiedAt, dc.ModifiedOp)
	return nil
}

func (dc *DcQuota) DecodeJson(data string) {
	err := json.Unmarshal([]byte(data), dc)

	if err != nil {
		log.Println("DecodeJson Error: err=%s\n", err)
	}
}

func (dc *DcQuota) EncodeJson() string {
	data, err := json.Marshal(dc)
	if err != nil {
		log.Println("EncodeJson Error: err=%s\n", err)
		return ""
	}
	return string(data)
}
