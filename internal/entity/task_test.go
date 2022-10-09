package entity_test

import (
	"encoding/json"
	"testing"
	"time"

	"example.com/fiber-m3o-validator/internal/entity"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
)

func TestTask_MarshallAndUnmarshall(t *testing.T) {
	now := time.Now()

	type Test struct {
		TestName string
		Task     entity.Task
		Validate func(t *testing.T, test Test, newTask entity.Task)
	}

	tests := []Test{
		{
			TestName: "CompletedAt is nil",
			Task: entity.Task{
				AbstractEntity: entity.AbstractEntity{
					Id:        ulid.Make().String(),
					CreatedAt: now,
				},
				Title: "title",
			},
			Validate: func(t *testing.T, _ Test, newTask entity.Task) {
				assert.Nil(t, newTask.CompletedAt, "completed_at")
			},
		},
		{
			TestName: "CompletedAt is present",
			Task: entity.Task{
				AbstractEntity: entity.AbstractEntity{
					Id:        ulid.Make().String(),
					CreatedAt: now,
				},
				Title:       "title",
				CompletedAt: &now,
			},
			Validate: func(t *testing.T, test Test, newTask entity.Task) {
				assert.Equal(t, test.Task.CompletedAt.Format(time.RFC3339), newTask.CompletedAt.Format(time.RFC3339), "completed_at")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			// Marshal
			taskMap := entity.TaskToMap(test.Task)
			data, err := json.Marshal(taskMap)
			assert.NoError(t, err)

			// Unmarshal
			newTaskMap := map[string]interface{}{}
			err = json.Unmarshal(data, &newTaskMap)
			assert.NoError(t, err)
			newTask, err := entity.TaskFromMap(newTaskMap)
			assert.NoError(t, err)

			// Test
			assert.Equal(t, test.Task.Id, newTask.Id, "id")
			assert.Equal(t, test.Task.CreatedAt.Format(time.RFC3339), newTask.CreatedAt.Format(time.RFC3339), "created_at")
			assert.Equal(t, test.Task.Title, newTask.Title, "title")
			test.Validate(t, test, newTask)
		})
	}
}
