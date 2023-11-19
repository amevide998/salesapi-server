package Routes

import (
	"github.com/gofiber/fiber/v2"
	"sales-api/Controller"
)

func Setup(app *fiber.App) {
	// auth routes
	app.Post("/cashier/:cashierId/login", Controller.Login)
	app.Get("/cashier/:cashierId/logout", Controller.Logout)
	app.Post("/cashier/:cashierId/passcode", Controller.Passcode)

	// cashier routes
	app.Post("/cashiers", Controller.CreateCashier)
	app.Get("/cashiers", Controller.GetCashierList)
	app.Get("/cashiers/:cashierId", Controller.GetCashierDetails)
	app.Delete("/cashiers/:cashierId", Controller.DeleteCashier)
	app.Put("/cashiers/:cashierId", Controller.UpdateCashier)
}
