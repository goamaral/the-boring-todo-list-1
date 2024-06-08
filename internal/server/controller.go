package server

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type controller struct {
	validate *validator.Validate
}

func newController() controller {
	return controller{
		validate: validator.New(),
	}
}

/* ERRORS */
func SendErrorResponse(c *fiber.Ctx, status int, err error) error {
	// TODO: Use logger
	fmt.Printf("Error: %s", err.Error())
	return c.Status(status).JSON(fiber.Map{"error": err.Error()})
}

func SendValidationErrorsResponse(c *fiber.Ctx, errs validator.ValidationErrors) error {
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
				currMap[err.Tag()] = decapitalize(err.Param())
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
