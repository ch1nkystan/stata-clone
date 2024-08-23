package main

import (
	"github.com/prosperofair/stata/pkg/pgsql"
)

type Worker struct {
	pg  *pgsql.Client
	cfg *Config
}

func NewWorker(pg *pgsql.Client, cfg *Config) *Worker {
	return &Worker{pg: pg, cfg: cfg}
}
