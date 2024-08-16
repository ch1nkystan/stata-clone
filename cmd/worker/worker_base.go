package main

import (
	"github.com/prosperofair/stata/pkg/pgsql"
)

type Worker struct {
	pg  *pgsql.PGSQLClient
	cfg *Config
}

func NewWorker(pg *pgsql.PGSQLClient, cfg *Config) *Worker {
	return &Worker{
		pg: pg,

		cfg: cfg,
	}
}
