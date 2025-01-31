package main

import "time"

const (
	WorkerFBToolFetcher  = "fbtool-fetcher"
	WorkerOnlineSnapshot = "online-snapshot"
)

func NewConfig() Config {
	return Config{
		Worker:   WorkerConfig{},
		Logger:   LoggerConfig{},
		Depot:    DepotConfig{},
		Postgres: PostgresConfig{},
	}
}

type Config struct {
	Worker   WorkerConfig
	Logger   LoggerConfig
	Depot    DepotConfig
	Postgres PostgresConfig
}

type WorkerConfig struct {
	Name string `env:"WORKER_NAME" envDefault:""`

	SingleRun bool          `env:"WORKER_SINGLE_RUN" envDefault:"true"`
	DryRun    bool          `env:"WORKER_DRY_RUN" envDefault:"false"`
	Timeout   time.Duration `env:"WORKER_RUN_TIMEOUT" envDefault:"30s"`
}

type LoggerConfig struct {
	Level    string `env:"LOG_LEVEL" envDefault:"info"`
	Encoding string `env:"LOG_ENCODING" envDefault:"json"`
}

type DepotConfig struct {
	Host  string `env:"DEPOT_HOST" envDefault:"https://depot.drapps.lol"`
	Token string `env:"DEPOT_TOKEN" envDefault:"fcbae6821cb21dd6e2aba928e281e7da"`
}

type PostgresConfig struct {
	Database string `env:"POSTGRES_DATABASE" envDefault:"stata"`
	Host     string `env:"POSTGRES_HOST" envDefault:"127.0.0.1"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"pass"`
	Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
}
