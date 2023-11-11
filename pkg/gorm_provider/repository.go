package gorm_provider

import (
	"context"
)

type AbstractRepository[T AbstractEntity] interface {
	NewTransaction(ctx context.Context, fc func(context.Context) error) error
	NewQuery(ctx context.Context, opts ...any) Query[T]
	Create(ctx context.Context, record *T, opts ...any) error
	Find(ctx context.Context, opts ...any) ([]T, error)
	FindInBatches(ctx context.Context, bacthSize int, fn func([]T) error, opts ...any) error
	FindOne(ctx context.Context, opts ...any) (T, error)
	First(ctx context.Context, opts ...any) (T, bool, error)
	Update(ctx context.Context, update any, opts ...any) error
	Delete(ctx context.Context, opts ...any) error
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

func (am Repository[T]) NewQuery(ctx context.Context, opts ...any) Query[T] {
	db := GetDbFromContextOr(ctx, am.Provider.GetDb()).WithContext(ctx).Table(am.TableName)
	return Query[T]{
		DB: ApplyQueryOptions(db, opts...),
	}
}

func (am Repository[T]) Create(ctx context.Context, record *T, opts ...any) error {
	return am.NewQuery(ctx, opts...).Create(record)
}

func (am Repository[T]) Find(ctx context.Context, opts ...any) ([]T, error) {
	return am.NewQuery(ctx, opts...).Find()
}

func (am Repository[T]) FindInBatches(ctx context.Context, bacthSize int, fn func([]T) error, opts ...any) error {
	return am.NewQuery(ctx, opts...).FindInBatches(bacthSize, fn)
}

func (am Repository[T]) FindOne(ctx context.Context, opts ...any) (T, error) {
	return am.NewQuery(ctx, opts...).FindOne()
}

func (am Repository[T]) First(ctx context.Context, opts ...any) (T, bool, error) {
	return am.NewQuery(ctx, opts...).First()
}

func (am Repository[T]) Update(ctx context.Context, update any, opts ...any) error {
	return am.NewQuery(ctx, opts...).Update(update)
}

func (am Repository[T]) Delete(ctx context.Context, opts ...any) error {
	return am.NewQuery(ctx, opts...).Delete()
}
