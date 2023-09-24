package gorm_provider

import (
	"context"
	"database/sql/driver"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QueryFieldFilter[T any] struct {
	Val     T
	Defined bool
}

func NewQueryFieldFilter[T any](val T) QueryFieldFilter[T] {
	return QueryFieldFilter[T]{Defined: true, Val: val}
}

func (o QueryFieldFilter[T]) Value() (driver.Value, error) {
	return driver.DefaultParameterConverter.ConvertValue(o.Val)
}

func (o QueryFieldFilter[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{SQL: "?", Vars: []any{o.Val}}
}

type QuerySliceFieldFilter[T any] struct {
	Defined bool
	Val     []T
}

func NewQuerySliceFieldFilter[T any](val []T) QuerySliceFieldFilter[T] {
	return QuerySliceFieldFilter[T]{Defined: true, Val: val}
}

func (o QuerySliceFieldFilter[T]) Value() (driver.Value, error) {
	return driver.DefaultParameterConverter.ConvertValue(o.Val)
}

func (o QuerySliceFieldFilter[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{SQL: "?", Vars: []any{o.Val}}
}

func (o *QuerySliceFieldFilter[T]) Append(items ...T) {
	o.Defined = true
	o.Val = append(o.Val, items...)
}

func (o *QuerySliceFieldFilter[T]) Concat(items []T) {
	o.Defined = true
	o.Val = append(o.Val, items...)
}
