package server_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/internal/test"
	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

func TestTask_CreateTask(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	user := test.AddUser(t, gormProvider, entity.User{})
	s := server.NewServer(jwtProvider, gormProvider)
	accessToken, err := server.GenerateAccessToken(jwtProvider, user.UUID)
	require.NoError(t, err)

	t.Run("Created", func(t *testing.T) {
		title := "title"

		res := server.NewTest(
			t, s, fiber.MethodPost, "/tasks",
			server.CreateTaskRequest{Title: title},
		).
			WithCookie("accessToken", accessToken).
			Send()
		require.Equal(t, fiber.StatusFound, res.StatusCode)

		task, err := repository.NewTaskRepository(gormProvider).First(
			context.Background(),
			gorm_provider.OrderOption("created_at DESC"),
		)
		require.NoError(t, err)
		assert.Equal(t, "/tasks/"+task.UUID.String(), res.Header.Get("Location"))
	})

	t.Run("BadRequest", func(t *testing.T) {
		res := server.NewTest(
			t, s, fiber.MethodPost, "/tasks",
			server.CreateTaskRequest{},
		).
			WithCookie("accessToken", accessToken).
			Send()
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})
}

func TestTask_ListTasks(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	user := test.AddUser(t, gormProvider, entity.User{})
	// task := test.AddTask(t, gormProvider, entity.Task{AuthorID: user.ID})

	userB := test.AddUser(t, gormProvider, entity.User{})
	test.AddTask(t, gormProvider, entity.Task{AuthorID: userB.ID})

	s := server.NewServer(jwtProvider, gormProvider)
	accessToken, err := server.GenerateAccessToken(jwtProvider, user.UUID)
	require.NoError(t, err)

	res := server.NewTest(t, s, fiber.MethodGet, "/tasks", server.ListTasksRequest{}).
		WithCookie("accessToken", accessToken).
		Send()
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	// TODO: Check response
	// require.Len(t, res.Tasks, 1)
	// assert.Equal(t, task.UUID, res.Tasks[0].UUID)
}

func TestTask_GetTask(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	user := test.AddUser(t, gormProvider, entity.User{})
	task := test.AddTask(t, gormProvider, entity.Task{AuthorID: user.ID})

	s := server.NewServer(jwtProvider, gormProvider)
	accessToken, err := server.GenerateAccessToken(jwtProvider, user.UUID)
	require.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		res := server.NewTest(t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", task.UUID), nil).
			WithCookie("accessToken", accessToken).
			Send()
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		res := server.NewTest(t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", ulid.Make().String()), nil).
			WithCookie("accessToken", accessToken).
			Send()
		assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
	})
}

func TestTask_PatchTask(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	user := test.AddUser(t, gormProvider, entity.User{})
	taskUuid := test.AddTask(t, gormProvider, entity.Task{AuthorID: user.ID}).UUID

	s := server.NewServer(jwtProvider, gormProvider)
	accessToken, err := server.GenerateAccessToken(jwtProvider, user.UUID)
	require.NoError(t, err)

	title := "updated title"
	doneAt := time.Date(2023, 9, 24, 12, 11, 0, 0, time.UTC)
	patch := repository.TaskPatch{
		Title:  gorm_provider.NewOptionalField(title),
		DoneAt: gorm_provider.NewOptionalField(&doneAt),
	}

	req := server.PatchTaskRequest{Patch: patch}
	res := server.NewTest(t, s, fiber.MethodPatch, fmt.Sprintf("/tasks/%s", taskUuid), req).
		WithCookie("accessToken", accessToken).
		Send()
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	task, err := repository.NewTaskRepository(gormProvider).First(context.Background(), clause.Eq{Column: "uuid", Value: taskUuid})
	require.NoError(t, err)
	assert.Equal(t, title, task.Title)
	assert.Equal(t, doneAt, *task.DoneAt)
}

func TestTask_DeleteTask(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	user := test.AddUser(t, gormProvider, entity.User{})
	taskUuid := test.AddTask(t, gormProvider, entity.Task{AuthorID: user.ID}).UUID

	s := server.NewServer(jwtProvider, gormProvider)
	accessToken, err := server.GenerateAccessToken(jwtProvider, user.UUID)
	require.NoError(t, err)

	res := server.NewTest(t, s, fiber.MethodDelete, fmt.Sprintf("/tasks/%s", taskUuid), nil).
		WithCookie("accessToken", accessToken).
		Send()
	assert.Equal(t, fiber.StatusOK, res.StatusCode)

	_, found, err := repository.NewTaskRepository(gormProvider).FindOne(context.Background(), clause.Eq{Column: "uuid", Value: taskUuid})
	require.NoError(t, err)
	require.False(t, found)
}
