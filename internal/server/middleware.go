package server

import (
	"fmt"
	"strings"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func NewJWTAuthMiddleware(jwtProvider jwt_provider.Provider) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := string(c.Request().Header.Peek("Authorization"))
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			return SendErrorResponse(c, fiber.StatusUnauthorized, ErrAuthorizationHeader)
		}

		claims, err := jwtProvider.GetClaims(authHeaderParts[1])
		if err != nil {
			// TODO: Use logger
			fmt.Printf("Error: %s", err.Error())
			return SendErrorResponse(c, fiber.StatusUnauthorized, ErrAuthorizationHeader)
		}

		userUUID, err := claims.GetSubject()
		if err != nil {
			// TODO: Use logger
			fmt.Printf("Error: %s", err.Error())
			return SendErrorResponse(c, fiber.StatusUnauthorized, ErrAuthorizationHeader)
		}

		c.Locals("userUUID", userUUID)

		return c.Next()
	}
}

func GetAuthUser(c *fiber.Ctx, userRepo repository.UserRepository, opts ...any) (entity.User, error) {
	return userRepo.First(c.Context(), clause.Eq{Column: "uuid", Value: c.Locals("userUUID")}, opts)
}
