package repository_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
	ulid "github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	gormprovider "example.com/the-boring-to-do-list-1/pkg/provider/gorm"
)

func NewTaskRepository(t *testing.T) repository.TaskRepository {
	_, b, _, _ := runtime.Caller(0)
	folderPath := filepath.Dir(b)
	if err := godotenv.Load(folderPath + "/../../.env"); err != nil {
		t.Error(err)
	}

	schema, err := os.Open(folderPath + "/../../db/1_schema.sql")
	if err != nil {
		t.Error(err)
	}

	// seed, err := os.Open(folderPath + "/../../db/2_seed.sql")
	seed := bytes.NewReader(nil)

	gormProvider := gormprovider.NewTestProvider(t, schema, seed)
	return repository.NewTaskRepository(gormProvider)
}

func AddTask(t *testing.T, repo gormprovider.Repository, task *entity.Task) *entity.Task {
	if err := repo.NewQuery(context.Background()).Create(task).Error; err != nil {
		t.Error(err)
	}
	return task
}

func TestTaskRepository_CreateTask(t *testing.T) {
	repo := NewTaskRepository(t)
	title := "Test Task"

	task := entity.Task{Title: title}
	err := repo.CreateTask(context.Background(), &task)
	require.NoError(t, err)
	assert.NotZero(t, task.Id)
	assert.NotZero(t, task.CreatedAt)
}

func TestTaskRepository_ListTasks(t *testing.T) {
	taskIds := []string{ulid.Make().String(), ulid.Make().String()}

	type Test struct {
		TestName string
		Opts     []gormprovider.QueryOption
	}
	runTest := func(t *testing.T, test Test) []entity.Task {
		repo := NewTaskRepository(t)

		for _, taskId := range taskIds {
			AddTask(t, repo.GetGormRepository(), &entity.Task{AbstractEntity: entity.AbstractEntity{Id: taskId}})
		}

		tasks, err := repo.ListTasks(context.Background(), test.Opts...)
		require.NoError(t, err)
		return tasks
	}

	t.Run("WithoutOptions", func(t *testing.T) {
		tasks := runTest(t, Test{})
		assert.Len(t, tasks, 2)
	})

	t.Run("WithPageId", func(t *testing.T) {
		tasks := runTest(t, Test{
			Opts: []gormprovider.QueryOption{
				gormprovider.PaginationOption{PageId: taskIds[0]},
			},
		})
		require.Len(t, tasks, 1)
		assert.Equal(t, tasks[0].Id, taskIds[1])
	})
}
