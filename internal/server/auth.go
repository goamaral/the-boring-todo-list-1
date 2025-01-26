package server

import (
	"time"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/env"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

var ErrInvalidCredentials = fiber.NewError(fiber.StatusBadRequest, "invalid credentials")
var ErrUserAlreadyExists = fiber.NewError(fiber.StatusUnprocessableEntity, "user already exists")

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
	router.Get("/login", ctrl.NewLogin)
	router.Post("/login", ctrl.Login)
	router.Get("/register", ctrl.NewRegister)
	router.Post("/register", ctrl.Register)
	router.Get("/logout", ctrl.Logout)

	return ctrl
}

func GenerateAccessToken(jwtProvider jwt_provider.Provider, userUUID uuid.UUID) (string, error) {
	return jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   userUUID.String(),
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(15 * time.Minute)},
	})
}

func GenerateRefreshToken(jwtProvider jwt_provider.Provider, userUUID uuid.UUID) (string, error) {
	return jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   userUUID.String(),
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)},
	})
}

func (ct *authController) NewLogin(c *fiber.Ctx) error {
	return c.Render("auth/login", nil)
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
		return err
	}

	// Validate request
	err = ct.validate.Struct(req)
	if err != nil {
		return SendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Get user password
	user, found, err := ct.UserRepo.FindOne(
		c.Context(),
		gorm_provider.SelectOption("uuid", "encrypted_password"),
		clause.Eq{Column: "username", Value: req.Username},
	)
	if err != nil {
		return err
	}
	if !found {
		return ErrInvalidCredentials
	}

	// Compare password
	ok, err := user.ComparePassword(req.Password)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidCredentials
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
		return err
	}

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
			return ErrUserAlreadyExists
		}
		return err
	}

	return ct.afterAuthenticate(c, user.UUID)
}

func (ct *authController) Logout(c *fiber.Ctx) error {
	c.ClearCookie("accessToken")
	c.ClearCookie("refreshToken")
	return c.Redirect("/auth/login")
}

func (ct *authController) afterAuthenticate(c *fiber.Ctx, userUUID uuid.UUID) error {
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
		Secure:   env.Get("ENV") == "production",
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Secure:   env.Get("ENV") == "production",
		HTTPOnly: true,
	})

	return c.Redirect("/tasks")
}
