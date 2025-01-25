package server

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/handlebars/v2"

	"example.com/the-boring-to-do-list-1/pkg/env"
	"example.com/the-boring-to-do-list-1/pkg/fs"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

type Server struct {
	fiberApp       *fiber.App
	AuthController *authController
	TaskController *taskController
}

func NewServer(jwtProvider jwt_provider.Provider, gormProvider gorm_provider.AbstractProvider) Server {
	viewEngine := handlebars.New(fs.ResolveRelativePath("../views"), ".hbs")
	viewEngine.LayoutName = "yield"
	if env.GetOrDefault("ENV", "production") != "production" {
		viewEngine.Reload(true)
	}

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// TODO: Use logger
			fmt.Printf("Error: %s", err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		},
		Views:       viewEngine,
		ViewsLayout: "layouts/public",
	})
	fiberApp.Use(logger.New(logger.Config{Format: "[${time} ${latency}] ${status} ${method} ${path}\n"}))
	if env.GetOrDefault("ENV", "production") != "test" {
		fiberApp.Use(recover.New())
	}
	fiberApp.Static("/static", "./public")
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	return Server{
		fiberApp:       fiberApp,
		AuthController: newAuthController(fiberApp, jwtProvider, gormProvider),
		TaskController: newTaskController(fiberApp, jwtProvider, gormProvider),
	}
}

func (s Server) Run() error {
	return s.fiberApp.Listen("0.0.0.0:3000")
}

func (s Server) Test(req *http.Request) (*http.Response, error) {
	return s.fiberApp.Test(req, -1)
}
