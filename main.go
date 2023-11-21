package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	routes "sales-api/Routes"
	db "sales-api/config"
	_ "sales-api/docs"
)

// @title Sales api docs
// @version 1.0
// @description This is the complete api documentation for sales api
// @host localhost:8000
func main() {

	//db
	db.Connect()

	// fiber instance
	app := fiber.New()

	// Initialize default config
	app.Use(cors.New())

	//Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// swagger
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	//app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
	//	URL:         "http://example.com/doc.json",
	//	DeepLinking: false,
	//	// Expand ("list") or Collapse ("none") tag groups by default
	//	DocExpansion: "none",
	//	// Prefill OAuth ClientId on Authorize popup
	//	OAuth: &swagger.OAuthConfig{
	//		AppName:  "OAuth Provider",
	//		ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
	//	},
	//	// Ability to change OAuth2 redirect uri location
	//	OAuth2RedirectUrl: "http://localhost:8000/swagger/oauth2-redirect.html",
	//}))

	routes.Setup(app)

	// listen on port
	err := app.Listen(":8000")
	if err != nil {
		return
	}
}
