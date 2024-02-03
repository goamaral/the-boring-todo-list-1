package gorm_provider

import "fmt"

type DSN struct {
	Protocol string
	User     string
	Host     string
	Port     string
	DBName   string
	Password string
	SSLMode  string
}

func (dsn DSN) String() string {
	return fmt.Sprintf("%s://%s@%s:%s/%s?password=%s&sslmode=%s&timezone=UTC", dsn.Protocol, dsn.User, dsn.Host, dsn.Port, dsn.DBName, dsn.Password, dsn.SSLMode)
}
