package repository_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
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

	seed, err := os.Open(folderPath + "/../../db/2_seed.sql")
	if err != nil {
		t.Error(err)
	}

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

func TestTaskRepository_TaskFilter(t *testing.T) {
	// TODO
}
