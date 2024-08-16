package main

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gocraft/dbr/v2"
	"github.com/joho/godotenv"
	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
	"go.uber.org/zap"
)

func main() {
	log.Info("setting default timezone to UTC...")
	if err := os.Setenv("TZ", "UTC"); err != nil {
		log.Fatal("failed to set UTC timezone", zap.Error(err))
	}

	if err := godotenv.Load(); err != nil {
		log.Error("failed to load .env file", zap.Error(err))
	}

	cfg := NewConfig()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("failed to parse .env to config")
	}

	log.SetLogEncoding(cfg.Logger.Encoding)
	log.SetLogLevel(cfg.Logger.Level)

	log.Info("creating pgdb connection...")
	conn, err := createPostgresConnection(cfg.Postgres)
	if err != nil {
		log.Fatal("failed to make pg connection", zap.Error(err))
	}

	log.Info("loading pgsql client...")
	pg := pgsql.NewPGSQLClient(conn)

	worker := NewWorker(pg, &cfg)

	switch cfg.Worker.Name {
	case WorkerFBToolFetcher:
		if err := runWorker(worker.fbtoolFetcher, cfg.Worker.SingleRun); err != nil {
			log.Fatal("failed to run worker", zap.Error(err))
		}
	}
}

func runWorker(f func() error, singleRun bool) error {
	timeout := 30 * time.Second

	for {
		log.Info("starting iteration...")
		if err := f(); err != nil {
			return err
		}

		log.Info("finished iteration, sleeping...", zap.Duration("timeout", timeout))

		if singleRun {
			log.Info("single run mode, exiting...")
			return nil
		}

		time.Sleep(timeout)
	}
}

func createPostgresConnection(cfg PostgresConfig) (*dbr.Connection, error) {
	cs := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require sslrootcert=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		"ca-certificate.crt",
	)

	conn, err := dbr.Open("postgres", cs, nil)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
