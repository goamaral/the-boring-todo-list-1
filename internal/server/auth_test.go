package server_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/internal/test"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

func TestAuth_Login(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	s := server.NewServer(jwtProvider, gormProvider)

	username := "username"
	password := "password"
	user := entity.User{Username: username}
	require.NoError(t, user.SetEncryptedPassword(password))
	user = test.AddUser(t, gormProvider, user)

	t.Run("OK", func(t *testing.T) {
		req := server.LoginRequest{Username: username, Password: password}
		res := server.NewTest(t, s, fiber.MethodPost, "/auth/login", req).Send()
		assert.Equal(t, fiber.StatusFound, res.StatusCode)
		assert.Equal(t, "/tasks", res.Header.Get("Location"))
		assert.NotZero(t, lo.CountBy(res.Cookies(), func(c *http.Cookie) bool { return c.Name == "accessToken" }))
		assert.NotZero(t, lo.CountBy(res.Cookies(), func(c *http.Cookie) bool { return c.Name == "refreshToken" }))
	})

	t.Run("BadRequest/InvalidCredentials/UserNotFound", func(t *testing.T) {
		t.Skip() // TODO
	})

	t.Run("BadRequest/InvalidCredentials/InvalidPassword", func(t *testing.T) {
		t.Skip() // TODO
	})
}

func TestAuth_Register(t *testing.T) {
	ctx := context.Background()
	gormProvider := test.NewGormProvider(t, ctx)
	jwtProvider := jwt_provider.NewTestProvider(t)

	s := server.NewServer(jwtProvider, gormProvider)

	password := "password"

	t.Run("OK", func(t *testing.T) {
		req := server.RegisterRequest{Username: test.RandomString(), Password: password, ConfirmPassword: password}
		res := server.NewTest(t, s, fiber.MethodPost, "/auth/register", req).Send()
		assert.Equal(t, fiber.StatusFound, res.StatusCode)
		assert.Equal(t, "/tasks", res.Header.Get("Location"))
		assert.NotZero(t, lo.CountBy(res.Cookies(), func(c *http.Cookie) bool { return c.Name == "accessToken" }))
		assert.NotZero(t, lo.CountBy(res.Cookies(), func(c *http.Cookie) bool { return c.Name == "refreshToken" }))
	})

	t.Run("BadRequest/Validation/InvalidConfirmPassword", func(t *testing.T) {
		req := server.RegisterRequest{Username: test.RandomString(), Password: password, ConfirmPassword: password + "a"}
		res := server.NewTest(t, s, fiber.MethodPost, "/auth/register", req).Send()
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		// TODO: assert.Equal(t, "password", res["confirmPassword"].(map[string]any)["eqfield"])
	})

	t.Run("BadRequest/UserAlreadyExists", func(t *testing.T) {
		user := entity.User{}
		require.NoError(t, user.SetEncryptedPassword(password))
		user = test.AddUser(t, gormProvider, user)

		req := server.RegisterRequest{Username: user.Username, Password: password, ConfirmPassword: password}
		res := server.NewTest(t, s, fiber.MethodPost, "/auth/register", req).Send()
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		// TODO: assert.Equal(t, server.ErrUserAlreadyExists.Error(), res["error"])
	})
}
