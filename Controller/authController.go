package Controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"sales-api/Model"
	"sales-api/Response"
	db "sales-api/config"
	"sales-api/dto"
	"strconv"
	"time"
)

// login controller
// @Description login Cashier
// @Summary login cashier
// @Tags authentication
// @Produce json
// @Param cashierId path string true "cashier Id"
// @Param request body dto.Passcode true "request"
// @Success 200 {object} Response.WebResponse[dto.Token]
// @Router /cashier/{cashierId}/login [post]
func Login(c *fiber.Ctx) error {
	cashierId := c.Params("cashierId")
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		response := Response.NewWebErrorResponse("invalid post request")
		return c.Status(400).JSON(response)
	}

	//check if passcode is empty
	if data["passcode"] == "" {
		response := Response.NewWebErrorResponse("passcode cannot be empty")
		return c.Status(400).JSON(response)
	}
	var cashier Model.Cashier
	db.DB.Where("id = ?", cashierId).First(&cashier)

	//check if cashier exist
	if cashier.Id == 0 {
		response := Response.NewWebErrorResponse("cashier not found")
		return c.Status(404).JSON(response)
	}

	if cashier.Passcode != data["passcode"] {
		response := Response.NewWebErrorResponse("passcode not match")
		return c.Status(401).JSON(response)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Issuer":    strconv.Itoa(int(cashier.Id)),
		"ExpiresAt": time.Now().Add(time.Hour * 24).Unix(), //1 day
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		response := Response.NewWebErrorResponse(err.Error())
		return c.Status(401).JSON(response)
	}

	tokenResponse := dto.Token{
		Token: tokenString,
	}
	response := Response.NewWebResponse(tokenResponse)
	return c.Status(http.StatusOK).JSON(response)
}

// logout controller
// @Description logout Cashier
// @Summary logout cashier
// @Tags authentication
// @Produce json
// @Param cashierId path string true "cashier Id"
// @Param request body dto.Passcode true "request"
// @Success 200 {object} Response.WebResponse[string]
// @Router /cashier/{cashierId}/logout [post]
func Logout(c *fiber.Ctx) error {
	cashierId := c.Params("cashierId")
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	//check if passcode is empty
	if data["passcode"] == "" {
		response := Response.NewWebResponse("passcode cannot be empty")
		return c.Status(400).JSON(response)
	}

	var cashier Model.Cashier
	db.DB.Where("Id = ?", cashierId).First(&cashier)

	//check if cashier exist
	if cashier.Id == 0 {
		response := Response.NewWebResponse("Cashier Not found")
		return c.Status(404).JSON(response)
	}
	//check if passcode match
	if cashier.Passcode != data["passcode"] {
		response := Response.NewWebResponse("passcode not match")
		return c.Status(401).JSON(response)
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	response := Response.NewWebResponse("")
	return c.Status(200).JSON(response)
}

// passcode controller
// @Description passcode Cashier
// @Summary passcode cashier
// @Tags authentication
// @Produce json
// @Param cashierId path string true "cashier Id"
// @Param request body dto.Passcode true "request"
// @Success 200 {object} Response.WebResponse[dto.Passcode]
// @Router /cashier/{cashierId}/passcode [post]
func Passcode(c *fiber.Ctx) error {
	cashierId := c.Params("cashierId")
	var cashier Model.Cashier
	db.DB.Select("id,name,passcode").Where("id=?", cashierId).First(&cashier)

	if cashier.Name == "" || cashier.Id == 0 {
		response := Response.NewWebResponse("cashier not found")
		return c.Status(404).JSON(response)
	}

	passcode := dto.Passcode{
		Passcode: cashier.Passcode,
	}

	response := Response.NewWebResponse(passcode)
	return c.Status(200).JSON(response)
}
