package server

import (
	"errors"
	"strings"
	"time"

	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/gormprovider"
	"example.com/the-boring-to-do-list-1/pkg/jwtprovider"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var ErrAuthorizationHeader = errors.New("authorization header is missing/invalid")
var ErrInvalidCredentials = errors.New("invalid credentials")

type authController struct {
	absctractController
	jwtProvider *jwtprovider.Provider
	userRepo    repository.UserRepository
}

func newAuthController(baseRouter fiber.Router, jwtProvider *jwtprovider.Provider, userRepo repository.UserRepository) *authController {
	ctrl := &authController{
		absctractController: newAbstractController(),
		jwtProvider:         jwtProvider,
		userRepo:            userRepo,
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

func (ct authController) Login(c *fiber.Ctx) error {
	// Parse request
	req := LoginRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Validate request
	err = ct.validate.Struct(req)
	if err != nil {
		return sendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Get user password
	user, found, err := ct.userRepo.Get(c.Context(), repository.UserFilter{Username: &req.Username}, gormprovider.SelectOption("id", "username"))
	if err != nil {
		return err
	}
	if !found {
		return sendErrorResponse(c, fiber.StatusBadRequest, ErrInvalidCredentials)
	}

	// Compare password
	ok, err := user.ComparePassword(req.Password)
	if err != nil {
		return err
	}
	if !ok {
		return sendErrorResponse(c, fiber.StatusBadRequest, ErrInvalidCredentials)
	}

	// Generate JWT access token
	accessToken, err := ct.jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   user.Id,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(15 * time.Minute)},
	})
	if err != nil {
		return sendDefaultStatusResponse(c, fiber.StatusInternalServerError)
	}

	// Generate JWT refresh token
	refreshToken, err := ct.jwtProvider.GenerateSignedToken(jwt.RegisteredClaims{
		Subject:   "",
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)},
	})
	if err != nil {
		return sendDefaultStatusResponse(c, fiber.StatusInternalServerError)
	}

	return c.JSON(LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func (ac authController) JWTAuthMiddleware(c *fiber.Ctx) error {
	authHeader := string(c.Request().Header.Peek("Authorization"))
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return sendErrorResponse(c, fiber.StatusUnauthorized, ErrAuthorizationHeader)
	}

	claims, err := ac.jwtProvider.GetClaims(authHeaderParts[1])
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, ErrAuthorizationHeader)
	}

	userId, err := claims.GetSubject()
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnauthorized, ErrAuthorizationHeader)
	}

	c.Locals("userId", userId)

	return c.Next()
}
