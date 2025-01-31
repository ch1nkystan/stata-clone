package main

import (
	"fmt"
	"time"

	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
	"go.uber.org/zap"
)

type Worker struct {
	pg  *pgsql.Client
	cfg WorkerConfig
}

func NewWorker(pg *pgsql.Client, cfg WorkerConfig) *Worker {
	return &Worker{
		pg:  pg,
		cfg: cfg,
	}
}

func (w *Worker) run() {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		if err := w.pg.CreateOnlineHistory(now.Add(-w.cfg.Interval), now); err != nil {
			log.Error(fmt.Sprintf("failed to create online history: %v", err), zap.Error(err))
		}

	}
}
