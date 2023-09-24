package server

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"example.com/the-boring-to-do-list-1/pkg/env"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

type Server struct {
	fiberApp       *fiber.App
	AuthController *authController
	TaskController *taskController
}

func NewServer(jwt_provider jwt_provider.Provider, gorm_provider gorm_provider.AbstractProvider) Server {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// TODO: Add logger
			fmt.Printf("Error: %s", err.Error())
			return sendDefaultStatusResponse(c, fiber.StatusInternalServerError)
		},
	})
	fiberApp.Use(logger.New(logger.Config{Format: "[${time} ${latency}] ${status} ${method} ${path}\n"}))
	if env.GetOrDefault("ENV", "production") != "test" {
		fiberApp.Use(recover.New())
	}
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	return Server{
		fiberApp:       fiberApp,
		AuthController: newAuthController(fiberApp, jwt_provider, gorm_provider),
		TaskController: newTaskController(fiberApp, gorm_provider),
	}
}

func (s Server) Run() error {
	return s.fiberApp.Listen("0.0.0.0:3000")
}

func (s Server) Test(req *http.Request) (*http.Response, error) {
	return s.fiberApp.Test(req, -1)
}
