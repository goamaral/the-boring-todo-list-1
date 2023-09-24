package repository

import (
	"example.com/the-boring-to-do-list-1/internal/entity"
	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

type AbstractUserRepository interface {
	gorm_provider.AbstractRepository[entity.User]
}

const (
	usersTableName = "users"
)

type UserRepository struct {
	gorm_provider.Repository[entity.User]
}

func NewUserRepository(gormProvider gorm_provider.AbstractProvider) UserRepository {
	return UserRepository{Repository: gorm_provider.NewRepository[entity.User](gormProvider, usersTableName)}
}
