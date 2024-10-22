package routers

import (
	"github.com/gofiber/fiber/v3"
	"storage-api/src/application/controllers"
)

func GeneralRouter(router fiber.Router) fiber.Router {
	controller := controllers.GeneralController()

	router.Get("/", controller.GetHomeHandler)

	return router
}
