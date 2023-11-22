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

// Product controller
// @Description create product
// @Summary create product
// @Tags Product
// @Produce json
// @Param request body dto.ProdDiscount true "request"
// @Param Authorization header string true "authorization"
// @Success 201 {object} Response.WebResponse[Model.Product]
// @Router /products [post]
func CreateProduct(c *fiber.Ctx) error {
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

	var data dto.ProdDiscount
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Product error in post request %v", err)
	}

	var p []Model.Product
	db.DB.Find(&p)

	discount := Model.Discount{
		Qty:       data.Discount.Qty,
		Type:      data.Discount.Types,
		Result:    data.Discount.Result,
		ExpiredAt: data.Discount.ExpiredAt,
	}
	db.DB.Create(&discount)

	product := Model.Product{
		Name:       data.Name,
		Image:      data.Image,
		CategoryId: data.CategoryId,
		DiscountId: discount.Id,
		Price:      data.Price,
		Stock:      data.Stock,
	}
	db.DB.Create(&product)

	db.DB.Table("products").Where("id = ?", product.Id).Update("sku", "SK00"+strconv.Itoa(product.Id))

	response := Response.NewWebResponse(product)
	return c.Status(http.StatusCreated).JSON(response)

}

// Product controller
// @Description get product details
// @Summary get product details
// @Tags Product
// @Produce json
// @Param productId path string true "productId"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[Model.ProductResult]
// @Router /products/{productId} [get]
func GetProductDetails(c *fiber.Ctx) error {
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

	productId := c.Params("productId")
	productsRes := make([]*Model.ProductResult, 0)
	var products []Model.Product
	db.DB.Where("id = ? ", productId).Find(&products)

	var category Model.Category
	var discount Model.Discount
	for i := 0; i < len(products); i++ {
		db.DB.Where("id = ?", products[i].CategoryId).Find(&category)

		db.DB.Where("id = ?", products[i].DiscountId).Find(&discount)

		//productsRes =
		productsRes = append(productsRes,
			&Model.ProductResult{
				Id:       products[i].Id,
				Sku:      products[i].Sku,
				Name:     products[i].Name,
				Stock:    products[i].Stock,
				Price:    products[i].Price,
				Image:    products[i].Image,
				Category: category,
				Discount: discount,
			},
		)
	}

	response := Response.NewWebResponse(productsRes)
	return c.Status(http.StatusOK).JSON(response)
}

// Product controller
// @Description update products
// @Summary update products
// @Tags Product
// @Produce json
// @Param productId path string true "productId"
// @Param request body Model.Product true "request"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[Model.Product]
// @Router /products/{productId} [put]
func UpdateProduct(c *fiber.Ctx) error {
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

	productId := c.Params("productId")
	var product Model.Product

	db.DB.Find(&product, "id = ?", productId)

	if product.Name == "" {
		response := Response.NewWebErrorResponse("product Not Found")
		return c.Status(404).JSON(response)
	}

	var updateProductData Model.Product
	err := c.BodyParser(&updateProductData)
	if err != nil {
		return err
	}

	if updateProductData.Name == "" {
		response := Response.NewWebErrorResponse("Product name is required")
		return c.Status(400).JSON(response)
	}

	product.Name = updateProductData.Name
	product.CategoryId = updateProductData.CategoryId
	product.Image = updateProductData.Image
	product.Price = updateProductData.Price
	product.Stock = updateProductData.Stock

	db.DB.Save(&product)
	response := Response.NewWebResponse(product)
	return c.Status(200).JSON(response)
}

// Product controller
// @Description product list
// @Summary product list
// @Tags Product
// @Produce json
// @Param limit query string true "limit"
// @Param skip query string true "skip"
// @Param categoryId query string true "categoryId"
// @Param q query string true "q"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[dto.ProductList]
// @Router /products [get]
func ProductList(c *fiber.Ctx) error {
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

	limit := c.Query("limit")
	skip := c.Query("skip")
	categoryId := c.Query("categoryId")
	productName := c.Query("q")
	intLimit, _ := strconv.Atoi(limit)
	intSkip, _ := strconv.Atoi(skip)
	var products []Model.Product

	productsRes := make([]*Model.ProductResult, 0)

	if productName == "" {
		var count int64
		db.DB.Where("category_id = ?", categoryId).Limit(intLimit).Offset(intSkip).Find(&products).Count(&count)

		var category Model.Category
		var discount Model.Discount
		for i := 0; i < len(products); i++ {

			db.DB.Table("categories").Where("id = ?", products[i].CategoryId).Find(&category)

			db.DB.Where("id = ?", products[i].DiscountId).Limit(intLimit).Offset(intSkip).Find(&discount).Count(&count)
			count = int64(len(products))
			//productsRes =
			productsRes = append(productsRes,
				&Model.ProductResult{
					Id:       products[i].Id,
					Sku:      products[i].Sku,
					Name:     products[i].Name,
					Stock:    products[i].Stock,
					Price:    products[i].Price,
					Image:    products[i].Image,
					Category: category,
					Discount: discount,
				},
			)
		}

		meta := dto.Pagination{
			Total: count,
			Limit: intLimit,
			Skip:  intSkip,
		}

		dataResult := dto.ProductList{
			ProductRes: productsRes,
			Meta:       meta,
		}
		response := Response.NewWebResponse(dataResult)
		return c.Status(200).JSON(response)
	} else {

		var count int64
		if categoryId != "" {
			db.DB.Where("category_id = ? AND name= ?", categoryId, productName).Limit(intLimit).Offset(intSkip).Find(&products).Count(&count)
		} else {
			db.DB.Where(" name= ?", productName).Limit(intLimit).Offset(intSkip).Find(&products).Count(&count)
		}
		var category Model.Category
		var discount Model.Discount
		for i := 0; i < len(products); i++ {
			db.DB.Where("id = ?", products[i].CategoryId).Find(&category)
			db.DB.Where("id = ?", products[i].DiscountId).Limit(intLimit).Offset(intSkip).Find(&discount).Count(&count)
			count = int64(len(products))
			productsRes = append(productsRes,
				&Model.ProductResult{
					Id:       products[i].Id,
					Sku:      products[i].Sku,
					Name:     products[i].Name,
					Stock:    products[i].Stock,
					Price:    products[i].Price,
					Image:    products[i].Image,
					Category: category,
					Discount: discount,
				},
			)
		}

		meta := dto.Pagination{
			Total: count,
			Limit: intLimit,
			Skip:  intSkip,
		}

		dataResult := dto.ProductList{
			ProductRes: productsRes,
			Meta:       meta,
		}

		response := Response.NewWebResponse(dataResult)
		return c.Status(200).JSON(response)
	}
}

// Product controller
// @Description delete product
// @Summary delete product
// @Tags Product
// @Produce json
// @Param productId path string true "product id"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[string]
// @Router /products/{productId} [delete]
func DeleteProduct(c *fiber.Ctx) error {

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

	productId := c.Params("productId")
	var product Model.Product

	db.DB.First(&product, productId)
	if product.Id == 0 {
		response := Response.NewWebErrorResponse("Product Not Found")
		return c.Status(404).JSON(response)
	}

	db.DB.Delete(&product)

	response := Response.NewWebResponse("")
	return c.Status(200).JSON(response)
}
