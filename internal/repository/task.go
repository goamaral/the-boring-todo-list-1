package repository

import (
	"example.com/the-boring-to-do-list-1/internal/entity"
	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

type AbstractTaskRepository interface {
	gorm_provider.AbstractRepository[entity.Task]
}

const (
	tasksTableName = "tasks"
)

type TaskRepository struct {
	gorm_provider.Repository[entity.Task]
}

func NewTaskRepository(gormProvider gorm_provider.AbstractProvider) TaskRepository {
	return TaskRepository{Repository: gorm_provider.NewRepository[entity.Task](gormProvider, tasksTableName)}
}
