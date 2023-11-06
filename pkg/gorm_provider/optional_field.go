package gorm_provider

import (
	"context"
	"database/sql/driver"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AbstractOptionalField interface {
	Defined() bool
	Value() (driver.Value, error)
	GormValue(ctx context.Context, db *gorm.DB) clause.Expr
}

type OptionalField[T any] struct {
	Val       T
	IsDefined bool
}

func NewOptionalField[T any](val T) OptionalField[T] {
	return OptionalField[T]{IsDefined: true, Val: val}
}

func (o OptionalField[T]) Defined() bool {
	return o.IsDefined
}

func (o OptionalField[T]) Value() (driver.Value, error) {
	return driver.DefaultParameterConverter.ConvertValue(o.Val)
}

func (o OptionalField[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{SQL: "?", Vars: []any{o.Val}}
}

type OptionalSliceField[T any] struct {
	IsDefined bool
	Val       []T
}

func NewOptionalSliceField[T any](val []T) OptionalSliceField[T] {
	return OptionalSliceField[T]{IsDefined: true, Val: val}
}

func (o OptionalSliceField[T]) Defined() bool {
	return o.IsDefined
}

func (o OptionalSliceField[T]) Value() (driver.Value, error) {
	return driver.DefaultParameterConverter.ConvertValue(o.Val)
}

func (o OptionalSliceField[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{SQL: "?", Vars: []any{o.Val}}
}

func (o *OptionalSliceField[T]) Append(items ...T) {
	o.IsDefined = true
	o.Val = append(o.Val, items...)
}

func (o *OptionalSliceField[T]) Concat(items []T) {
	o.IsDefined = true
	o.Val = append(o.Val, items...)
}
