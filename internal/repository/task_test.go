package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/internal/test"
	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

func AddTask(t *testing.T, repo repository.TaskRepository, task *entity.Task) *entity.Task {
	require.NoError(t, repo.Create(context.Background(), task))
	return task
}

func TestTaskRepository_TaskFilter(t *testing.T) {
	repo := repository.NewTaskRepository(test.NewTestProvider(t))
	taskA := AddTask(t, repo, &entity.Task{})
	taskB := AddTask(t, repo, &entity.Task{DoneAt: lo.ToPtr(time.Now())})

	type Test struct {
		TaskFilter  repository.TaskFilter
		ExpectedIds []uint
	}
	runTest := func(test Test) func(t *testing.T) {
		return func(t *testing.T) {
			tasks, err := repo.Find(context.Background(), test.TaskFilter)
			require.NoError(t, err)
			for i, expectedTaskId := range test.ExpectedIds {
				assert.Equal(t, expectedTaskId, tasks[i].ID)
			}
			assert.Equal(t, len(tasks), len(test.ExpectedIds))
		}
	}

	t.Run("Blank", runTest(Test{
		TaskFilter:  repository.TaskFilter{},
		ExpectedIds: []uint{taskA.ID, taskB.ID},
	}))

	t.Run("ID", runTest(Test{
		TaskFilter:  repository.TaskFilter{ID: gorm_provider.NewQueryFieldFilter(taskA.ID)},
		ExpectedIds: []uint{taskA.ID},
	}))

	t.Run("IDGt", runTest(Test{
		TaskFilter:  repository.TaskFilter{IDGt: gorm_provider.NewQueryFieldFilter(taskA.ID)},
		ExpectedIds: []uint{taskB.ID},
	}))

	t.Run("UUID", runTest(Test{
		TaskFilter:  repository.TaskFilter{UUID: gorm_provider.NewQueryFieldFilter(taskA.UUID)},
		ExpectedIds: []uint{taskA.ID},
	}))

	t.Run("Done", runTest(Test{
		TaskFilter:  repository.TaskFilter{Done: gorm_provider.NewQueryFieldFilter(true)},
		ExpectedIds: []uint{taskB.ID},
	}))
}

func TestTaskRepository_TaskPatch(t *testing.T) {
	repo := repository.NewTaskRepository(test.NewTestProvider(t))
	newTitle := "New title"
	newDoneAt := time.Date(2023, 9, 24, 12, 31, 0, 0, time.UTC)

	type Test struct {
		Task      entity.Task
		TaskPatch repository.TaskPatch
		Assert    func(t *testing.T, updatedTask entity.Task)
	}
	runTest := func(test Test) func(t *testing.T) {
		return func(t *testing.T) {
			AddTask(t, repo, &test.Task)
			taskFilter := repository.TaskFilter{ID: gorm_provider.NewQueryFieldFilter(test.Task.ID)}
			require.NoError(t, repo.Update(context.Background(), test.TaskPatch, taskFilter, gorm_provider.DebugOption()))
			updatedTask, err := repo.FindOne(context.Background(), taskFilter)
			require.NoError(t, err)
			assert.Greater(t, updatedTask.UpdatedAt.Unix(), test.Task.UpdatedAt.Unix())
			if test.Assert != nil {
				test.Assert(t, updatedTask)
			}
		}
	}

	t.Run("Blank", runTest(Test{}))

	t.Run("Title", runTest(Test{
		TaskPatch: repository.TaskPatch{Title: gorm_provider.NewQueryFieldFilter(newTitle)},
		Assert: func(t *testing.T, updatedTask entity.Task) {
			assert.Equal(t, newTitle, updatedTask.Title)
		},
	}))

	t.Run("DoneAt", runTest(Test{
		TaskPatch: repository.TaskPatch{DoneAt: gorm_provider.NewQueryFieldFilter(&newDoneAt)},
		Assert: func(t *testing.T, updatedTask entity.Task) {
			assert.Equal(t, newDoneAt.Unix(), updatedTask.DoneAt.Unix())
		},
	}))

	t.Run("DoneAt to nil", runTest(Test{
		Task:      entity.Task{DoneAt: lo.ToPtr(time.Now())},
		TaskPatch: repository.TaskPatch{DoneAt: gorm_provider.NewQueryFieldFilter[*time.Time](nil)},
		Assert: func(t *testing.T, updatedTask entity.Task) {
			assert.Nil(t, updatedTask.DoneAt)
		},
	}))
}
