package server_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	m3o "go.m3o.com"
	"go.m3o.com/db"

	"example.com/fiber-m3o-validator/mocks"
	"example.com/fiber-m3o-validator/server"
)

func newServer(t *testing.T, setupMocks func(*mocks.M3ODb) error) server.Server {
	m3oDbClient := mocks.NewM3ODb(t)

	err := setupMocks(m3oDbClient)
	if err != nil {
		t.Fatalf("failed to setup mocks: %v", err)
	}

	return server.NewServer(&m3o.Client{Db: m3oDbClient})
}

func Test_Tasks_Post(t *testing.T) {
	t.Run("Created", func(t *testing.T) {
		server := newServer(t, func(mockM3ODbClient *mocks.M3ODb) error {
			mockM3ODbClient.On("Create", mock.Anything).Return(&db.CreateResponse{Id: "task-id"}, nil)
			return nil
		})

		reqBody := fiber.Map{"task": fiber.Map{"title": "test title"}}
		resBody, res := sendRequest(t, server, fiber.MethodPost, "/tasks", reqBody)
		if assert.Equal(t, fiber.StatusCreated, res.StatusCode, resBody) {
			assert.NotZero(t, resBody["id"])
		}
	})
}
