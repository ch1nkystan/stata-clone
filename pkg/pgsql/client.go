package pgsql

import (
	"errors"

	"github.com/gocraft/dbr/v2"
)

var ErrAlreadyExists = errors.New("err_already_exists")

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
