package main

func NewConfig() Config {
	return Config{
		Migrator:       MigratorConfig{},
		Logger:         LoggerConfig{},
		SourcePostgres: SourcePostgresConfig{},
		TargetPostgres: TargetPostgresConfig{},
	}
}

type Config struct {
	Migrator       MigratorConfig
	Logger         LoggerConfig
	SourcePostgres SourcePostgresConfig
	TargetPostgres TargetPostgresConfig
}

type MigratorConfig struct {
	ChunkSize int  `env:"MIGRATOR_CHUNK_SIZE" envDefault:"3000"`
	DryRun    bool `env:"DRY_RUN" envDefault:"false"`
}

type LoggerConfig struct {
	Level    string `env:"LOG_LEVEL" envDefault:"info"`
	Encoding string `env:"LOG_ENCODING" envDefault:"json"`
}

type SourcePostgresConfig struct {
	Database string `env:"SOURCE_POSTGRES_DATABASE" envDefault:"stata"`
	Host     string `env:"SOURCE_POSTGRES_HOST" envDefault:"127.0.0.1"`
	User     string `env:"SOURCE_POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"SOURCE_POSTGRES_PASSWORD" envDefault:"pass"`
	Port     string `env:"SOURCE_POSTGRES_PORT" envDefault:"5432"`
}

type TargetPostgresConfig struct {
	Database string `env:"TARGET_POSTGRES_DATABASE" envDefault:"stata"`
	Host     string `env:"TARGET_POSTGRES_HOST" envDefault:"127.0.0.1"`
	User     string `env:"TARGET_POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"TARGET_POSTGRES_PASSWORD" envDefault:"pass"`
	Port     string `env:"TARGET_POSTGRES_PORT" envDefault:"5432"`
}
