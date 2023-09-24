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
)

func AddTask(t *testing.T, repo repository.TaskRepository, task *entity.Task) *entity.Task {
	require.NoError(t, repo.Create(context.Background(), task))
	return task
}

func TestTaskRepository_TaskFilter(t *testing.T) {
	repo := repository.NewTaskRepository(test.NewTestProvider(t))
	taskA := AddTask(t, repo, &entity.Task{})
	taskB := AddTask(t, repo, &entity.Task{CompletedAt: lo.ToPtr(time.Now())})

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
			assert.GreaterOrEqual(t, len(tasks), len(test.ExpectedIds))
		}
	}

	t.Run("Blank", runTest(Test{
		TaskFilter:  repository.TaskFilter{},
		ExpectedIds: []uint{taskA.ID, taskB.ID},
	}))

	t.Run("UUID", runTest(Test{
		TaskFilter:  repository.TaskFilter{UUID: &taskA.UUID},
		ExpectedIds: []uint{taskA.ID},
	}))

	t.Run("IDGt", runTest(Test{
		TaskFilter:  repository.TaskFilter{IDGt: &taskA.ID},
		ExpectedIds: []uint{taskB.ID},
	}))

	t.Run("IsComplete", runTest(Test{
		TaskFilter:  repository.TaskFilter{IsComplete: lo.ToPtr(true)},
		ExpectedIds: []uint{taskB.ID},
	}))
}
