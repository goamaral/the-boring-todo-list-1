package server_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"example.com/fiber-m3o-validator/internal/entity"
	"example.com/fiber-m3o-validator/internal/server"
	"example.com/fiber-m3o-validator/mocks"
)

func TestTask_CreateTask(t *testing.T) {
	t.Run("Created", func(t *testing.T) {
		expectedId := ulid.Make().String()

		taskService := mocks.NewTaskService(t)
		taskService.On("CreateTask", mock.Anything).
			Return(entity.Task{AbstractEntity: entity.AbstractEntity{Id: expectedId}}, nil)

		s := server.NewServer(taskService)
		reqBody := server.CreateTaskRequest{Task: server.NewTask{Title: "test title"}}
		resBody := server.CreateResponse{}
		res := sendRequest(t, s, fiber.MethodPost, "/tasks", reqBody, &resBody)
		if assert.Equal(t, fiber.StatusCreated, res.StatusCode, resBody) {
			assert.NotZero(t, resBody.Id)
		}
	})
}

func TestTask_ListTasks(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		expectedTaskId := ulid.Make().String()
		var pageSize uint = 10

		taskService := mocks.NewTaskService(t)
		taskService.On("ListTasks", "", pageSize).
			Return([]entity.Task{{AbstractEntity: entity.AbstractEntity{Id: expectedTaskId}}}, nil)

		s := server.NewServer(taskService)
		reqBody := server.ListTasksRequest{PageSize: pageSize}
		resBody := server.ListTasksResponse{}
		res := sendRequest(t, s, fiber.MethodGet, "/tasks", reqBody, &resBody)
		if assert.Equal(t, fiber.StatusOK, res.StatusCode, resBody) && assert.Len(t, resBody.Tasks, 1) {
			assert.Equal(t, expectedTaskId, resBody.Tasks[0].Id)
		}
	})
}
