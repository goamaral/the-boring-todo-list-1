package server

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type absctractController struct {
	validate *validator.Validate
}

func newAbstractController() absctractController {
	return absctractController{
		validate: validator.New(),
	}
}

func sendDefaultStatusResponse(c *fiber.Ctx, status int) error {
	var msg string

	switch status {
	case http.StatusOK:
		msg = "OK"
	case http.StatusNotFound:
		msg = "Not found"
	case http.StatusInternalServerError:
		msg = "Internal error"
	default:
		return errors.Errorf("http status (%d) has no default response defined", status)
	}

	return c.Status(status).JSON(msg)
}

type CreateResponse struct {
	Id string `json:"id"`
}

func sendCreatedResponse(c *fiber.Ctx, id string) error {
	return c.Status(fiber.StatusCreated).JSON(CreateResponse{Id: id})
}

/* ERRORS */
func sendErrorResponse(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}

func sendValidationErrorsResponse(c *fiber.Ctx, errs validator.ValidationErrors) error {
	res := fiber.Map{}
	var parent fiber.Map

	for _, err := range errs {
		nsParts := strings.Split(err.Namespace(), ".")
		for i, nsPart := range nsParts {
			if i == 0 {
				parent = res
				continue
			}

			nsPart := decapitalize(nsPart)

			curr := parent[nsPart]
			currMap := fiber.Map{}
			if curr != nil {
				currMap = curr.(fiber.Map)
			}
			if i == len(nsParts)-1 {
				currMap[err.Tag()] = err.Param()
			}
			parent[nsPart] = currMap
			parent = currMap
		}
	}

	return c.Status(fiber.StatusBadRequest).JSON(res)
}

func decapitalize(s string) string {
	ns := []byte(s)
	if len(s) > 0 && ns[0] < 'a' {
		lowerCaseOffset := byte(32) // 'a' - 'A'
		ns[0] = ns[0] + lowerCaseOffset
	}
	return string(ns)
}
