package server_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"example.com/fiber-m3o-validator/internal/entity"
	"example.com/fiber-m3o-validator/internal/server"
	"example.com/fiber-m3o-validator/internal/service"
	"example.com/fiber-m3o-validator/mocks"
)

func TestTask_CreateTask(t *testing.T) {
	id := ulid.Make().String()
	title := "title"

	taskService := mocks.NewTaskService(t)
	taskService.On("CreateTask", mock.Anything).
		Return(func(task entity.Task) entity.Task {
			assert.Equal(t, title, task.Title)

			return entity.Task{AbstractEntity: entity.AbstractEntity{Id: id}}
		}, nil)

	s := server.NewServer(taskService)
	reqBody := server.CreateTaskRequest{Task: server.NewTask{Title: title}}
	resBody := server.CreateResponse{}
	res, err := sendRequest(t, s, fiber.MethodPost, "/tasks", reqBody, &resBody)
	if assert.NoError(t, err, err) && assert.Equal(t, fiber.StatusCreated, res.StatusCode, resBody) {
		assert.NotZero(t, resBody.Id)
	}
}

func TestTask_ListTasks(t *testing.T) {
	id := ulid.Make().String()
	var pageSize uint = 10

	taskService := mocks.NewTaskService(t)
	taskService.On("ListTasks", pageSize, mock.Anything).
		Return(func(_ uint, opts *service.ListTasksOpts) []entity.Task {
			assert.Empty(t, opts.PageId)

			return []entity.Task{{AbstractEntity: entity.AbstractEntity{Id: id}}}
		}, nil)

	s := server.NewServer(taskService)
	reqBody := server.ListTasksRequest{PageSize: pageSize}
	resBody := server.ListTasksResponse{}
	res, err := sendRequest(t, s, fiber.MethodGet, "/tasks", reqBody, &resBody)
	if assert.NoError(t, err, err) && assert.Equal(t, fiber.StatusOK, res.StatusCode, resBody) && assert.Len(t, resBody.Tasks, 1) {
		assert.Equal(t, id, resBody.Tasks[0].Id)
	}
}
