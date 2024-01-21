package server_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/internal/test"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

func TestAuth_Login(t *testing.T) {
	jwtProvider := jwt_provider.NewTestProvider(t)
	gormProvider := test.NewGormProvider(t)

	s := server.NewServer(jwtProvider, gormProvider)

	username := "username"
	password := "password"
	user := entity.User{Username: username}
	require.NoError(t, user.SetEncryptedPassword(password))
	user = test.AddUser(t, gormProvider, user)

	t.Run("OK", func(t *testing.T) {
		req := server.LoginRequest{Username: username, Password: password}
		res := server.NewTest[server.LoginResponse](t, s, fiber.MethodPost, "/auth/login", req).
			Send().
			UnmarshalBody()

		assert.NotZero(t, res.AccessToken)
		assert.NotZero(t, res.RefreshToken)
	})

	t.Run("BadRequest/InvalidCredentials/UserNotFound", func(t *testing.T) {
		t.Skip() // TODO
	})

	t.Run("BadRequest/InvalidCredentials/InvalidPassword", func(t *testing.T) {
		t.Skip() // TODO
	})
}
