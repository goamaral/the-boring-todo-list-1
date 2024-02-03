package initializer

import (
	"gorm.io/driver/postgres"

	"example.com/the-boring-to-do-list-1/pkg/env"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

func NewDSN(dsn gorm_provider.DSN) gorm_provider.DSN {
	dsn.Protocol = "postgresql"
	return dsn
}

func DefaultDSN() gorm_provider.DSN {
	return NewDSN(gorm_provider.DSN{
		Host:     env.GetOrDefault("DB_HOST", "localhost"),
		Port:     env.GetOrDefault("DB_PORT", "5432"),
		DBName:   env.GetOrDefault("DB_NAME", "postgres"),
		User:     env.GetOrDefault("DB_USER", "postgres"),
		Password: env.GetOrDefault("DB_PASS", "postgres"),
		SSLMode:  env.GetOrDefault("DB_SSLMODE", "disable"),
	})
}

func NewProvider(dsn gorm_provider.DSN) (gorm_provider.Provider, error) {
	return gorm_provider.NewProvider(postgres.Open(NewDSN(dsn).String()))
}
