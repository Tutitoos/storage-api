package routers

import (
	"github.com/gofiber/fiber/v3"
	"storage-api/src/application/controllers"
)

func CloudflareRouter(router fiber.Router) fiber.Router {
	controller := controllers.CloudflareController()

	router.Get("/", controller.GetHomeHandler)
	router.Get("/files/*", controller.GetFilesHandler)
	router.Get("/file/*", controller.GetFileHandler)
	router.Delete("/file/*", controller.DeleteFileHandler)
	router.Post("/file", controller.UploadFileHandler)

	return router
}
