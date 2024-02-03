package server_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	accessToken, err := jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
	})
	require.NoError(t, err)

	t.Run("Created", func(t *testing.T) {
		title := "title"

		res := server.NewTest[server.CreateResponse](t, s, fiber.MethodPost, "/tasks", server.CreateTaskRequest{Task: server.NewTask{Title: title}}).
			WithAuthorizationHeader(accessToken).
			Send().
			ExpectsStatusCode(fiber.StatusCreated).
			UnmarshalBody()

		task, err := repository.NewTaskRepository(gormProvider).First(context.Background(), clause.Eq{Column: "uuid", Value: res.UUID})
		require.NoError(t, err)
		assert.Equal(t, task.AuthorID, user.ID)
	})

	t.Run("BadRequest", func(t *testing.T) {
		server.NewTest[any](t, s, fiber.MethodPost, "/tasks", server.CreateTaskRequest{Task: server.NewTask{}}).
			WithAuthorizationHeader(accessToken).
			Send().
			ExpectsStatusCode(fiber.StatusBadRequest)
	})
}

func TestTask_ListTasks(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	user := test.AddUser(t, gormProvider, entity.User{})
	task := test.AddTask(t, gormProvider, entity.Task{AuthorID: user.ID})

	userB := test.AddUser(t, gormProvider, entity.User{})
	test.AddTask(t, gormProvider, entity.Task{AuthorID: userB.ID})

	s := server.NewServer(jwtProvider, gormProvider)
	accessToken, err := jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
	})
	require.NoError(t, err)

	res := server.NewTest[server.ListTasksResponse](t, s, fiber.MethodGet, "/tasks", server.ListTasksRequest{}).
		WithAuthorizationHeader(accessToken).
		Send().
		UnmarshalBody()

	require.Len(t, res.Tasks, 1)
	assert.Equal(t, task.UUID, res.Tasks[0].UUID)
}

func TestTask_GetTask(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	user := test.AddUser(t, gormProvider, entity.User{})
	task := test.AddTask(t, gormProvider, entity.Task{AuthorID: user.ID})

	s := server.NewServer(jwtProvider, gormProvider)
	accessToken, err := jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
	})
	require.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		res := server.NewTest[server.GetTaskResponse](t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", task.UUID), nil).
			WithAuthorizationHeader(accessToken).
			Send().
			UnmarshalBody()

		assert.Equal(t, task.UUID, res.Task.UUID)
	})

	t.Run("NotFound", func(t *testing.T) {
		server.NewTest[server.GetTaskResponse](t, s, fiber.MethodGet, fmt.Sprintf("/tasks/%s", ulid.Make().String()), nil).
			WithAuthorizationHeader(accessToken).
			Send().
			ExpectsStatusCode(fiber.StatusNotFound)
	})
}

func TestTask_PatchTask(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	user := test.AddUser(t, gormProvider, entity.User{})
	taskUuid := test.AddTask(t, gormProvider, entity.Task{AuthorID: user.ID}).UUID

	s := server.NewServer(jwtProvider, gormProvider)
	accessToken, err := jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
	})
	require.NoError(t, err)

	title := "updated title"
	doneAt := time.Date(2023, 9, 24, 12, 11, 0, 0, time.UTC)
	patch := repository.TaskPatch{
		Title:  gorm_provider.NewOptionalField(title),
		DoneAt: gorm_provider.NewOptionalField(&doneAt),
	}

	req := server.PatchTaskRequest{Patch: patch}
	server.NewTest[any](t, s, fiber.MethodPatch, fmt.Sprintf("/tasks/%s", taskUuid), req).
		WithAuthorizationHeader(accessToken).
		Send().
		ExpectsStatusCode(fiber.StatusOK)

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
	accessToken, err := jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
	})
	require.NoError(t, err)

	server.NewTest[any](t, s, fiber.MethodDelete, fmt.Sprintf("/tasks/%s", taskUuid), nil).
		WithAuthorizationHeader(accessToken).
		Send().
		ExpectsStatusCode(fiber.StatusOK)

	_, found, err := repository.NewTaskRepository(gormProvider).FindOne(context.Background(), clause.Eq{Column: "uuid", Value: taskUuid})
	require.NoError(t, err)
	require.False(t, found)
}
