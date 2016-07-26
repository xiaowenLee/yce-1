package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const (
	MAX_POOL_SIZE        = 10
	DATABASE_DRIVER      = "mysql"
	DB_HOST              = "172.21.1.11:32306"
	DB_USER              = "root"
	DB_PASSWORD          = "root"
	DB_NAME              = "yce"
	DB_CONNECTION_SUFFIX = "?parseTime=true"
	DELAY_MILLISECONDS   = 5000
)

type MysqlClient struct {
	db       *sql.DB
	host     string
	user     string
	password string
	database string
	pool     int
}

func NewMysqlClient(host, user, password, database string, pool int) *MysqlClient {
	return &MysqlClient{host: host, password: password, database: database, pool: pool}
}

func (c *MysqlClient) Open() {
	endpoint := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ")/" + DB_NAME + DB_CONNECTION_SUFFIX

	db, err := sql.Open(DATABASE_DRIVER, endpoint)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Set Connection Pool
	db.SetMaxOpenConns(c.pool)
	db.SetMaxIdleConns((int)(c.pool / 2))

	c.db = db
}

func (c *MysqlClient) Close() {
	c.db.Close()
}

func (c *MysqlClient) Conn() *sql.DB {
	return c.db
}

// Ping the connection, keep connection alive
func (c *MysqlClient) Ping() {
	select {
	case <-time.After(time.Millisecond * time.Duration(DELAY_MILLISECONDS)):
		err := c.db.Ping()
		if err != nil {
			log.Fatal(err)
			c.Open()
		}
	}
}
