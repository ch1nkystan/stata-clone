package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/gocraft/dbr/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/prosperofair/pkg/depot"
	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
	"github.com/prosperofair/stata/pkg/server"
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
	conn, err := createPostgresConnection(cfg.Postgres)
	if err != nil {
		log.Fatal("failed to make pg connection", zap.Error(err))
	}

	log.Info("running migrations...")
	if err := runMigrations(conn.DB, "./migrations"); err != nil {
		log.Fatal("failed to run migrations", zap.Error(err))
	}

	log.Info("loading pgsql client...")
	pg := pgsql.NewPGSQLClient(conn)

	log.Info("loading depot client...")
	dc := depot.NewClient(&depot.Config{
		Host:  cfg.Depot.Host,
		Token: cfg.Depot.Token,
	})

	log.Info("loading server...")
	s := server.New(&server.Config{
		BackendTokens:  convertTokens(cfg.Server.BackendTokens),
		FrontendTokens: convertTokens(cfg.Server.FrontendTokens),

		ExposeMetrics: cfg.Server.ExposeMetrics,
	}, &server.Deps{
		PG:    pg,
		Depot: dc,
	})

	if err := s.App.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatal("failed to run server", zap.Error(err))
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

func convertTokens(tokens []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, t := range tokens {
		m[t] = struct{}{}
	}

	return m
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath, // File source URL
		"postgres",               // Database name
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
