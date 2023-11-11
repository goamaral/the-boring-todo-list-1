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
	"gorm.io/gorm/clause"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/internal/server"
	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	mock_gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider/mocks"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

func TestTask_CreateTask(t *testing.T) {
	t.Run("Created", func(t *testing.T) {
		title := "title"

		taskRepo := mock_gorm_provider.NewAbstractRepository[entity.Task](t)
		taskRepo.Mock.Test(nil)
		taskRepo.EXPECT().
			Create(mock.Anything, mock.Anything).
			RunAndReturn(func(_ context.Context, task *entity.Task, _ ...any) error {
				task.UUID = ulid.Make().String()
				return nil
			})

		s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
		s.TaskController.TaskRepo = taskRepo

		reqBody := server.CreateTaskRequest{Task: server.NewTask{Title: title}}
		testRequest[server.CreateResponse](t, s, fiber.MethodPost, "/tasks", buildReqBodyReader(t, reqBody)).
			Test(fiber.StatusCreated, func(resBody server.CreateResponse) {
				assert.NotZero(t, resBody.UUID)
			})
	})

	t.Run("BadRequest", func(t *testing.T) {
		s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
		reqBody := server.CreateTaskRequest{Task: server.NewTask{}}
		testRequest[string](t, s, fiber.MethodPost, "/tasks", buildReqBodyReader(t, reqBody)).Test(fiber.StatusBadRequest, nil)
	})
}

func TestTask_ListTasks(t *testing.T) {
	uuid := ulid.Make().String()
	reqBody := server.ListTasksRequest{}

	taskRepo := mock_gorm_provider.NewAbstractRepository[entity.Task](t)
	taskRepo.Mock.Test(nil)
	taskRepo.EXPECT().
		Find(mock.Anything, clause.Gt{Column: "id", Value: uint(0)}, clause.Eq{Column: "done_at", Value: nil}).
		Return([]entity.Task{{EntityWithUUID: gorm_provider.EntityWithUUID{UUID: uuid}}}, nil)

	s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
	s.TaskController.TaskRepo = taskRepo

	testRequest[server.ListTasksResponse](t, s, fiber.MethodGet, "/tasks", buildReqBodyReader(t, reqBody)).
		Test(fiber.StatusOK, func(resBody server.ListTasksResponse) {
			require.Len(t, resBody.Tasks, 1)
			assert.Equal(t, uuid, resBody.Tasks[0].UUID)
		})
}

func TestTask_GetTask(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		task := entity.Task{EntityWithUUID: gorm_provider.EntityWithUUID{UUID: ulid.Make().String()}}

		taskRepo := mock_gorm_provider.NewAbstractRepository[entity.Task](t)
		taskRepo.Mock.Test(nil)
		taskRepo.EXPECT().
			First(mock.Anything, clause.Eq{Column: "uuid", Value: task.UUID}).
			Return(task, true, nil)

		s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
		s.TaskController.TaskRepo = taskRepo

		testRequest[server.GetTaskResponse](t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", task.UUID), nil).
			Test(fiber.StatusOK, func(resBody server.GetTaskResponse) {
				assert.Equal(t, task, resBody.Task)
			})
	})

	t.Run("NotFound", func(t *testing.T) {
		uuid := ulid.Make().String()

		taskRepo := mock_gorm_provider.NewAbstractRepository[entity.Task](t)
		taskRepo.Mock.Test(nil)
		taskRepo.EXPECT().
			First(mock.Anything, clause.Eq{Column: "uuid", Value: uuid}).
			Return(entity.Task{}, false, nil)

		s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
		s.TaskController.TaskRepo = taskRepo

		testRequest[string](t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", uuid), nil).Test(fiber.StatusNotFound, nil)
	})
}

func TestTask_PatchTask(t *testing.T) {
	uuid := ulid.Make().String()
	title := "updated title"
	doneAt := time.Date(2023, 9, 24, 12, 11, 0, 0, time.UTC)
	patch := repository.TaskPatch{
		Title:  gorm_provider.NewOptionalField(title),
		DoneAt: gorm_provider.NewOptionalField(&doneAt),
	}

	reqBody := server.PatchTaskRequest{Patch: patch}
	taskRepo := mock_gorm_provider.NewAbstractRepository[entity.Task](t)
	taskRepo.Mock.Test(nil)
	taskRepo.EXPECT().
		Update(mock.Anything, patch, clause.Eq{Column: "uuid", Value: uuid}).
		Return(nil)

	s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
	s.TaskController.TaskRepo = taskRepo

	testRequest[string](t, s, fiber.MethodPatch, fmt.Sprintf("/tasks/%s", uuid), buildReqBodyReader(t, reqBody)).
		Test(fiber.StatusOK, nil)
}

func TestTask_DeleteTask(t *testing.T) {
	uuid := ulid.Make().String()
	taskRepo := mock_gorm_provider.NewAbstractRepository[entity.Task](t)
	taskRepo.Mock.Test(nil)
	taskRepo.EXPECT().
		Delete(mock.Anything, clause.Eq{Column: "uuid", Value: uuid}).
		Return(nil)

	s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
	s.TaskController.TaskRepo = taskRepo

	testRequest[string](t, s, fiber.MethodDelete, fmt.Sprintf("/tasks/%s", uuid), nil).Test(fiber.StatusOK, nil)
}
