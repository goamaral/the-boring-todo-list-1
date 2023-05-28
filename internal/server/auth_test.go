package server_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/mocks"
	"example.com/the-boring-to-do-list-1/pkg/gormprovider"
	"example.com/the-boring-to-do-list-1/pkg/jwtprovider"
)

func NewUserReposiotry(t *testing.T) *mocks.UserRepository {
	userRepo := mocks.NewUserRepository(t)
	userRepo.Mock.Test(nil)
	return userRepo
}

func TestAuth_Login(t *testing.T) {
	username := "username"
	password := "password"

	t.Run("OK", func(t *testing.T) {
		user := entity.User{Username: username}
		require.NoError(t, user.SetEncryptedPassword(password))

		jwtProvider := jwtprovider.NewTestProvider(t)

		userRepo := NewUserReposiotry(t)
		userRepo.On("Get", mock.Anything, repository.UserFilter{Username: &username}, gormprovider.SelectOption("id", "username")).
			Return(user, true, nil)

		s := server.NewServer(jwtProvider, nil, userRepo)
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
