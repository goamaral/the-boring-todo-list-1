package server

import "github.com/gofiber/fiber/v2"

type server struct {
	fiberApp *fiber.App
}

type Server interface {
}

func NewServer() *server {
	s := &server{fiberApp: fiber.New()}
	s.setupRoutes()

	return s
}

/* PUBLIC */
func (s server) Run() error {
	return s.fiberApp.Listen(":3000")
}

/* PRIVATE */
func (s server) setupRoutes() {
	s.fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
