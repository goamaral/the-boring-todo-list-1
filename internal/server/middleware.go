package server

import (
	"fmt"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func NewJWTAuthMiddleware(jwtProvider jwt_provider.Provider) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies("accessToken")
		fmt.Println("accessToken", accessToken)
		if accessToken == "" {
			return Logout(c)
		}

		claims, err := jwtProvider.GetClaims(accessToken)
		if err != nil {
			// TODO: If expired, try to use refresh token
			// TODO: Use logger
			fmt.Printf("Error: %s", err.Error())
			return Logout(c)
		}

		userUUID, err := claims.GetSubject()
		if err != nil {
			// TODO: Use logger
			fmt.Printf("Error: %s", err.Error())
			return Logout(c)
		}

		c.Locals("userUUID", userUUID)

		return c.Next()
	}
}

func GetAuthUser(c *fiber.Ctx, userRepo repository.UserRepository, opts ...any) (entity.User, error) {
	return userRepo.First(c.Context(), clause.Eq{Column: "uuid", Value: c.Locals("userUUID")}, opts)
}
