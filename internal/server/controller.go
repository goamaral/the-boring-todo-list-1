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

func sendCreatedResponse(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusCreated).JSON(CreateResponse{Id: id})
}

func parseBody(c *fiber.Ctx, body any) error {
	if c.Request().Header.ContentLength() != 0 {
		err := c.BodyParser(body)
		if err != nil {
			return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
		}
	}
	return nil
}
