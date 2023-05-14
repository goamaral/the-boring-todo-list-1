package server_test

import (
	"context"
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
	t.Run("Created", func(t *testing.T) {
		title := "title"

		taskRepo := NewTaskRepository(t)
		taskRepo.On("Create", mock.Anything, mock.Anything).Return(func(_ context.Context, task *entity.Task) error {
			task.Id = ulid.Make().String()
			return nil
		})

		s := server.NewServer(taskRepo)
		reqBody := server.CreateTaskRequest{Task: server.NewTask{Title: title}}
		testRequest[server.CreateResponse](t, s, fiber.MethodPost, "/tasks", buildReqBodyReader(t, reqBody)).
			Test(fiber.StatusCreated, func(resBody server.CreateResponse) {
				assert.NotZero(t, resBody.Id)
			})
	})

	t.Run("BadRequest", func(t *testing.T) {
		s := server.NewServer(nil)
		reqBody := server.CreateTaskRequest{Task: server.NewTask{}}
		testRequest[string](t, s, fiber.MethodPost, "/tasks", buildReqBodyReader(t, reqBody)).Test(fiber.StatusBadRequest, nil)
	})
}

func TestTask_ListTasks(t *testing.T) {
	id := ulid.Make().String()
	reqBody := server.ListTasksRequest{}

	taskRepo := NewTaskRepository(t)
	taskRepo.On("List", mock.Anything, gormprovider.PaginationOption{PageId: reqBody.PageId, PageSize: reqBody.PageSize}, repository.TaskFilter{}).
		Return([]entity.Task{{AbstractEntity: gormprovider.AbstractEntity{Id: id}}}, nil)

	s := server.NewServer(taskRepo)
	testRequest[server.ListTasksResponse](t, s, fiber.MethodGet, "/tasks", buildReqBodyReader(t, reqBody)).
		Test(fiber.StatusOK, func(resBody server.ListTasksResponse) {
			require.Len(t, resBody.Tasks, 1)
			assert.Equal(t, id, resBody.Tasks[0].Id)
		})
}

func TestTask_GetTask(t *testing.T) {
	task := entity.Task{AbstractEntity: gormprovider.AbstractEntity{Id: ulid.Make().String()}}

	taskRepo := NewTaskRepository(t)
	taskRepo.On("Get", mock.Anything, repository.TaskFilter{Id: &task.Id}).Return(task, nil)

	s := server.NewServer(taskRepo)
	testRequest[server.GetTaskResponse](t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", task.Id), nil).
		Test(fiber.StatusOK, func(resBody server.GetTaskResponse) {
			assert.Equal(t, task, resBody.Task)
		})
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
		repository.TaskFilter{Id: &id},
	).Return(nil)

	s := server.NewServer(taskRepo)
	testRequest[string](t, s, fiber.MethodPut, fmt.Sprintf("/tasks/%s", id), buildReqBodyReader(t, reqBody)).
		Test(fiber.StatusOK, nil)
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
		repository.TaskFilter{Id: &id},
	).Return(nil)

	s := server.NewServer(taskRepo)
	testRequest[string](t, s, fiber.MethodPatch, fmt.Sprintf("/tasks/%s", id), buildReqBodyReader(t, reqBody)).
		Test(fiber.StatusOK, nil)
}

func TestTask_DeleteTask(t *testing.T) {
	id := ulid.Make().String()
	taskRepo := NewTaskRepository(t)
	taskRepo.On("Delete", mock.Anything, repository.TaskFilter{Id: &id}).Return(nil)

	s := server.NewServer(taskRepo)
	testRequest[string](t, s, fiber.MethodDelete, fmt.Sprintf("/tasks/%s", id), nil).Test(fiber.StatusOK, nil)
}
