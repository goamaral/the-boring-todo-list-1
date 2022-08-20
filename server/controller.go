package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func sendErrorResponse(c *fiber.Ctx, status int, err error, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": errors.Wrap(err, message).Error()})
}

func sendCreatedResponse(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}
