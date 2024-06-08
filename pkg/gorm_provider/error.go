package gorm_provider

import (
	"github.com/jackc/pgx/v5/pgconn"
)

type ErrorCode string

var UniqueConstraintViolation = ErrorCode("23505")

func HasErrorCode(err error, code ErrorCode) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}
	return pgErr.Code == string(code)
}
