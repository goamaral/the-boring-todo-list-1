package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"

	"example.com/the-boring-to-do-list-1/internal/config"
	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/internal/test"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

func TestMain(m *testing.M) {
	config.LoadTestEnv()
	os.Exit(m.Run())
}

func TestTaskRepository_TaskPatch(t *testing.T) {
	gormProvider := test.NewGormProvider(t)
	user := test.AddUser(t, gormProvider, entity.User{})
	repo := repository.NewTaskRepository(gormProvider)
	newTitle := "New title"
	newDoneAt := time.Date(2023, 9, 24, 12, 31, 0, 0, time.UTC)

	type Test struct {
		Task      entity.Task
		TaskPatch repository.TaskPatch
		Assert    func(t *testing.T, updatedTask entity.Task)
	}
	runTest := func(test Test) func(t *testing.T) {
		return func(t *testing.T) {
			test.Task.AuthorID = user.ID
			require.NoError(t, repo.Create(context.Background(), &test.Task))
			filter := clause.Eq{Column: "id", Value: test.Task.ID}
			require.NoError(t, repo.Update(context.Background(), test.TaskPatch, filter))
			updatedTask, err := repo.First(context.Background(), filter)
			require.NoError(t, err)
			assert.Greater(t, updatedTask.UpdatedAt.UnixMicro(), test.Task.UpdatedAt.UnixMicro())
			if test.Assert != nil {
				test.Assert(t, updatedTask)
			}
		}
	}

	t.Run("Blank", runTest(Test{}))

	t.Run("Title", runTest(Test{
		TaskPatch: repository.TaskPatch{Title: gorm_provider.NewOptionalField(newTitle)},
		Assert: func(t *testing.T, updatedTask entity.Task) {
			assert.Equal(t, newTitle, updatedTask.Title)
		},
	}))

	t.Run("DoneAt", runTest(Test{
		TaskPatch: repository.TaskPatch{DoneAt: gorm_provider.NewOptionalField(&newDoneAt)},
		Assert: func(t *testing.T, updatedTask entity.Task) {
			assert.Equal(t, newDoneAt.Unix(), updatedTask.DoneAt.Unix())
		},
	}))

	t.Run("DoneAt to nil", runTest(Test{
		Task:      entity.Task{DoneAt: lo.ToPtr(time.Now())},
		TaskPatch: repository.TaskPatch{DoneAt: gorm_provider.NewOptionalField[*time.Time](nil)},
		Assert: func(t *testing.T, updatedTask entity.Task) {
			assert.Nil(t, updatedTask.DoneAt)
		},
	}))
}
