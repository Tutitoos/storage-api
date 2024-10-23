package application

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"net/http"
	"storage-api/src/application/middlewares"
	"storage-api/src/application/routers"
	"storage-api/src/domain"
	"time"
)

func Api() {
	if err := domain.CustomLogger("logs/app.log"); err != nil {
		log.Fatalf("error al crear logger: %v", err)
	}
	defer domain.Logger.Close()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   10 << 20,
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			result := domain.ResultData[string]()
			result.AddError(http.StatusInternalServerError, err.Error())
			return ctx.Status(http.StatusInternalServerError).JSON(result)
		},
	})

	app.Use(recover.New())

	app.Use(middlewares.IPWhitelistMiddleware)

	app.Use(logger.New(logger.Config{
		Format:     "${time} :: ${ip} :: ${ips} :: ${status} :: ${method} ${path} :: ${latency}\n",
		TimeFormat: "02-01-2006 03:04:05 PM",
		TimeZone:   "Europe/Madrid",
	}))

	app.Use(func(ctx fiber.Ctx) error {
		start := time.Now()
		err := ctx.Next()
		latency := time.Since(start)

		logMessage := fmt.Sprintf("%s :: %s :: %s :: %s %s :: %d :: %s",
			ctx.IP(),
			ctx.IPs(),
			ctx.Method(),
			ctx.Path(),
			ctx.OriginalURL(),
			ctx.Response().StatusCode(),
			latency,
		)

		domain.Logger.Debug(logMessage)

		if ctx.Response().StatusCode() != http.StatusOK {
			responseData := ctx.Response().Body()
			return ctx.SendString(string(responseData))
		}

		if err != nil {
			result := domain.ResultData[string]()
			result.AddError(ctx.Response().StatusCode(), err.Error())
			return ctx.Status(ctx.Response().StatusCode()).JSON(result)
		}

		return ctx.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"https://elcoheteboom.com"},
		AllowMethods:     []string{"GET,POST,PUT,DELETE,OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "application/json", "multipart/form-data"},
	}))

	//app.Use(middlewares.IPWhitelistMiddleware)
	app.Use(middlewares.AuthMiddleware)

	router := app.Group("/v1")

	routers.GeneralRouter(router)
	routers.CloudflareRouter(router)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", domain.CONFIG.Port)))
}
