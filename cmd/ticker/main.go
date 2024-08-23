package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gocraft/dbr/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
	"go.uber.org/zap"
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

	log.Info("loading pgsql client...")
	pg := pgsql.NewClient(conn)

	// select all active assets

	assets := []string{"ETH", "BNB", "BTC"}

	for _, asset := range assets {
		log.Info("creating ticker...", zap.String("asset", asset))
		name := fmt.Sprintf("%sUSDT", asset)
		ticker := NewTicker(pg, name, cfg.Ticker.Tick, cfg.Ticker.DryRun)

		go ticker.run()
		time.Sleep(1 * time.Second)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Info("waiting for SIGTERM signal...")

	sig := <-sigs
	log.Info("signal received", zap.String("signal", sig.String()))
	log.Info("worker finished...")
}

func createPostgresConnection(cfg PostgresConfig) (*dbr.Connection, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	cs := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
		cfg.SSLCertPath,
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
