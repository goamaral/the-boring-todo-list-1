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
	var pageSize uint = 10

	type Test struct {
		TestName               string
		Ops                    *service.ListTasksOpts
		ValidateM3OReadRequest func(t *testing.T, test Test, req *db.ReadRequest)
	}

	tests := []Test{
		{
			TestName: "Without options",
			ValidateM3OReadRequest: func(t *testing.T, test Test, req *db.ReadRequest) {
				assert.Equal(t, entity.TasksTableName, req.Table, "Table")
				assert.Equal(t, int32(pageSize), req.Limit, "Limit")
				assert.Empty(t, req.OrderBy, "OrderBy")
				assert.Empty(t, req.Query, "Query")
			},
		},
		{
			TestName: "With PageId option",
			Ops:      &service.ListTasksOpts{PageId: ulid.Make().String()},
			ValidateM3OReadRequest: func(t *testing.T, test Test, req *db.ReadRequest) {
				assert.Equal(t, entity.TasksTableName, req.Table, "Table")
				assert.Equal(t, int32(pageSize), req.Limit, "Limit")
				assert.Equal(t, "id", req.OrderBy, "OrderBy")
				assert.Contains(t, req.Query, fmt.Sprintf("'%s'", test.Ops.PageId), "Query")
			},
		},
	}

	for _, test := range tests {
		expectedId := ulid.Make().String()

		mockM3ODbClient := mocks.NewM3ODb(t)
		mockM3ODbClient.On("Read", mock.Anything).
			Return(func(req *db.ReadRequest) *db.ReadResponse {
				test.ValidateM3OReadRequest(t, test, req)

				return &db.ReadResponse{Records: []map[string]interface{}{
					entity.TaskToMap(entity.Task{AbstractEntity: entity.AbstractEntity{Id: expectedId}}),
				}}
			}, nil)

		s := service.NewTaskService(mockM3ODbClient)
		tasks, err := s.ListTasks(pageSize, test.Ops)
		if assert.NoError(t, err) && assert.Len(t, tasks, 1) {
			assert.Equal(t, expectedId, tasks[0].Id)
		}
	}
}
