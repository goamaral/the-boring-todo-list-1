package repository

import (
	"example.com/the-boring-to-do-list-1/internal/entity"
	gormprovider "example.com/the-boring-to-do-list-1/pkg/gormprovider"
)

const (
	tasksTableName = "tasks"
)

type taskRepository struct {
	gormprovider.AbstractRepository[entity.Task]
}

type TaskRepository interface {
	gormprovider.Repository[entity.Task]
}

func NewTaskRepository(gormProvider *gormprovider.Provider) TaskRepository {
	return &taskRepository{AbstractRepository: gormprovider.NewAbstractRepository[entity.Task](gormProvider, tasksTableName, "id")}
}
