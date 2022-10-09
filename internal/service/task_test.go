package service_test

import (
	"fmt"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.m3o.com/db"

	"example.com/fiber-m3o-validator/internal/entity"
	"example.com/fiber-m3o-validator/internal/service"
	"example.com/fiber-m3o-validator/mocks"
)

func TestTaskService_CreateTask(t *testing.T) {
	mockM3ODbClient := mocks.NewM3ODb(t)
	mockM3ODbClient.On("Create", mock.Anything).
		Return(func(req *db.CreateRequest) *db.CreateResponse {
			assert.NotZero(t, req.Id)
			assert.Equal(t, entity.TasksTableName, req.Table)
			return &db.CreateResponse{Id: req.Id}
		}, nil)

	s := service.NewTaskService(mockM3ODbClient)
	task, err := s.CreateTask(entity.Task{Title: "Test Task"})
	if assert.NoError(t, err) {
		assert.NotZero(t, task.Id)
		assert.NotZero(t, task.CreatedAt)
	}
}

func TestTaskService_ListTasks(t *testing.T) {
	pageId := ulid.Make().String()
	expectedId := ulid.Make().String()
	var pageSize uint = 10

	mockM3ODbClient := mocks.NewM3ODb(t)
	mockM3ODbClient.On("Read", mock.Anything).
		Return(func(req *db.ReadRequest) *db.ReadResponse {
			assert.Equal(t, int32(pageSize), req.Limit)
			assert.Equal(t, "id", req.OrderBy)
			assert.Contains(t, req.Query, fmt.Sprintf("'%s'", pageId))
			assert.Equal(t, entity.TasksTableName, req.Table)

			return &db.ReadResponse{Records: []map[string]interface{}{
				entity.TaskToMap(entity.Task{AbstractEntity: entity.AbstractEntity{Id: expectedId}}),
			}}
		}, nil)

	s := service.NewTaskService(mockM3ODbClient)
	tasks, err := s.ListTasks(pageId, pageSize)
	if assert.NoError(t, err) && assert.Len(t, tasks, 1) {
		assert.Equal(t, expectedId, tasks[0].Id)
	}
}
