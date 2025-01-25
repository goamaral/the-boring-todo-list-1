package server

import (
	"errors"
	"fmt"
	"time"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/env"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm/clause"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserAlreadyExists = errors.New("user already exists")

type authController struct {
	controller
	jwtProvider jwt_provider.Provider
	UserRepo    repository.UserRepository
}

func newAuthController(baseRouter fiber.Router, jwtProvider jwt_provider.Provider, gorm_provider gorm_provider.AbstractProvider) *authController {
	ctrl := &authController{
		controller:  newController(),
		jwtProvider: jwtProvider,
		UserRepo:    repository.NewUserRepository(gorm_provider),
	}

	router := baseRouter.Group("/auth")
	router.Post("/login", ctrl.Login)
	router.Get("/register", ctrl.NewRegister)
	router.Post("/register", ctrl.Register)

	return ctrl
}

func GenerateAccessToken(jwtProvider jwt_provider.Provider, userUUID string) (string, error) {
	return jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   userUUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(15 * time.Minute)},
	})
}

func GenerateRefreshToken(jwtProvider jwt_provider.Provider, userUUID string) (string, error) {
	return jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   userUUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)},
	})
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (ct *authController) Login(c *fiber.Ctx) error {
	// Parse request
	req := LoginRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Validate request
	err = ct.validate.Struct(req)
	if err != nil {
		return SendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Get user password
	user, found, err := ct.UserRepo.FindOne(
		c.Context(),
		gorm_provider.SelectOption("encrypted_password"),
		clause.Eq{Column: "username", Value: req.Username},
	)
	if err != nil {
		return err
	}
	if !found {
		return SendErrorResponse(c, fiber.StatusBadRequest, ErrInvalidCredentials)
	}

	// Compare password
	ok, err := user.ComparePassword(req.Password)
	if err != nil {
		return err
	}
	if !ok {
		return SendErrorResponse(c, fiber.StatusBadRequest, ErrInvalidCredentials)
	}

	return ct.afterAuthenticate(c, user.UUID)
}

func (ct *authController) NewRegister(c *fiber.Ctx) error {
	return c.Render("auth/register", nil)
}

type RegisterRequest struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"eqfield=Password"`
}

func (ct *authController) Register(c *fiber.Ctx) error {
	// Parse request
	req := RegisterRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return SendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}
	fmt.Println(req)

	// Validate request
	err = ct.validate.Struct(req)
	if err != nil {
		return SendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Create user
	user := entity.User{Username: req.Username}
	err = user.SetEncryptedPassword(req.Password)
	if err != nil {
		return err
	}
	err = ct.UserRepo.Create(c.Context(), &user)
	if err != nil {
		if gorm_provider.HasErrorCode(err, gorm_provider.UniqueConstraintViolation) {
			return SendErrorResponse(c, fiber.StatusBadRequest, ErrUserAlreadyExists)
		}
		return err
	}

	return ct.afterAuthenticate(c, user.UUID)
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("accessToken")
	c.ClearCookie("refreshToken")
	return c.Redirect("/auth/login")
}

func (ct *authController) afterAuthenticate(c *fiber.Ctx, userUUID string) error {
	// Generate JWT access token
	accessToken, err := GenerateAccessToken(ct.jwtProvider, userUUID)
	if err != nil {
		return err
	}

	// Generate JWT refresh token
	refreshToken, err := GenerateRefreshToken(ct.jwtProvider, userUUID)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Domain:   c.Hostname(),
		Secure:   env.Get("ENV") == "production",
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Domain:   c.Hostname(),
		Secure:   env.Get("ENV") == "production",
		HTTPOnly: true,
	})

	return c.Redirect("/tasks")
}
