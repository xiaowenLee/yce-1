package quota


const (
	QUOTA_SELECT = "SELECT id, name, cpu, mem, rbd, price, " +
		"status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM quota WHERE id=?"

	QUOTA_SELECT_ALL = "SELECT id, name, cpu, mem, rbd, price, " +
		"status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM quota"

	QUOTA_SELECT_ALL_ORDER_BY_CPU = "SELECT id, name, cpu, mem, rbd, price, " +
		"status, createdAt, modifiedAt, modifiedOp, comment " +
		"FROM quota " +
		"ORDER BY cpu ASC"

	QUOTA_INSERT = "INSERT INTO " +
		"quota(name, cpu, mem, rbd, price, status, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	QUOTA_UPDATE = "UPDATE quota " +
		"SET name=?, cpu=?, mem=?, rbd=?, price=?, status=?, modifiedAt=?, modifiedOp=?, comment=? " +
		"WHERE id=?"

	QUOTA_DELETE = "UPDATE quota " +
		"SET status=?, modifiedAt=?, modifiedOp=? " +
		"WHERE id=?"

	VALID   = 1
	INVALID = 0
)
