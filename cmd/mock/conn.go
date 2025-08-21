package main

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"time"
)

func createPostgresConnection(cfg PostgresConfig) (*dbr.Connection, error) {
	cs := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require sslrootcert=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLCertPath,
	)

	conn, err := dbr.Open("postgres", cs, nil)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.DB.SetConnMaxLifetime(5 * time.Minute)

	return conn, nil
}
