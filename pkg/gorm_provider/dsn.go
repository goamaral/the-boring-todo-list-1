package gorm_provider

import (
	"fmt"
)

type DSN struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

func (dsn DSN) String() string {
	return fmt.Sprintf("postgresql://%s@%s:%s/%s?password=%s", dsn.User, dsn.Host, dsn.Port, dsn.DBName, dsn.Password)
}

func (dsn DSN) ConnectionString() string {
	return fmt.Sprintf("user=%s host=%s port=%s database=%s password=%s sslmode=disable", dsn.User, dsn.Host, dsn.Port, dsn.DBName, dsn.Password)
}
