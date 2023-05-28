package server

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/jwtprovider"
)

type server struct {
	fiberApp       *fiber.App
	authController *authController
	taskController *taskController
}

type Server interface {
	Run() error
	Test(req *http.Request) (resp *http.Response, err error)
}

func NewServer(jwtProvider *jwtprovider.Provider, taskRepo repository.TaskRepository, userRepo repository.UserRepository) *server {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// TODO: Add logger
			fmt.Printf("Error: %s", err.Error())
			return sendDefaultStatusResponse(c, fiber.StatusInternalServerError)
		},
	})
	fiberApp.Use(logger.New(logger.Config{Format: "[${time} ${latency}] ${status} ${method} ${path}\n"}))
	fiberApp.Use(recover.New())
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	s := &server{
		fiberApp:       fiberApp,
		authController: newAuthController(fiberApp, jwtProvider, userRepo),
		taskController: newTaskController(fiberApp, taskRepo),
	}

	return s
}

/* PUBLIC */
func (s server) Run() error {
	return s.fiberApp.Listen("0.0.0.0:3000")
}

func (s server) Test(req *http.Request) (*http.Response, error) {
	return s.fiberApp.Test(req, -1)
}
