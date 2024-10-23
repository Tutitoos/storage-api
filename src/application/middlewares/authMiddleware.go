package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"net/http"
	"storage-api/src/domain"
)

func AuthMiddleware(ctx fiber.Ctx) error {
	result := domain.ResultData[string]()

	authToken := ctx.Get("Authorization")
	if authToken == "" {
		result.AddMessage("Authorization token is missing")

		domain.Logger.Error("Authorization token is missing")

		return ctx.Status(http.StatusUnauthorized).JSON(result)
	}

	if authToken != domain.CONFIG.Token {
		result.AddMessage("Authorization token is invalid")

		domain.Logger.Error("Authorization token is invalid")

		return ctx.Status(http.StatusUnauthorized).JSON(result)
	}

	return ctx.Next()
}
