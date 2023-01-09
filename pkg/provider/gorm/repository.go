package gormprovider

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db        *gorm.DB
	tableName string
}

type Repository interface {
	TableName() string
	NewQuery(ctx context.Context) *gorm.DB
	NewQueryWithOpts(ctx context.Context, opts ...QueryOption) *gorm.DB
}

func (repo *repository) TableName() string {
	return repo.tableName
}

func (repo *repository) NewQuery(ctx context.Context) *gorm.DB {
	return repo.db.WithContext(ctx).Table(repo.tableName)
}

func (repo *repository) NewQueryWithOpts(ctx context.Context, opts ...QueryOption) *gorm.DB {
	return ApplyQueryOpts(repo.NewQuery(ctx), opts...)
}
