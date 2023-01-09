package repository

import (
	"context"

	ulid "github.com/oklog/ulid/v2"

	"example.com/the-boring-to-do-list-1/internal/entity"
	gormprovider "example.com/the-boring-to-do-list-1/pkg/provider/gorm"
)

const (
	tasksTableName = "tasks"
)

type taskRepository struct {
	gormprovider.Repository
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *entity.Task) error
	ListTasks(ctx context.Context, opts ...gormprovider.QueryOption) ([]entity.Task, error)
}

func NewTaskRepository(gormProvider gormprovider.Provider) *taskRepository {
	return &taskRepository{Repository: gormProvider.NewRepository("tasks")}
}

/* PUBLIC */
func (repo *taskRepository) CreateTask(ctx context.Context, task *entity.Task) error {
	task.Id = ulid.Make().String()
	return repo.Repository.NewQuery(ctx).Create(&task).Error
}

func (repo *taskRepository) ListTasks(ctx context.Context, opts ...gormprovider.QueryOption) ([]entity.Task, error) {
	var tasks []entity.Task

	err := repo.NewQueryWithOpts(ctx, opts...).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
