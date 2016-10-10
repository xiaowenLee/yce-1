package deployment


const (
	DEPLOYMENT_SELECT = "SELECT id, name, actionType, actionVerb, actionUrl, " +
		"actionAt, actionOp, dcList, success, reason, json, comment, orgId " +
		"FROM deployment where id=?"

	DEPLOYMENT_BYNAME = "SELECT id, name, actionType, actionVerb, actionUrl, " +
		"actionAt, actionOp, dcList, success, reason, json, comment, orgId " +
		"FROM deployment where name=? ORDER BY id DESC LIMIT 30"

	DEPLOYMENT_INSERT = "INSERT INTO deployment(name, actionType, actionVerb, actionUrl, " +
		"actionAt, actionOp, dcList, success, reason, json, comment, orgId) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	DEPLOYMENT_ACTIONTYPE_STAT = "SELECT id, name, actionType, actionVerb, actionUrl, " +
		"actionAt, actionOp, dcList, success, reason, json, comment, orgId " +
		"FROM deployment where actionType=?"

	VALID   = 1
	INVALID = 0
)

var OPERATION_LOG = `
SELECT date_format(d.actionAt, '%Y-%m-%d') AS date, sum(d.actionType) AS total, 2 AS op
FROM deployment d
WHERE date_format(d.actionAt, '%Y-%m-%d') IN
  (
    SELECT DISTINCT date_format(d.actionAt, '%Y-%m-%d')
    FROM deployment d
    WHERE date_sub(curdate(), INTERVAL 7 DAY) <= date_format(d.actionAt, '%Y-%m-%d')
  )
  AND d.actionType = 2
GROUP BY date_format(d.actionAt, '%Y-%m-%d')

UNION

SELECT date_format(d.actionAt, '%Y-%m-%d') AS date, sum(d.actionType) AS total, 3 AS op
FROM deployment d
WHERE date_format(d.actionAt, '%Y-%m-%d') IN
  (
    SELECT DISTINCT date_format(d.actionAt, '%Y-%m-%d')
    FROM deployment d
    WHERE date_sub(curdate(), INTERVAL 7 DAY) <= date_format(d.actionAt, '%Y-%m-%d')
  )
  AND d.actionType = 3
GROUP BY date_format(d.actionAt, '%Y-%m-%d')

UNION

SELECT date_format(d.actionAt, '%Y-%m-%d') AS date, sum(d.actionType) AS total, 4 AS op
FROM deployment d
WHERE date_format(d.actionAt, '%Y-%m-%d') IN
  (
    SELECT DISTINCT date_format(d.actionAt, '%Y-%m-%d')
    FROM deployment d
    WHERE date_sub(curdate(), INTERVAL 7 DAY) <= date_format(d.actionAt, '%Y-%m-%d')
  )
  AND d.actionType = 4
GROUP BY date_format(d.actionAt, '%Y-%m-%d')

UNION

SELECT date_format(d.actionAt, '%Y-%m-%d') AS date, sum(d.actionType) AS total, 8 AS op
FROM deployment d
WHERE date_format(d.actionAt, '%Y-%m-%d') IN
  (
    SELECT DISTINCT date_format(d.actionAt, '%Y-%m-%d')
    FROM deployment d
    WHERE date_sub(curdate(), INTERVAL 7 DAY) <= date_format(d.actionAt, '%Y-%m-%d')
  )
  AND d.actionType = 8
GROUP BY date_format(d.actionAt, '%Y-%m-%d')

UNION

SELECT date_format(d.actionAt, '%Y-%m-%d') AS date, sum(d.actionType) AS total, 9 AS op
FROM deployment d
WHERE date_format(d.actionAt, '%Y-%m-%d') IN
  (
    SELECT DISTINCT date_format(d.actionAt, '%Y-%m-%d')
    FROM deployment d
    WHERE date_sub(curdate(), INTERVAL 7 DAY) <= date_format(d.actionAt, '%Y-%m-%d')
  )
  AND d.actionType = 9
GROUP BY date_format(d.actionAt, '%Y-%m-%d')
`
