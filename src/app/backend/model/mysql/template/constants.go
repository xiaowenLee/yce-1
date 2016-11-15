package template

const (
	VALID = 1
	INVALID = 0

	TEMPLATE_INSERT = "INSERT INTO " +
		"template(name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	TEMPLATE_INSERT_ON_DUPLICATE_UPDATE = "INSERT INTO " +
		"template(name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE status=?, deployment=?, service=?, endpoints=?"
	TEMPLATE_UPDATE = "UPDATE template " +
		"SET name=?, deployment=?, service=?, endpoints=? " +
		"WHERE id=?"
	TEMPLATE_DELETE = "UPDATE template " +
		"SET status=?, modifiedAt=?, modifiedOp=? " +
		"WHERE id=?"

	QUERY_BY_ID = "SELECT id, name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM template " +
		"WHERE id=? AND status=?"
	QUERY_BY_NAME = "SELECT id, name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM template " +
		"WHERE name=? AND status=?"
	QUERY_ALL_BY_ORGID = "SELECT id, name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM template " +
		"WHERE orgId=? AND status=?"
	QUERY_DUPLICATED_NAME = "SELECT id, name, orgId, deployment, service, endpoints, status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM template " +
		"WHERE name=? AND orgId=?"
)
