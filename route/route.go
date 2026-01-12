package route

import (
	logcustom "utility/config/log"
	"utility/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetUpRoute() *fiber.App {
	app1 := fiber.New(fiber.Config{
		ErrorHandler: logcustom.NewZapErrorHandler("logs"),
	})
	app1.Use(logcustom.NewRequestResponseLogger("logs"))
	// app1.Get("/get", func(cx *fiber.Ctx) error {
	// 	return cx.Status(200).JSON(fiber.Map{
	// 		"status": true,
	// 	})
	// })
	api := app1.Group("/api")
	apiCard := api.Group("/card")
	apiCard.Post("/get-brand-card", func(cx *fiber.Ctx) error {
		return cx.Status(200).JSON(fiber.Map{
			"status": true,
		})
	})
	app1.Use(recover.New())
	app1.Use(NotFoundRoute)
	return app1
}

func NotFoundRoute(c *fiber.Ctx) error {
	return c.Status(404).JSON(helper.StructMasterResponse{
		Status:     false,
		StatusCode: "R0001",
		Message:    "No Register Route",
	}.StructMasterResponseFinal())
}
