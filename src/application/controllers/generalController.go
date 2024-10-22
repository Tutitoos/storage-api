package controllers

import (
	"github.com/gofiber/fiber/v3"
	"storage-api/src/domain"
)

type IGeneralController struct {
}

func GeneralController() *IGeneralController {
	return &IGeneralController{}
}

func (c *IGeneralController) GetHomeHandler(ctx fiber.Ctx) error {
	result := domain.ResultData[string]()

	result.AddMessage("API is up and running!")

	return ctx.JSON(result)
}
