package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	m3o "go.m3o.com"
)

type server struct {
	fiberApp       *fiber.App
	taskController *taskController
}

type Server interface {
	Run() error
	Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error)
}

func NewServer(m3oClient *m3o.Client) *server {
	fiberApp := fiber.New()

	s := &server{
		fiberApp:       fiberApp,
		taskController: newTaskController(fiberApp.Group("/tasks"), m3oClient.Db),
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
