package organization


const (
	ORG_SELECT = "SELECT id, name, cpuQuota, memQuota, budget, balance, status, dcIdList," +
		"createdAt, modifiedAt, modifiedOp, comment " +
		"FROM organization WHERE id=?"
	ORG_SELECT_ALL = "SELECT id, name, cpuQuota, memQuota, budget, balance, status, dcIdList," +
		"createdAt, modifiedAt, modifiedOp, comment " +
		"FROM organization where status=?"

	ORG_SELECT_NAME = "SELECT id, name, cpuQuota, memQuota, budget, balance, status, dcIdList," +
		"createdAt, modifiedAt, modifiedOp, comment " +
		"FROM organization WHERE name=? and status=?"

	ORG_INSERT = "INSERT INTO organization(name, cpuQuota, memQuota, budget, " +
		"balance, status, dcIdList, createdAt, modifiedAt, modifiedOp, comment) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	ORG_UPDATE = "UPDATE organization SET name=?, cpuQuota=?, memQuota=?, budget=?, " +
		"balance=?, status=?, dcIdList=?, modifiedAt=?, modifiedOp=?, comment=? " +
		"WHERE id=?"

	ORG_DELETE = "UPDATE organization SET status=?, modifiedAt=?, modifiedOp=? WHERE id=?"

	VALID   = 1
	INVALID = 0
)
