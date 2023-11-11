package test

import (
	"os"
	"testing"

	"example.com/the-boring-to-do-list-1/internal/config"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	postgres_gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider/postgres"
	"github.com/stretchr/testify/require"
)

func NewGormProvider(t *testing.T) gorm_provider.Provider {
	schema, err := os.Open(config.RelativePath("../../db/1_schema.sql"))
	require.NoError(t, err)

	seed, err := os.Open(config.RelativePath("../../db/2_seed.sql"))
	require.NoError(t, err)

	return postgres_gorm_provider.NewTestProvider(t, schema, seed, postgres_gorm_provider.DefaultDSN())
}
