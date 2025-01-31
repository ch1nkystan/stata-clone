package main

import (
	"github.com/prosperofair/pkg/log"
	"go.uber.org/zap"
)

func (w *Worker) usersOnlineSnapshot() error {

	if err := w.pg.CreateUsersOnlineSnapshot(); err != nil {
		log.Error("failed to create users online snapshot", zap.Error(err))
	}

	return nil
}
