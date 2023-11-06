package server_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/server"
	mock_repository "example.com/the-boring-to-do-list-1/mocks/repository"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

func TestAuth_Login(t *testing.T) {
	username := "username"
	password := "password"

	t.Run("OK", func(t *testing.T) {
		user := entity.User{Username: username}
		require.NoError(t, user.SetEncryptedPassword(password))

		userRepo := mock_repository.NewAbstractUserRepository(t)
		userRepo.Mock.Test(nil)
		userRepo.EXPECT().
			First(
				mock.Anything,
				gorm_provider.SelectClause("id", "username"),
				clause.Eq{Column: "username", Value: username},
			).
			Return(user, true, nil)

		s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
		s.AuthController.UserRepo = userRepo

		reqBody := server.LoginRequest{Username: username, Password: password}
		testRequest[server.LoginResponse](t, s, fiber.MethodPost, "/auth/login", buildReqBodyReader(t, reqBody)).
			Test(fiber.StatusOK, func(resBody server.LoginResponse) {
				assert.NotZero(t, resBody.AccessToken)
				assert.NotZero(t, resBody.RefreshToken)
			})
	})

	t.Run("BadRequest/InvalidCredentials/UserNotFound", func(t *testing.T) {
		t.Skip() // TODO
	})

	t.Run("BadRequest/InvalidCredentials/InvalidPassword", func(t *testing.T) {
		t.Skip() // TODO
	})
}
