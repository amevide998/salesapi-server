package Controller

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"sales-api/Model"
	db "sales-api/config"
	"strconv"
	"time"
)

func CreateCashier(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "invalid input",
			})
	}

	if data["name"] == "" {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "name is required",
			})
	}

	if data["passcode"] == "" {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "passcode is required",
			})
	}

	// saving to database
	cashier := Model.Cashier{
		Name:      data["name"],
		Passcode:  data["passcode"],
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	db.DB.Create(&cashier)
	return c.Status(http.StatusCreated).JSON(
		fiber.Map{
			"status":  true,
			"message": "success created cashier",
			"data":    cashier.Name,
		})

}

func UpdateCashier(c *fiber.Ctx) error {
	var cashier Model.Cashier

	cashierId := c.Params("cashierId")

	db.DB.Find(&cashier, "id = ? ", cashierId)

	if cashier.Name == "" {
		return c.Status(http.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "cashier not found",
			})
	}

	var updateCashier Model.Cashier

	err := c.BodyParser(&updateCashier)
	if err != nil {
		return err
	}

	if updateCashier.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"success": false,
				"message": "name is required",
			})
	}

	cashier.Name = updateCashier.Name

	db.DB.Save(&cashier)

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "success updated cashier",
			"data":    cashier.Name,
		})
}

func GetCashierList(c *fiber.Ctx) error {
	var cashier []Model.Cashier

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))

	var count int64
	db.DB.Select("*").Limit(limit).Offset(skip).Find(&cashier).Count(&count)
	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "success get cashier list",
			"data":    cashier,
		})
}

func DeleteCashier(c *fiber.Ctx) error {
	var cashier Model.Cashier

	cashierId := c.Params("cashierId")

	db.DB.Where("id = ?", cashierId).First(&cashier)

	if cashier.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "cashier not found",
			})
	}

	db.DB.Where("id = ?", cashierId).Delete(&cashier)
	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "success delete cashier",
			"data":    cashier.Id,
		})

}

func GetCashierDetails(c *fiber.Ctx) error {
	var cashier Model.Cashier
	cashierId := c.Params("cashierId")

	db.DB.Select("*").Where("id = ?", cashierId).First(&cashier)

	if cashier.Id == 0 {
		return c.Status(http.StatusNotFound).JSON(
			fiber.Map{
				"success": false,
				"message": "Cashier Not Found",
				"error":   map[string]interface{}{},
			})
	}

	cashierData := make(map[string]interface{})

	cashierData["cashierId"] = cashier.Id
	cashierData["name"] = cashier.Name
	cashierData["createdAt"] = cashier.CreatedAt
	cashierData["updatedAt"] = cashier.UpdatedAt

	return c.Status(http.StatusOK).JSON(
		fiber.Map{
			"success": true,
			"message": "success get cashier data",
			"data":    cashierData,
		})

}
