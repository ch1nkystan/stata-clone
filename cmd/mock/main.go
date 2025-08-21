package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
	"go.uber.org/zap"
	"time"
)

func main() {
	minUsersPerDay := flag.Int("minUsersPerDay", 10, "minimum number of users per day")
	maxUsersPerDay := flag.Int("maxUsersPerDay", 50, "maximum number of users per day")
	deepLinkCount := flag.Int("deepLinksCount", 5, "deep links count")
	referralUserPercentage := flag.Float64("referralUsersPercentage", 0.15, "referral users percentage")
	referralUserCountMin := flag.Int("referralUsersCountMin", 1, "referral users count min")
	referralUserCountMax := flag.Int("referralUsersCountMax", 5, "referral users count max")
	botID := flag.Int("botID", 140, "bot id")
	periodDays := flag.Int("periodDays", 7, "period in days")

	flag.Parse()

	logger := log.GetLogger()

	cfg := NewConfig()
	if err := env.Parse(&cfg); err != nil {
		logger.Fatal("failed to parse .env to config")
	}

	logger.Info("creating pgdb connection...")
	conn, err := createPostgresConnection(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	logger.Info("running migrations...")
	if err = pgsql.RunMigrations(conn.DB, "./migrations"); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	pg := pgsql.NewClient(conn)

	generator, err := newGenerator(
		pg,
		*minUsersPerDay,
		*maxUsersPerDay,
		*deepLinkCount,
		*referralUserPercentage,
		*referralUserCountMin,
		*referralUserCountMax,
		time.Duration(*periodDays)*24*time.Hour,
		*botID,
		logger,
	)
	if err != nil {
		logger.Fatal("failed to create mock data", zap.Error(err))
	}

	logger.Info("starting creating mock data")
	if err = generator.Generate(); err != nil {
		panic(err)
	}
	logger.Info("mock data generated")
}
