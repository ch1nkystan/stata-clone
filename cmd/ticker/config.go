package main

import "time"

const ()

func NewConfig() Config {
	return Config{
		Logger:   LoggerConfig{},
		Postgres: PostgresConfig{},
		Ticker:   TickerConfig{},
	}
}

type Config struct {
	Logger   LoggerConfig
	Postgres PostgresConfig
	Ticker   TickerConfig
}

type LoggerConfig struct {
	Level    string `env:"LOG_LEVEL" envDefault:"info"`
	Encoding string `env:"LOG_ENCODING" envDefault:"json"`
}

type PostgresConfig struct {
	Database    string `env:"POSTGRES_DATABASE" envDefault:"stata"`
	Host        string `env:"POSTGRES_HOST" envDefault:"127.0.0.1"`
	User        string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password    string `env:"POSTGRES_PASSWORD" envDefault:"pass"`
	Port        string `env:"POSTGRES_PORT" envDefault:"5432"`
	SSLMode     string `env:"POSTGRES_SSL_MODE" envDefault:"require"`
	SSLCertPath string `env:"POSTGRES_SSL_CERT_PATH" envDefault:"ca-certificate.crt"`
}

type TickerConfig struct {
	Tick    time.Duration `env:"TICKER_TIME_DURATION" envDefault:"300s"`
	Symbols []string      `env:"TICKER_SYMBOL" envDefault:"BTCUSDT,ETHUSDT,BNBUSDT"`
	DryRun  bool          `env:"TICKER_DRY_RUN" envDefault:"false"`
}
