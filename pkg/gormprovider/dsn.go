package gormprovider

import (
	"fmt"

	"example.com/the-boring-to-do-list-1/pkg/env"
)

type DSN struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
}

func NewDSN() DSN {
	return DSN{
		Host:     env.GetOrDefault("DB_HOST", "localhost"),
		Port:     env.GetOrDefault("DB_PORT", "5432"),
		DbName:   env.GetOrDefault("DB_NAME", "postgres"),
		User:     env.GetOrDefault("DB_USER", "postgres"),
		Password: env.GetOrDefault("DB_PASSWORD", "postgres"),
	}
}

func (dsn DSN) SetHost(host string) DSN {
	dsn.Host = host
	return dsn
}

func (dsn DSN) SetPort(port string) DSN {
	dsn.Port = port
	return dsn
}

func (dsn DSN) SetDbName(dbName string) DSN {
	dsn.DbName = dbName
	return dsn
}

func (dsn DSN) SetUser(user string) DSN {
	dsn.User = user
	return dsn
}

func (dsn DSN) SetPassword(password string) DSN {
	dsn.Password = password
	return dsn
}

func (dsn DSN) String() string {
	return fmt.Sprintf("postgresql://%s@%s:%s/%s?password=%s", dsn.User, dsn.Host, dsn.Port, dsn.DbName, dsn.Password)
}

func (dsn DSN) ConnectionString() string {
	return fmt.Sprintf("user=%s host=%s port=%s database=%s password=%s sslmode=disable", dsn.User, dsn.Host, dsn.Port, dsn.DbName, dsn.Password)
}
