package dcquota

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
