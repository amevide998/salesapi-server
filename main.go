package main

import (
	"github.com/gofiber/fiber/v2"
	routes "sales-api/Routes"
	db "sales-api/config"
)

func main() {

	//db
	db.Connect()

	// fiber instance
	app := fiber.New()
	//app.Use()

	routes.Setup(app)

	// listen on port
	err := app.Listen(":8000")
	if err != nil {
		return
	}
}
