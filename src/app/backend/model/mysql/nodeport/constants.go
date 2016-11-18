package nodeport

//选择时需要同时满足dcId和port
const (
	NP_SELECT = "SELECT svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE port=? AND dcId=?"

	NP_SELECT_BY_DC_AND_PORT = "SELECT port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE port=? AND dcId=?"

	NP_INSERT = "INSERT INTO " + "nodeport(port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment) " + "VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

	NP_UPDATE = "UPDATE nodeport " + "SET port=?, dcId=?, svcName=?, status=?, modifiedAt=?, modifiedOp=?, comment=? " + "WHERE port=? AND dcId=?"

	NP_DELETE = "UPDATE nodeport " + "SET status=?, modifiedAt=?, modifiedOp=?, comment=? " + "WHERE port=? AND dcId=?"

	NP_INSERT_ON_DUPLICATE_KEY_UPDATE = "INSERT INTO nodeport(port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment) " + "VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE svcName=?, status=?, modifiedAt=?, modifiedOp=? "

	NP_SELECT_VALID = "SELECT port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE status=? AND dcId=?"

	NP_SELECT_NEW = "SELECT port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE port=? AND dcId=? ORDER BY port DESC"

	NP_SELECT_OCCUPIED = "SELECT port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment " + "FROM nodeport WHERE status=?"

	NP_SELECT_FREE_BY_DC_AND_PORT = "SELECT port=?, dcId=?, svcName=?, status=?, createdAt=?, modifiedAt=?, modifiedOp=?, comment=?" + " FROM nodeport WHERE status=? AND port=? AND dcId=?"

	NP_SELECT_FREE_BY_DC = "SELECT port, dcId, svcName, status, createdAt, modifiedAt, modifiedOp, comment" + " FROM nodeport WHERE status=? AND dcId=?"

	VALID      = 1
	INVALID    = 0
	PORT_START = 30000
	PORT_LIMIT = 32767

	OCCUPIED = 0
	FREE = 1
)
