package main

import "time"

func NewConfig() Config {
	return Config{
		Logger:   LoggerConfig{},
		Postgres: PostgresConfig{},
		Worker:   WorkerConfig{},
	}
}

type Config struct {
	Logger   LoggerConfig
	Postgres PostgresConfig
	Worker   WorkerConfig
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

type WorkerConfig struct {
	Interval time.Duration `env:"ONLINE_SNAPSHOT_INTERVAL" envDefault:"5m"`
}
