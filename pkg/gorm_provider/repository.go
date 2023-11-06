package gorm_provider

import (
	"context"

	"gorm.io/gorm/clause"
)

type AbstractRepository[T AbstractEntity] interface {
	NewTransaction(ctx context.Context, fc func(context.Context) error) error
	NewQuery(ctx context.Context, clauses ...clause.Expression) Query[T]
	Create(ctx context.Context, record *T, clauses ...clause.Expression) error
	Find(ctx context.Context, clauses ...clause.Expression) ([]T, error)
	FindInBatches(ctx context.Context, bacthSize int, fn func([]T) error, clauses ...clause.Expression) error
	FindOne(ctx context.Context, clauses ...clause.Expression) (T, error)
	First(ctx context.Context, clauses ...clause.Expression) (T, bool, error)
	Update(ctx context.Context, update any, clauses ...clause.Expression) error
	Delete(ctx context.Context, clauses ...clause.Expression) error
}

type Repository[T AbstractEntity] struct {
	Provider  AbstractProvider
	TableName string
}

func NewRepository[T AbstractEntity](provider AbstractProvider, tableName string) Repository[T] {
	return Repository[T]{provider, tableName}
}

func (am Repository[T]) NewTransaction(ctx context.Context, fc func(context.Context) error) error {
	return am.Provider.NewTransaction(ctx, fc)
}

func (am Repository[T]) NewQuery(ctx context.Context, clauses ...clause.Expression) Query[T] {
	return Query[T]{DB: GetDbFromContextOr(ctx, am.Provider.GetDb()).Table(am.TableName).Clauses(clauses...)}
}

func (am Repository[T]) Create(ctx context.Context, record *T, clauses ...clause.Expression) error {
	return am.NewQuery(ctx, clauses...).Create(record)
}

func (am Repository[T]) Find(ctx context.Context, clauses ...clause.Expression) ([]T, error) {
	return am.NewQuery(ctx, clauses...).Find()
}

func (am Repository[T]) FindInBatches(ctx context.Context, bacthSize int, fn func([]T) error, clauses ...clause.Expression) error {
	return am.NewQuery(ctx, clauses...).FindInBatches(bacthSize, fn)
}

func (am Repository[T]) FindOne(ctx context.Context, clauses ...clause.Expression) (T, error) {
	return am.NewQuery(ctx, clauses...).FindOne()
}

func (am Repository[T]) First(ctx context.Context, clauses ...clause.Expression) (T, bool, error) {
	return am.NewQuery(ctx, clauses...).First()
}

func (am Repository[T]) Update(ctx context.Context, update any, clauses ...clause.Expression) error {
	return am.NewQuery(ctx, clauses...).Update(update)
}

func (am Repository[T]) Delete(ctx context.Context, clauses ...clause.Expression) error {
	return am.NewQuery(ctx, clauses...).Delete()
}
