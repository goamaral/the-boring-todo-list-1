package server

import (
	"errors"
	"time"

	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm/clause"
)

var ErrAuthorizationHeader = errors.New("authorization header is missing/invalid")
var ErrInvalidCredentials = errors.New("invalid credentials")

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

	return ctrl
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

	// Generate JWT access token
	accessToken, err := ct.jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(15 * time.Minute)},
	})
	if err != nil {
		return SendDefaultStatusResponse(c, fiber.StatusInternalServerError)
	}

	// Generate JWT refresh token
	refreshToken, err := ct.jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)},
	})
	if err != nil {
		return SendDefaultStatusResponse(c, fiber.StatusInternalServerError)
	}

	return c.JSON(LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
