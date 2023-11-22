package Controller

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"sales-api/Model"
	"sales-api/Response"
	db "sales-api/config"
	"sales-api/dto"
	"strconv"
	"time"
)

// cashier controller
// @Description create cashier
// @Summary create cashier
// @Tags cashiers
// @Produce json
// @Param request body Model.Cashier true "request"
// @Success 201 {object} Response.WebResponse[string]
// @Router /cashiers [post]
func CreateCashier(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		response := Response.NewWebErrorResponse("invalid input")
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	if data["name"] == "" {
		response := Response.NewWebErrorResponse("name is required")
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	if data["passcode"] == "" {
		response := Response.NewWebErrorResponse("passcode is required")
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	// saving to database
	cashier := Model.Cashier{
		Name:      data["name"],
		Passcode:  data["passcode"],
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	db.DB.Create(&cashier)

	response := Response.NewWebResponse(cashier.Name)
	return c.Status(http.StatusCreated).JSON(response)

}

// cashier controller
// @Description update cashier
// @Summary get update cashier
// @Tags cashiers
// @Produce json
// @Param request body Model.Cashier true "request"
// @Success 200 {object} Response.WebResponse[string]
// @Router /cashiers/{cashierId} [put]
func UpdateCashier(c *fiber.Ctx) error {
	var cashier Model.Cashier

	cashierId := c.Params("cashierId")

	db.DB.Find(&cashier, "id = ? ", cashierId)

	if cashier.Name == "" {
		response := Response.NewWebResponse("cashier not found")
		return c.Status(http.StatusNotFound).JSON(response)
	}

	var updateCashier Model.Cashier

	err := c.BodyParser(&updateCashier)
	if err != nil {
		return err
	}

	if updateCashier.Name == "" {
		response := Response.NewWebResponse("name is required")
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	cashier.Name = updateCashier.Name

	db.DB.Save(&cashier)

	response := Response.NewWebResponse(cashier.Name, "success updated cashier")

	return c.Status(http.StatusOK).JSON(response)
}

// cashier controller
// @Description get cashier list
// @Summary get cashier list
// @Tags cashiers
// @Produce json
// @Success 200 {object} Response.WebResponse[[]Model.Cashier]
// @Router /cashiers [get]
func GetCashierList(c *fiber.Ctx) error {
	var cashier []Model.Cashier

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	var count int64
	db.DB.Select("*").Limit(limit).Offset(skip).Find(&cashier).Count(&count)

	response := Response.NewWebResponse(cashier)

	return c.Status(http.StatusOK).JSON(response)
}

// cashier controller
// @Description delete cashier
// @Summary delete cashier
// @Tags cashiers
// @Produce json
// @Param cashierId path string true "cashier id"
// @Success 200 {object} Response.WebResponse[string]
// @Router /cashiers/{cashierId} [delete]
func DeleteCashier(c *fiber.Ctx) error {
	var cashier Model.Cashier

	cashierId := c.Params("cashierId")

	db.DB.Where("id = ?", cashierId).First(&cashier)

	if cashier.Id == 0 {
		response := Response.NewWebErrorResponse("cashier not found")
		return c.Status(http.StatusNotFound).JSON(response)
	}

	db.DB.Where("id = ?", cashierId).Delete(&cashier)

	response := Response.NewWebResponse(cashier.Id)
	return c.Status(http.StatusOK).JSON(response)

}

// cashier controller
// @Description get cashier detail
// @Summary get cashier detail
// @Tags cashiers
// @Produce json
// @Param cashierId path string true "cashier id"
// @Success 200 {object} Response.WebResponse[dto.CashierDetails]
// @Router /cashiers/{cashierId} [get]
func GetCashierDetails(c *fiber.Ctx) error {
	var cashier Model.Cashier
	cashierId := c.Params("cashierId")

	db.DB.Select("*").Where("id = ?", cashierId).First(&cashier)

	if cashier.Id == 0 {
		response := Response.NewWebResponse("cashier not found")
		return c.Status(http.StatusNotFound).JSON(response)
	}

	cashierData := dto.CashierDetails{
		CashierId: int(cashier.Id),
		Name:      cashier.Name,
		CreatedAt: cashier.CreatedAt,
		UpdatedAt: cashier.UpdatedAt,
	}

	response := Response.NewWebResponse(cashierData)
	return c.Status(http.StatusOK).JSON(response)

}
