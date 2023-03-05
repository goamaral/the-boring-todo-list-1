package gormprovider

import (
	"fmt"
	"io"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/oklog/ulid/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DefaultPageSize = 10

type provider struct {
	*gorm.DB
}

type Provider interface {
	NewRepository(tableName string) Repository
}

func NewProvider(args ...any) (Provider, error) {
	dsn := NewDSN()
	for _, arg := range args {
		switch arg.(type) {
		case DSN:
			dsn = arg.(DSN)
		default:
			panic(fmt.Sprintf("Argument of type %T not supported", arg))
		}
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &provider{DB: db}, nil
}

func NewTestProvider(t *testing.T, schema io.Reader, seed io.Reader) Provider {
	dsn := NewDSN()
	testDsn := NewDSN()
	testDsn.DbName = fmt.Sprintf("%s_test_%s", dsn.DbName, strings.ToLower(ulid.Make().String()))

	// Create test database
	db, err := connect(dsn)
	err = db.Exec("CREATE DATABASE " + testDsn.DbName).Error
	if err != nil {
		t.Fatal(err)
	}

	// Connect to new database
	testDb, err := connect(testDsn)
	if err != nil {
		t.Fatal(err)
	}

	// Drop database and close connections
	t.Cleanup(func() {
		err = disconnect(testDb)
		if err != nil {
			t.Fatal(err)
		}
		err = db.Exec("DROP DATABASE " + testDsn.DbName).Error
		if err != nil {
			t.Fatal(err)
		}
		err = disconnect(db)
		if err != nil {
			t.Fatal(err)
		}
	})

	// Load schema
	schemaSql, err := io.ReadAll(schema)
	if err != nil {
		t.Fatal(err)
	}
	err = testDb.Exec(string(schemaSql)).Error
	if err != nil {
		t.Fatal(err)
	}

	// Load seeds
	seedSql, err := io.ReadAll(seed)
	if err != nil {
		t.Fatal(err)
	}
	err = testDb.Exec(string(seedSql)).Error
	if err != nil {
		t.Fatal(err)
	}

	// Reconnect so loaded schema is visible
	testDb, err = reconnect(testDb, testDsn)
	if err != nil {
		t.Fatal(err)
	}

	return &provider{DB: testDb}
}

func (p *provider) NewRepository(tableName string) Repository {
	return &repository{db: p.DB, tableName: tableName}
}

func connect(dsn DSN) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
}

func disconnect(db *gorm.DB) error {
	rawDb, err := db.DB()
	if err != nil {
		return err
	}
	return rawDb.Close()
}

func reconnect(db *gorm.DB, dsn DSN) (*gorm.DB, error) {
	err := disconnect(db)
	if err != nil {
		return db, err
	}
	return connect(dsn)
}
