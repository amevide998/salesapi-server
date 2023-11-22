package Controller

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"sales-api/Middleware"
	"sales-api/Model"
	"sales-api/Response"
	db "sales-api/config"
	"sales-api/dto"
	"strconv"
	"strings"
)

// revenue controller
// @Description get revenues
// @Summary get revenues
// @Tags report
// @Produce json
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[dto.Revenue]
// @Router /revenues [get]
func GetRevenues(c *fiber.Ctx) error {

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

	order := []Model.Order{}

	db.DB.Find(&order)
	TotalRevenues := make([]*Model.RevenueResponse, 0)

	Resp1 := Model.RevenueResponse{}
	Resp2 := Model.RevenueResponse{}
	Resp3 := Model.RevenueResponse{}

	sum1 := 0
	sum2 := 0
	sum3 := 0
	for _, v := range order {
		if v.PaymentTypesId == 1 {
			payment := Model.Payment{}
			paymentTypes := Model.PaymentType{}

			db.DB.Where("id=?", 1).First(&paymentTypes)
			db.DB.Where("payment_type_id=?", 1).First(&payment)

			sum1 += v.TotalPaid
			Resp1.Name = paymentTypes.Name
			Resp1.Logo = payment.Logo
			Resp1.TotalAmount = sum1
			Resp1.PaymentTypeId = v.PaymentTypesId
		}

		if v.PaymentTypesId == 2 {

			payment := Model.Payment{}
			paymentTypes := Model.PaymentType{}

			db.DB.Where("id=?", 2).First(&paymentTypes)
			db.DB.Where("payment_type_id=?", 2).First(&payment)

			sum2 += v.TotalPaid
			Resp2.Name = paymentTypes.Name
			Resp2.Logo = payment.Logo
			Resp2.TotalAmount = sum2
			Resp2.PaymentTypeId = v.PaymentTypesId
		}
		if v.PaymentTypesId == 3 {

			payment := Model.Payment{}
			paymentTypes := Model.PaymentType{}

			db.DB.Where("id=?", 3).First(&paymentTypes)
			db.DB.Where("payment_type_id=?", 3).First(&payment)

			sum3 += v.TotalPaid
			Resp3.Name = paymentTypes.Name
			Resp3.Logo = payment.Logo
			Resp3.TotalAmount = sum2
			Resp3.PaymentTypeId = v.PaymentTypesId
		}
	}
	TotalRevenues = append(TotalRevenues, &Resp1)
	TotalRevenues = append(TotalRevenues, &Resp2)
	TotalRevenues = append(TotalRevenues, &Resp3)

	dataResult := dto.Revenue{
		TotalRevenue: int64(sum1 + sum2 + sum3),
		PaymentTypes: TotalRevenues,
	}
	response := Response.NewWebResponse(dataResult)
	return c.Status(200).JSON(response)
}

type Sold struct {
	ProductId   string `json:"productId"`
	Quantities  string `json:"quantities"`
	TotalAmount int    `json:"totalAmount"`
}

// revenue controller
// @Description get solds
// @Summary get solds
// @Tags report
// @Produce json
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[dto.Sold]
// @Router /solds [get]
func GetSolds(c *fiber.Ctx) error {

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

	orders := []Model.Order{}
	db.DB.Find(&orders)

	TotalSold := make([]*Model.SoldResponse, 0)
	TotalSoldFinal := make([]*Model.SoldResponse, 0)

	for _, v := range orders {
		quantities := strings.Split(v.Quantities, ",")
		quantities = quantities[1:]

		products := strings.Split(v.ProductId, ",")
		products = products[1:]

		for i := 0; i < len(products); i++ {
			prods := Model.Product{}
			p, err := strconv.Atoi(products[i])
			q, errq := strconv.Atoi(quantities[i])

			if err != nil {
				log.Fatalf("->", err)
			}
			if errq != nil {
				log.Fatalf("->", errq)
			}

			db.DB.Where("id", p).Find(&prods)
			TotalSold = append(TotalSold, &Model.SoldResponse{
				Name:        prods.Name,
				ProductId:   p,
				TotalQty:    q,
				TotalAmount: q * prods.Price,
			})
		}

	}
	duplicates := []int{}
	for _, v := range TotalSold {

		if contains(duplicates, v.ProductId) == false {
			duplicates = append(duplicates, v.ProductId)
		}
	}
	quantityArray := []int{}
	for _, v := range duplicates {
		qty := 0
		for _, x := range TotalSold {
			if v == x.ProductId {
				qty = qty + x.TotalQty
			}
		}
		quantityArray = append(quantityArray, qty)

	}

	for i := 0; i < len(duplicates); i++ {

		prods := Model.Product{}
		db.DB.Where("id", duplicates[i]).Find(&prods)
		TotalSoldFinal = append(TotalSoldFinal, &Model.SoldResponse{
			Name:        prods.Name,
			TotalQty:    quantityArray[i],
			TotalAmount: quantityArray[i] * prods.Price,
			ProductId:   duplicates[i],
		})
	}

	dataResult := dto.Sold{
		TotalSold: TotalSoldFinal,
	}
	response := Response.NewWebResponse(dataResult)
	return c.Status(200).JSON(response)
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
