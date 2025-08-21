package main

func NewConfig() Config {
	return Config{
		Postgres: PostgresConfig{},
	}
}

type Config struct {
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Database    string `env:"POSTGRES_DATABASE" envDefault:"stata"`
	Host        string `env:"POSTGRES_HOST" envDefault:"db-postgresql-drapps-do-user-13508037-0.c.db.ondigitalocean.com"`
	User        string `env:"POSTGRES_USER" envDefault:"doadmin"`
	Password    string `env:"POSTGRES_PASSWORD" envDefault:"AVNS_Sh9FTNbSIV4Ma506inv"`
	Port        string `env:"POSTGRES_PORT" envDefault:"25060"`
	SSLCertPath string `env:"POSTGRES_SSL_CERT_PATH" envDefault:"ca-certificate.crt"`
}
