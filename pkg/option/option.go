package option

import (
	"database/sql"
)

type Option[T any] sql.Null[T]

func Some[T any](v T) Option[T] {
	return Option[T]{Valid: true, V: v}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) ValidOr(v T) T {
	if o.Valid {
		return o.V
	}
	return v
}
