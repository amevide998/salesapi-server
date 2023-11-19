package main

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	db "sales-api/config"
)

func main() {

	//db
	db.Connect()

	// fiber instance
	app := fiber.New()

	// http handler
	app.Get("/testApi", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "test api",
		})
	})

	// listen on port
	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
