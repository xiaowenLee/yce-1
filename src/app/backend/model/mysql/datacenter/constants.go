package datacenter

const (
	DC_SELECT_BY_ID = "SELECT id, name, host, port, secret, status, nodePort, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM datacenter WHERE id=? AND status=?"

	DC_SELECT_ALL = "SELECT id, name, host, port, secret, status, nodePort, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM datacenter where status=?"

	DC_SELECT_BY_NAME = "SELECT id, name, host, port, secret, status, nodePort, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM datacenter WHERE name=? AND status=?"

	DC_QUERY_DUPLICATED_NAME = "SELECT id, name, host, port, secret, status, nodePort, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM datacenter WHERE name=? and status=?"

	DC_INSERT = "INSERT INTO " +
		"datacenter(name, host, port, secret, status, nodePort, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	DC_UPDATE = "UPDATE datacenter " +
		"SET name=?, host=?, port=?, secret=?, status=?, nodePort=?, modifiedAt=?, modifiedOp=?, comment=? " +
		"WHERE id=?"

	DC_DELETE = "UPDATE datacenter " +
		"SET status=?, modifiedAt=?, modifiedOp=? " +
		"WHERE id=?"

	VALID   = 1
	INVALID = 0
)
