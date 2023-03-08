package repository_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"

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

func AddTask(t *testing.T, repo gormprovider.Repository[entity.Task], task *entity.Task) *entity.Task {
	if err := repo.Create(context.Background(), task); err != nil {
		t.Error(err)
	}
	return task
}

func TestTaskRepository_TaskFilter(t *testing.T) {
	// TODO
}
