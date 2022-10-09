package server

import (
	"github.com/gofiber/fiber/v2"
)

func sendErrorResponse(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}

type CreateResponse struct {
	Id string `json:"id"`
}

func sendCreateResponse(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusCreated).JSON(CreateResponse{Id: id})
}
