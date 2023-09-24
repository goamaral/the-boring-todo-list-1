package gorm_provider

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type AbstractProvider interface {
	GetDBFromContext(ctx context.Context) *gorm.DB
	NewTransaction(ctx context.Context, fc func(context.Context) error) error
}

type Provider struct {
	DB *gorm.DB
}

func NewProvider(dialector gorm.Dialector) (Provider, error) {
	db, err := gorm.Open(dialector)
	if err != nil {
		return Provider{}, err
	}
	return Provider{DB: db}, nil
}

func (p Provider) GetDBFromContext(ctx context.Context) *gorm.DB {
	ctxWithTx, ok := ctx.(TxContext)
	if ok && ctxWithTx.Tx != nil {
		return ctxWithTx.Tx
	}
	return p.DB
}

type TxContext struct {
	context.Context
	Tx *gorm.DB
}

func (p Provider) NewTransaction(ctx context.Context, fc func(context.Context) error) error {
	tx := p.GetDBFromContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	return tx.Transaction(func(tx *gorm.DB) error {
		return fc(TxContext{Context: ctx, Tx: tx})
	})
}

func NewTestProvider(t *testing.T, schema io.Reader, seed io.Reader, dialectorOpenFn func(string) gorm.Dialector, dsn DSN) Provider {
	dialector := dialectorOpenFn(dsn.String())
	testDsn := dsn
	testDsn.DBName = fmt.Sprintf("%s_test_%s", dsn.DBName, strconv.FormatUint(rand.Uint64(), 16))
	testDialector := dialectorOpenFn(testDsn.String())

	// Create test database
	db := connect(t, dialector)
	require.NoError(t, db.Exec("CREATE DATABASE "+testDsn.DBName).Error)

	// Connect to new database
	testDb := connect(t, testDialector)

	// Drop database and close connections
	t.Cleanup(func() {
		disconnect(t, testDb)
		require.NoError(t, db.Exec("DROP DATABASE "+testDsn.DBName).Error)
		disconnect(t, db)
	})

	// Load schema
	schemaSql, err := io.ReadAll(schema)
	require.NoError(t, err)
	require.NoError(t, testDb.Exec(string(schemaSql)).Error)

	// Load seeds
	seedSql, err := io.ReadAll(seed)
	require.NoError(t, err)
	require.NoError(t, testDb.Exec(string(seedSql)).Error)

	// Reconnect so loaded schema is visible
	testDb = reconnect(t, testDialector, testDb)

	return Provider{DB: testDb}
}

func connect(t *testing.T, dialector gorm.Dialector) *gorm.DB {
	db, err := gorm.Open(dialector)
	require.NoError(t, err)
	return db
}

func disconnect(t *testing.T, db *gorm.DB) {
	rawDb, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, rawDb.Close())
}

func reconnect(t *testing.T, dialector gorm.Dialector, db *gorm.DB) *gorm.DB {
	disconnect(t, db)
	return connect(t, dialector)
}
