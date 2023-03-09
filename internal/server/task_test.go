package server_test

import (
	"fmt"
	"testing"
	"time"

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

func NewTaskRepository(t *testing.T) *mocks.TaskRepository {
	taskRepo := mocks.NewTaskRepository(t)
	taskRepo.Mock.Test(nil)
	return taskRepo
}

func TestTask_CreateTask(t *testing.T) {
	title := "title"

	taskRepo := NewTaskRepository(t)
	taskRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

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

	taskRepo := NewTaskRepository(t)
	taskRepo.On("List", mock.Anything, gormprovider.PaginationOption{PageId: reqBody.PageId, PageSize: reqBody.PageSize}).
		Return([]entity.Task{{AbstractEntity: gormprovider.AbstractEntity{Id: id}}}, nil)

	s := server.NewServer(taskRepo)
	resBody := server.ListTasksResponse{}
	res, err := sendRequest(t, s, fiber.MethodGet, "/tasks", buildReqBodyReader(t, reqBody), &resBody)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, res.StatusCode, resBody)
	require.Len(t, resBody.Tasks, 1)
	assert.Equal(t, id, resBody.Tasks[0].Id)
}

func TestTask_GetTask(t *testing.T) {
	task := entity.Task{AbstractEntity: gormprovider.AbstractEntity{Id: ulid.Make().String()}}

	taskRepo := NewTaskRepository(t)
	taskRepo.On("Get", mock.Anything, repository.TaskFilter{Id: task.Id}).Return(task, nil)

	s := server.NewServer(taskRepo)
	resBody := server.GetTaskResponse{}
	res, err := sendRequest(t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", task.Id), nil, &resBody)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, res.StatusCode, resBody)
	assert.Equal(t, task, resBody.Task)
}

func TestTask_UpdateTask(t *testing.T) {
	id := ulid.Make().String()
	title := "updated title"
	completedAt := time.Now().UTC()
	reqBody := server.UpdateTaskRequest{
		Task: entity.Task{
			Title:       title,
			CompletedAt: &completedAt,
		},
	}

	taskRepo := NewTaskRepository(t)
	taskRepo.On("Update", mock.Anything,
		&entity.Task{
			AbstractEntity: gormprovider.AbstractEntity{Id: id},
			Title:          title,
			CompletedAt:    &completedAt,
		},
		repository.TaskFilter{Id: id},
	).Return(nil)

	s := server.NewServer(taskRepo)
	var resBody string
	res, err := sendRequest(t, s, fiber.MethodPut, fmt.Sprintf("/tasks/%s", id), buildReqBodyReader(t, reqBody), &resBody)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, res.StatusCode, resBody)
}

func TestTask_PatchTask(t *testing.T) {
	id := ulid.Make().String()
	title := "updated title"
	reqBody := server.UpdateTaskRequest{
		Task: entity.Task{Title: title},
	}
	taskRepo := NewTaskRepository(t)
	taskRepo.On("Patch", mock.Anything,
		&entity.Task{AbstractEntity: gormprovider.AbstractEntity{Id: id}, Title: title},
		repository.TaskFilter{Id: id},
	).Return(nil)

	s := server.NewServer(taskRepo)
	var resBody string
	res, err := sendRequest(t, s, fiber.MethodPatch, fmt.Sprintf("/tasks/%s", id), buildReqBodyReader(t, reqBody), &resBody)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, res.StatusCode, resBody)
}
