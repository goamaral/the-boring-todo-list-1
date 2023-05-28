package gormprovider

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type AbstractRepository[T any] struct {
	provider   *Provider
	tableName  string
	primaryKey string
}

type Repository[T any] interface {
	TableName() string
	NewQuery(ctx context.Context) *gorm.DB
	NewQueryWithOpts(ctx context.Context, opts ...QueryOption) *gorm.DB
	Create(ctx context.Context, entity *T) error
	List(ctx context.Context, opts ...QueryOption) ([]T, error)
	Get(ctx context.Context, opts ...QueryOption) (T, bool, error)
	Update(ctx context.Context, update *T, opts ...QueryOption) error
	Patch(ctx context.Context, patch *T, opts ...QueryOption) error
	Delete(ctx context.Context, opts ...QueryOption) error
}

func NewAbstractRepository[T any](provider *Provider, tableName string, primaryKey string) AbstractRepository[T] {
	return AbstractRepository[T]{provider: provider, tableName: tableName, primaryKey: primaryKey}
}

func (repo *AbstractRepository[T]) TableName() string {
	return repo.tableName
}

func (repo *AbstractRepository[T]) NewQuery(ctx context.Context) *gorm.DB {
	db := repo.provider.db
	txCtx, ok := ctx.(TxContext)
	if ok {
		db = txCtx.txDB
	}
	return db.WithContext(ctx).Table(repo.tableName)
}

func (repo *AbstractRepository[T]) NewQueryWithOpts(ctx context.Context, opts ...QueryOption) *gorm.DB {
	return ApplyQueryOpts(repo.NewQuery(ctx), opts...)
}

func (repo *AbstractRepository[T]) Create(ctx context.Context, entity *T) error {
	return repo.NewQuery(ctx).Create(entity).Error
}

func (repo *AbstractRepository[T]) List(ctx context.Context, opts ...QueryOption) ([]T, error) {
	var entities []T
	err := repo.NewQueryWithOpts(ctx, opts...).Find(&entities).Error
	return entities, err
}

func (repo *AbstractRepository[T]) Get(ctx context.Context, opts ...QueryOption) (T, bool, error) {
	var entity T
	err := repo.NewQueryWithOpts(ctx, opts...).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity, false, nil
	}
	if err != nil {
		return entity, false, err
	}
	return entity, true, nil
}

func (repo *AbstractRepository[T]) Update(ctx context.Context, update *T, opts ...QueryOption) error {
	return repo.NewQueryWithOpts(ctx, opts...).Select("*").Omit("created_at", "updated_at").Updates(update).Error
}

func (repo *AbstractRepository[T]) Patch(ctx context.Context, patch *T, opts ...QueryOption) error {
	return repo.NewQueryWithOpts(ctx, opts...).Omit("created_at", "updated_at").Updates(patch).Error
}

func (repo *AbstractRepository[T]) Delete(ctx context.Context, opts ...QueryOption) error {
	var entity T
	return repo.NewQueryWithOpts(ctx, opts...).Delete(&entity).Error
}

func (repo *AbstractRepository[T]) Count(ctx context.Context, opts ...QueryOption) (int64, error) {
	var count int64
	err := repo.NewQueryWithOpts(ctx, opts...).Count(&count).Error
	return count, err
}
