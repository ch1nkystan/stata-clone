package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/gocraft/dbr/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
)

func main() {
	log.Info("setting default timezone to UTC...")
	if err := os.Setenv("TZ", "UTC"); err != nil {
		log.Fatal("failed to set UTC timezone", zap.Error(err))
	}

	log.Info("loading .env file...")
	if err := godotenv.Load(); err != nil {
		log.Error("failed to load .env file", zap.Error(err))
	}

	log.Info("loading config...")
	cfg := NewConfig()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("failed to parse .env to config")
	}

	log.SetLogEncoding(cfg.Logger.Encoding)
	log.SetLogLevel(cfg.Logger.Level)

	log.Info("creating pgdb connection...")
	conn, err := createTargetPostgresConnection(cfg.TargetPostgres)
	if err != nil {
		log.Fatal("failed to make pg connection", zap.Error(err))
	}

	log.Info("loading pgsql client...")
	new := pgsql.NewClient(conn)

	conn2, err := createSourcePostgresConnection(cfg.SourcePostgres)
	if err != nil {
		log.Fatal("failed to make pg connection", zap.Error(err))
	}

	old := NewPGSQLClient(conn2)

	log.Info("loading migrator...")
	m := NewMigrator(new, old, cfg.Migrator.ChunkSize, cfg.Migrator.DryRun)
	if err := m.Run(); err != nil {
		log.Fatal("failed to run migrator", zap.Error(err))
	}
}

func createTargetPostgresConnection(cfg TargetPostgresConfig) (*dbr.Connection, error) {
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

func createSourcePostgresConnection(cfg SourcePostgresConfig) (*dbr.Connection, error) {
	cs := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=ams",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
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
