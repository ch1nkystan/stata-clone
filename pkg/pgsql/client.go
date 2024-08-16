package pgsql

import (
	"github.com/gocraft/dbr/v2"
)

type PGSQLClient struct {
	pgdb *dbr.Connection
}

func NewPGSQLClient(conn *dbr.Connection) *PGSQLClient {
	return &PGSQLClient{
		pgdb: conn,
	}
}

func (c *PGSQLClient) GetSession() *dbr.Session {
	return c.pgdb.NewSession(nil)
}
