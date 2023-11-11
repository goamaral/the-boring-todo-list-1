package postgres_gorm_provider

import (
	"io"
	"testing"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"

	"example.com/the-boring-to-do-list-1/pkg/env"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

func DefaultDSN() gorm_provider.DSN {
	return gorm_provider.DSN{
		Host:     env.GetOrDefault("DB_HOST", "localhost"),
		Port:     env.GetOrDefault("DB_PORT", "5432"),
		DBName:   env.GetOrDefault("DB_NAME", "postgres"),
		User:     env.GetOrDefault("DB_USER", "postgres"),
		Password: env.GetOrDefault("DB_PASS", "postgres"),
	}
}

func NewProvider(dsn gorm_provider.DSN) (gorm_provider.Provider, error) {
	return gorm_provider.NewProvider(postgres.Open(dsn.String()))
}

func NewTestProvider(t *testing.T, schema io.Reader, seed io.Reader, dsn gorm_provider.DSN) gorm_provider.Provider {
	return gorm_provider.NewTestProvider(t, schema, seed, postgres.Open, dsn)
}
