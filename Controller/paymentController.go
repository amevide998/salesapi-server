package Controller

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"sales-api/Middleware"
	"sales-api/Model"
	"sales-api/Response"
	db "sales-api/config"
	"sales-api/dto"
	"strconv"
)

// payment controller
// @Description create payment
// @Summary create payment
// @Tags Payment
// @Produce json
// @Param request body dto.CreatePayment true "request"
// @Param Authorization header string true "authorization"
// @Success 201 {object} Response.WebResponse[Model.Payment]
// @Router /payments [post]
func CreatePayment(c *fiber.Ctx) error {

	var data dto.CreatePayment
	paymentError := c.BodyParser(&data)
	if paymentError != nil {
		log.Fatalf("payment error in post request %v", paymentError)
	}
	if data.Name == "" || data.Type == "" {
		response := Response.NewWebErrorResponse("Payment Name and Type is required")
		return c.Status(400).JSON(response)
	}

	var paymentTypes Model.PaymentType
	db.DB.Where("name", data.Type).First(&paymentTypes)
	payment := Model.Payment{
		Name:          data.Type,
		Type:          data.Type,
		PaymentTypeId: int(paymentTypes.Id),
		Logo:          data.Logo,
	}
	db.DB.Create(&payment)

	response := Response.NewWebResponse(payment)
	return c.Status(200).JSON(response)
}

// payment controller
// @Description list payment
// @Summary list payment
// @Tags Payment
// @Produce json
// @Param limit query string true "limit"
// @Param skip query string true "skip"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[dto.ListPayment]
// @Router /payments [get]
func PaymentList(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	//Token authenticate

	//subtotal,_ := strconv.Atoi(c.Query("subtotal"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var payment []dto.Payment
	db.DB.Select("id ,name,type,payment_type_id,logo,created_at,updated_at").Limit(limit).Offset(skip).Find(&payment).Count(&count)
	metaMap := dto.Pagination{
		Total: count,
		Limit: limit,
		Skip:  skip,
	}
	categoriesData := dto.ListPayment{
		Payments: payment,
		Meta:     metaMap,
	}

	response := Response.NewWebResponse(categoriesData)
	return c.JSON(response)

}

// payment controller
// @Description payment details
// @Summary payment details
// @Tags Payment
// @Produce json
// @Param paymentId path string true "payment id"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[Model.Payment]
// @Router /payments/{paymentId} [get]
func GetPaymentDetails(c *fiber.Ctx) error {
	paymentId := c.Params("paymentId")

	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	//Token authenticate

	var payment Model.Payment
	db.DB.Where("id=?", paymentId).First(&payment)

	if payment.Name == "" {
		response := Response.NewWebErrorResponse("payment not found")
		return c.Status(404).JSON(response)
	}

	response := Response.NewWebResponse(payment)
	return c.Status(200).JSON(response)
}

// payment controller
// @Description payment details
// @Summary payment details
// @Tags Payment
// @Produce json
// @Param paymentId path string true "payment id"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[string]
// @Router /payments/{paymentId} [delete]
func DeletePayment(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		response := Response.NewWebErrorResponse("Unauthorized")
		return c.Status(401).JSON(response)
	}
	//Token authenticate

	paymentId := c.Params("paymentId")
	var payment Model.Payment

	db.DB.First(&payment, paymentId)
	if payment.Name == "" {
		response := Response.NewWebErrorResponse("No payment found against this payment id")
		return c.Status(http.StatusNotFound).JSON(response)
	}

	result := db.DB.Delete(&payment)
	if result.RowsAffected == 0 {
		response := Response.NewWebErrorResponse("payment removing failed")
		return c.Status(http.StatusNotFound).JSON(response)
	}

	response := Response.NewWebResponse("")
	return c.Status(200).JSON(response)
}

// payment controller
// @Description payment details
// @Summary payment details
// @Tags Payment
// @Produce json
// @Param paymentId path string true "payment id"
// @Param request body Model.Payment true "request"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[Model.Payment]
// @Router /payments/{paymentId} [put]
func UpdatePayment(c *fiber.Ctx) error {
	paymentId := c.Params("paymentId")
	var totalPayment Model.Payment
	db.DB.Find(&totalPayment)

	var payment Model.Payment

	db.DB.Find(&payment, "id = ?", paymentId)

	var updatePaymentData Model.Payment
	err := c.BodyParser(&updatePaymentData)
	if err != nil {
		return err
	}
	if updatePaymentData.Name == "" {
		response := Response.NewWebErrorResponse("Payment name is required")
		return c.Status(400).JSON(response)
	}

	var paymentTypeId int
	if updatePaymentData.Type == "CASH" {
		paymentTypeId = 1
	}
	if updatePaymentData.Type == "E-WALLET" {
		paymentTypeId = 2
	}
	if updatePaymentData.Type == "EDC" {
		paymentTypeId = 3
	}

	payment.Name = updatePaymentData.Name
	payment.Type = updatePaymentData.Type
	payment.PaymentTypeId = paymentTypeId
	payment.Logo = updatePaymentData.Logo

	db.DB.Save(&payment)
	response := Response.NewWebResponse(payment)
	return c.Status(200).JSON(response)

}
