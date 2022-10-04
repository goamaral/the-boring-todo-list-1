package server_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.m3o.com/db"

	"example.com/fiber-m3o-validator/mocks"
)

func TestTask_CreateTask(t *testing.T) {
	t.Run("Created", func(t *testing.T) {
		server := newServer(t, func(mockM3ODbClient *mocks.M3ODb) error {
			mockM3ODbClient.On("Create", mock.Anything).Return(&db.CreateResponse{Id: "task-id"}, nil)
			return nil
		})

		reqBody := fiber.Map{"task": fiber.Map{"title": "test title"}}
		resBody := fiber.Map{}
		res := sendRequest(t, server, fiber.MethodPost, "/tasks", reqBody, &resBody)
		if assert.Equal(t, fiber.StatusCreated, res.StatusCode, resBody) {
			assert.NotZero(t, resBody["id"])
		}
	})
}
