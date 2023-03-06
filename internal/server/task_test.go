package server_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/mocks"
	gormprovider "example.com/the-boring-to-do-list-1/pkg/provider/gorm"
)

func TestTask_CreateTask(t *testing.T) {
	id := ulid.Make().String()
	title := "title"

	taskRepo := mocks.NewTaskRepository(t)
	taskRepo.On("CreateTask", mock.Anything, mock.Anything).
		Return(func(_ context.Context, task *entity.Task) error {
			assert.Equal(t, title, task.Title)
			task.AbstractEntity = entity.AbstractEntity{Id: id}
			return nil
		}, nil)

	s := server.NewServer(taskRepo)
	reqBody := server.CreateTaskRequest{Task: server.NewTask{Title: title}}
	resBody := server.CreateResponse{}
	res, err := sendRequest(t, s, fiber.MethodPost, "/tasks", buildReqBodyReader(t, reqBody), &resBody)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, res.StatusCode, resBody)
	assert.NotZero(t, resBody.Id)
}

func TestTask_ListTasks(t *testing.T) {
	id := ulid.Make().String()
	reqBody := server.ListTasksRequest{}

	taskRepo := mocks.NewTaskRepository(t)
	taskRepo.On("ListTasks", mock.Anything, gormprovider.PaginationOption{PageId: reqBody.PageId, PageSize: reqBody.PageSize}).
		Return([]entity.Task{{AbstractEntity: entity.AbstractEntity{Id: id}}}, nil)

	s := server.NewServer(taskRepo)
	resBody := server.ListTasksResponse{}
	res, err := sendRequest(t, s, fiber.MethodGet, "/tasks", buildReqBodyReader(t, reqBody), &resBody)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, res.StatusCode, resBody)
	require.Len(t, resBody.Tasks, 1)
	assert.Equal(t, id, resBody.Tasks[0].Id)
}

func TestTask_GetTask(t *testing.T) {
	task := entity.Task{AbstractEntity: entity.AbstractEntity{Id: ulid.Make().String()}}

	taskRepo := mocks.NewTaskRepository(t)
	taskRepo.On("GetTask", mock.Anything, repository.TaskFilter{Id: task.Id}).
		Return(task, nil)

	s := server.NewServer(taskRepo)
	resBody := server.GetTaskResponse{}
	res, err := sendRequest(t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", task.Id), nil, &resBody)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, res.StatusCode, resBody)
	assert.Equal(t, task, resBody.Task)
}
