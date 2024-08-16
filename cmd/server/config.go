package main

func NewConfig() Config {
	return Config{
		Server:   ServerConfig{},
		Logger:   LoggerConfig{},
		Postgres: PostgresConfig{},
		Depot:    DepotConfig{},
	}
}

type Config struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Postgres PostgresConfig
	Depot    DepotConfig
}

type ServerConfig struct {
	ExposeMetrics bool   `env:"SERVER_EXPOSE_METRICS" envDefault:"true"`
	Port          string `env:"SERVER_PORT" envDefault:"8080"`

	BackendTokens  []string `env:"SERVER_BACKEND_TOKENS" envDefault:""`
	FrontendTokens []string `env:"SERVER_FRONTEND_TOKENS" envDefault:""`
}

type DepotConfig struct {
	Host  string `env:"DEPOT_HOST" envDefault:"https://depot.drapps.lol"`
	Token string `env:"DEPOT_TOKEN" envDefault:"fcbae6821cb21dd6e2aba928e281e7da"`
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
