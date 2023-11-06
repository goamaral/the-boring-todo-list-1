package test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	gorm_provider_postgres "example.com/the-boring-to-do-list-1/pkg/gorm_provider/postgres"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func RelativePath(relativePath string) string {
	_, file, _, _ := runtime.Caller(1)
	folderPath := filepath.Dir(file)
	return folderPath + "/" + relativePath
}

func LoadEnv() error {
	return godotenv.Load(RelativePath("../../secrets/.env.test"))
}

func LoadEnvT(t require.TestingT) {
	require.NoError(t, LoadEnv())
}

func NewTestProvider(t *testing.T) gorm_provider.Provider {
	LoadEnvT(t)

	schema, err := os.Open(RelativePath("../../db/1_schema.sql"))
	require.NoError(t, err)

	seed, err := os.Open(RelativePath("../../db/2_seed.sql"))
	require.NoError(t, err)

	return gorm_provider_postgres.NewTestProvider(t, schema, seed, gorm_provider_postgres.NewDSN())
}
