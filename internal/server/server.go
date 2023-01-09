package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"example.com/the-boring-to-do-list-1/internal/repository"
)

type server struct {
	fiberApp       *fiber.App
	taskController *taskController
}

type Server interface {
	Run() error
	Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error)
}

func NewServer(taskRepo repository.TaskRepository) *server {
	fiberApp := fiber.New()
	fiberApp.Use(logger.New(logger.Config{Format: "[${time} ${latency}] ${status} ${method} ${path}\n"}))
	fiberApp.Use(recover.New())
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	s := &server{
		fiberApp:       fiberApp,
		taskController: newTaskController(fiberApp, taskRepo),
	}

	return s
}

/* PUBLIC */
func (s server) Run() error {
	return s.fiberApp.Listen("0.0.0.0:3000")
}

func (s server) Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error) {
	return s.fiberApp.Test(req, msTimeout...)
}
