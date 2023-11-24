package Routes

import (
	"github.com/gofiber/fiber/v2"
	"sales-api/Controller"
)

func Setup(app *fiber.App) {
	// auth routes
	app.Post("/cashier/login", Controller.Login)
	app.Post("/cashier/:cashierId/logout", Controller.Logout)
	app.Post("/cashier/:cashierId/passcode", Controller.Passcode)

	// cashier routes
	app.Post("/cashiers", Controller.CreateCashier)
	app.Get("/cashiers", Controller.GetCashierList)
	app.Get("/cashiers/:cashierId", Controller.GetCashierDetails)
	app.Delete("/cashiers/:cashierId", Controller.DeleteCashier)
	app.Put("/cashiers/:cashierId", Controller.UpdateCashier)

	//Category routes
	app.Get("/categories", Controller.CategoryList)
	app.Get("/categories/:categoryId", Controller.GetCategoryDetails)
	app.Post("/categories", Controller.CreateCategory)
	app.Delete("/categories/:categoryId", Controller.DeleteCategory)
	app.Put("/categories/:categoryId", Controller.UpdateCategory)

	//Products routes
	app.Get("/products", Controller.ProductList)
	app.Get("/products/:productId", Controller.GetProductDetails)
	app.Post("/products", Controller.CreateProduct)
	app.Delete("/products/:productId", Controller.DeleteProduct)
	app.Put("/products/:productId", Controller.UpdateProduct)

	//Payment routes
	app.Get("/payments", Controller.PaymentList)
	app.Get("/payments/:paymentId", Controller.GetPaymentDetails)
	app.Post("/payments", Controller.CreatePayment)
	app.Delete("/payments/:paymentId", Controller.DeletePayment)
	app.Put("/payments/:paymentId", Controller.UpdatePayment)

	//Order routes
	app.Get("/orders", Controller.OrdersList)
	app.Get("/orders/:orderId", Controller.OrderDetail)
	app.Post("/orders", Controller.CreateOrder)
	app.Post("/orders/subtotal", Controller.SubTotalOrder)
	app.Get("/orders/:orderId/download", Controller.DownloadOrder)
	app.Get("/orders/:orderId/check-download", Controller.CheckOrder)

	//reports
	app.Get("/revenues", Controller.GetRevenues)
	app.Get("/solds", Controller.GetSolds)
}
