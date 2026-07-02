package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/dto"
	"github.com/thawatchai/full-stack-authentication/back-end/internal/service"
)

const (
	UserIDKey = "userID"
	EmailKey  = "email"
)

func Auth(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Error: "missing authorization header",
				Code:  "AUTH_MISSING",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Error: "invalid authorization header",
				Code:  "AUTH_INVALID",
			})
		}

		claims, err := authService.VerifyAccessToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Error: "invalid or expired token",
				Code:  "AUTH_INVALID",
			})
		}

		c.Locals(UserIDKey, claims.UserID)
		c.Locals(EmailKey, claims.Email)
		return c.Next()
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"
	errorCode := "INTERNAL_ERROR"

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	return c.Status(code).JSON(dto.ErrorResponse{
		Error: message,
		Code:  errorCode,
	})
}
