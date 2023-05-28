package repository

import (
	"example.com/the-boring-to-do-list-1/internal/entity"
	gormprovider "example.com/the-boring-to-do-list-1/pkg/gormprovider"
)

const (
	usersTableName = "users"
)

type UserRepository interface {
	gormprovider.Repository[entity.User]
}

func NewUserRepository(gormProvider *gormprovider.Provider) UserRepository {
	return &userRepository{AbstractRepository: gormprovider.NewAbstractRepository[entity.User](gormProvider, usersTableName, "id")}
}

type userRepository struct {
	gormprovider.AbstractRepository[entity.User]
}
