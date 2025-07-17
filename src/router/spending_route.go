package router

import (
	"app/src/controller"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func SpendingRoutes(r fiber.Router, spendingService *service.SpendingService) {
	spendingController := controller.NewSpendingController(*spendingService)
	spending := r.Group("/spending")

	spending.Post("/", func(c *fiber.Ctx) error {
		return spendingController.CreateSpending(c)
	})

	spending.Get("/list", func(c *fiber.Ctx) error {
		return spendingController.GetSpending(c)
	})

	spending.Get("/categories", func(c *fiber.Ctx) error {
		return spendingController.GetCategories(c)
	})

	spending.Get("/summary", func(c *fiber.Ctx) error {
		return spendingController.GetSummarySpending(c)
	})
}
