package test

import (
	"context"
	"testing"
	"time"

	"example.com/the-boring-to-do-list-1/internal/initializer"
	"example.com/the-boring-to-do-list-1/pkg/fs"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	testcontainers_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
)

func NewTestDB(t *testing.T, ctx context.Context) gorm_provider.DSN {
	dsn := initializer.DefaultDSN()
	container, err := testcontainers_postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:17"),
		testcontainers_postgres.WithDatabase(dsn.DBName),
		testcontainers_postgres.WithUsername(dsn.User),
		testcontainers_postgres.WithPassword(dsn.Password),
		testcontainers_postgres.WithInitScripts(fs.ResolveRelativePath("../../db/1_schema.sql"), fs.ResolveRelativePath("../../db/2_seed.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	port, err := container.MappedPort(ctx, nat.Port("5432/tcp"))
	require.NoError(t, err)

	dsn.Port = port.Port()

	return dsn
}

func NewGormProvider(t *testing.T, ctx context.Context) gorm_provider.Provider {
	gormProvider, err := gorm_provider.NewProvider(postgres.Open(NewTestDB(t, ctx).String()))
	require.NoError(t, err)
	return gormProvider
}
