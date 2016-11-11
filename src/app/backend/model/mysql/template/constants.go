package template

const (
	VALID = 1
	INVALID = 0

	TEMPLATE_INSERT = "INSERT INTO " +
		"user(name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?)"
	TEMPLATE_UPDATE = "UPDATE template " +
		"SET name=?, deployment=?, service=?, endpoints=? " +
		"WHERE id=?"
	TEMPLATE_DELETE = "UPDATE template " +
		"SET status=? " +
		"WHERE id=?"

	QUERY_BY_ID = "SELECT id, name, orgId, deployment, service, endpoint, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"FROM template " +
		"WHERE id=? AND status=1"
	QUERY_ALL_BY_ORGID = "SELECT id, name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM template " +
		"WHERE orgId=? AND status=1"
	QUERY_DUPLICATED_NAME = "SELECT id, name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM template " +
		"WHERE name=? AND orgId=?"
