package server_test

import (
	"os"
	"testing"

	"example.com/the-boring-to-do-list-1/internal/server"
	"example.com/the-boring-to-do-list-1/internal/test"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	test.LoadEnv()
	os.Exit(m.Run())
}

func TestServer_HealthCheck(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := server.NewServer(jwt_provider.NewTestProvider(t), nil)
		res := server.NewTest(t, s, fiber.MethodGet, "/health", nil).Send()
		assert.Equal(t, fiber.StatusNoContent, res.StatusCode)
	})
}
