package application

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"storage-api/src/application/middlewares"
	"storage-api/src/application/routers"
	"storage-api/src/domain"
)

func Api() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   10 << 20,
	})

	app.Use(logger.New(logger.Config{
		Format:     "${time} :: ${ip} :: ${ips} :: ${status} :: ${method} ${path} :: ${latency}\n",
		TimeFormat: "02-02-2006 15:04:05",
		TimeZone:   "Europe/Madrid",
	}))

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"https://s3-elcoheteboom.tutitoos.xyz"},
		AllowMethods:     []string{"GET,POST,PUT,DELETE,OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "application/json", "multipart/form-data"},
	}))

	app.Use(middlewares.AuthMiddleware)

	router := app.Group("/v1")

	routers.GeneralRouter(router)
	routers.CloudflareRouter(router)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", domain.CONFIG.Port)))
}
