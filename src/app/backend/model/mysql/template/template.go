package template

import (
	localtime "app/backend/common/util/time"
	mysql "app/backend/common/util/mysql"
	"encoding/json"
)



func NewTemplate(name string, orgId int32, deployment, service, endpoints string, modifiedOp int32, comment string) (*Template){
	return &Template {
		Name: name,
		OrgId: orgId,
		Deployment: deployment,
		Service: service,
		Endpoints: endpoints,
		Status: VALID,
		CreatedAt: localtime.NewLocalTime().String(),
		ModifiedAt: localtime.NewLocalTime().String(),
		ModifiedOp: modifiedOp,
		Comment: comment,
	}
}

func QueryAllTemplatesByOrgId(orgId int32) ([]Template, error) {
	templates := make([]Template, 0)

	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(QUERY_ALL_BY_ORGID)
	if err != nil {
		log.Fatalf("QueryAllTemplateByOrgId Error: err=$s", err)
		return nil, nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(orgId, VALID)
	if err != nil {
		log.Errorf("QueryAllTemplateByOrgId Error: err=%s", err)
		return nil, nil
	}
	defer rows.Close()

	for rows.Next() {
		t := new(Template)

		var comment []byte
		var deployment []byte
		var service []byte
		var endpoints []byte

		err = rows.Scan(&t.Id, &t.Name, &t.OrgId, &deployment, &service, &endpoints, &t.Status, &t.CreatedAt, &t.ModifiedAt, &t.ModifiedOp, &comment)
		if err != nil {
			log.Errorf("QueryAllTemplateByOrgId Error: err=%s", err)
			return nil, nil
		}

		t.Deployment = string(deployment)
		t.Service = string(service)
		t.Endpoints = string(endpoints)
		t.Comment = string(comment)

		templates = append(templates, *t)

		//log.Infof("QueryAllTemplateByOrgId: id=%d, name=%s, orgId=%d, deployment=%s, service=%s, endpoints=%s, status=%d, createdAt=%s, modifiedAt=%s, modifiedOp=%d, comment=%s",
		//	t.Id, t.Name, t.OrgId, t.Deployment, t.Service, t.Endpoints, t.Status, t.CreatedAt, t.ModifiedAt, t.ModifiedOp, t.Comment)
	}

	log.Infof("QueryAllTemplateByOrgId: len(templates)=%d", len(templates))
	return templates, nil
}

func (t *Template) QueryTemplateByTemplateNameAndOrgId(name string, orgId int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(QUERY_DUPLICATED_NAME)
	if err != nil {
		log.Fatalf("QueryTemplateByTemplateNameAndOrgId Error: err=%s", err)
		return nil
	}

	err = stmt.QueryRow(name, orgId).Scan(&t.Id, &t.Name, &t.OrgId, &t.Deployment, &t.Service, &t.Endpoints, &t.Status, &t.CreatedAt, &t.ModifiedAt, &t.ModifiedOp, &t.Comment)
	if err != nil {
		log.Errorf("QueryTemplateByTemplateNameAndOrgId Error: err=%s", err)
		return err
	}
	return nil
}

func (t *Template) QueryTemplateById(id int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(QUERY_BY_ID)
	if err != nil {
		log.Fatalf("QueryTemplateById Error: err=%s", err)
		return nil
	}
	defer stmt.Close()

	err = stmt.QueryRow(id, VALID).Scan(&t.Id, &t.Name, &t.OrgId, &t.Deployment, &t.Service, &t.Endpoints, &t.Status, &t.CreatedAt, &t.ModifiedAt, &t.ModifiedOp, &t.Comment)
	if err != nil {
		log.Errorf("QueryTemplateById Error: err=%s", err)
		return err
	}

	return nil
}

func (t *Template) QueryTemplateByName(name string) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(QUERY_BY_NAME)
	if err != nil {
		log.Fatalf("QueryTemplateByName Error: err=%s", err)
		return nil
	}
	defer stmt.Close()

	err = stmt.QueryRow(name, VALID).Scan(&t.Id, &t.Name, &t.OrgId, &t.Deployment, &t.Service, &t.Endpoints, &t.Status, &t.CreatedAt, &t.ModifiedAt, &t.ModifiedOp, &t.Comment)
	if err != nil {
		log.Errorf("QueryTemplateByName Error: err=%s", err)
		return err
	}

	return nil
}

func (t *Template) InsertTemplate(op int32) error {
	db := mysql.MysqlInstance().Conn()

	// Prepare insert-statment
	//stmt, err := db.Prepare(TEMPLATE_INSERT)
	stmt, err := db.Prepare(TEMPLATE_INSERT_ON_DUPLICATE_UPDATE)
	if err != nil {
		log.Fatalf("InsertTemplate Error: error=%s", err)
	}
	defer stmt.Close()

	// Update CreatedAt, modifiedAt, modifiedOp
	t.CreatedAt = localtime.NewLocalTime().String()
	t.ModifiedAt = localtime.NewLocalTime().String()
	t.ModifiedOp = op

	// Insert
	_, err = stmt.Exec(t.Name, t.OrgId, t.Deployment, t.Service, t.Endpoints, t.Status, t.CreatedAt, t.ModifiedAt, t.ModifiedOp, t.Comment, VALID, t.Deployment, t.Service, t.Endpoints)
	if err != nil {
		log.Errorf("InsertTemplate Error: error=%s", err)
		return err
	}

	return nil
}

func (t *Template) UpdateTemplate(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(TEMPLATE_UPDATE)
	if err != nil {
		log.Fatalf("UpdateTemplate Error: err=%s", err)
		return nil
	}
	defer stmt.Close()

	//update modifiedAt, modifiedOp
	t.ModifiedAt = localtime.NewLocalTime().String()
	t.ModifiedOp = op

	// update a template
	_, err = stmt.Exec(t.Name, t.Deployment, t.Service, t.Endpoints, t.Id)
	if err != nil {
		log.Errorf("UpdateTemplate Error: err=%s", err)
		return err
	}

	return nil
}

func (t *Template) DeleteTemplate(op int32) error {
	db := mysql.MysqlInstance().Conn()

	stmt, err := db.Prepare(TEMPLATE_DELETE)
	if err != nil {
		log.Fatalf("DeleteTemplate Error: err=%s", err)
		return nil
	}

	// update status from VALID to INVALID
	t.Status = INVALID
	t.ModifiedAt = localtime.NewLocalTime().String()
	t.ModifiedOp = op

	// do
	_, err = stmt.Exec(t.Status, t.ModifiedAt, t.ModifiedOp, t.Id)
	if err != nil {
		log.Errorf("DeleteTemplate Error: err=%s", err)
		return err
	}
	return nil
}

func (t *Template) EncodeJson() (string ,error) {
	data, err := json.Marshal(t)
	if err != nil {
		log.Errorf("EncodeJson Error: err=%s", err)
		return "", err
	}
	return string(data), nil
}

func (t *Template) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), t)
	if err != nil {
		log.Errorf("DecodeJson Error: err=%s", err)
		return err
	}

	return nil
}




