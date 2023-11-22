package Controller

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"sales-api/Middleware"
	"sales-api/Model"
	db "sales-api/config"
	"strconv"
)

type Products struct {
	Products     Model.Product
	CategoriesId string `json:"categories_Id"`
}
type ProdDiscount struct {
	Id         int      `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku        string   `json:"sku"`
	Name       string   `json:"name"`
	Stock      int      `json:"stock"`
	Price      int      `json:"price"`
	Image      string   `json:"image"`
	CategoryId int      `json:"categoryId"`
	Discount   Discount `json:"discount"`
}
type Discount struct {
	Qty       int    `json:"qty"`
	Types     string `json:"type"`
	Result    int    `json:"result"`
	ExpiredAt int    `json:"expiredAt"`
}

func CreateProduct(c *fiber.Ctx) error {
	var data ProdDiscount
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

	Response := map[string]interface{}{
		"success": true,
		"message": "Success",
		"data":    product,
	}
	return (c.JSON(Response))

}

func GetProductDetails(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
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

	Response := map[string]interface{}{
		"success": true,
		"message": "Success",
		"data":    productsRes,
	}
	return (c.JSON(Response))
}

func UpdateProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")
	var product Model.Product

	db.DB.Find(&product, "id = ?", productId)

	if product.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
			"error":   map[string]interface{}{},
		})
	}

	var updateProductData Model.Product
	c.BodyParser(&updateProductData)

	if updateProductData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Product name is required",
			"error":   map[string]interface{}{},
		})
	}

	product.Name = updateProductData.Name
	product.CategoryId = updateProductData.CategoryId
	product.Image = updateProductData.Image
	product.Price = updateProductData.Price
	product.Stock = updateProductData.Stock

	db.DB.Save(&product)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    product,
	})
}

func ProductList(c *fiber.Ctx) error {
	//Token authenticate
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
	}
	if err := Middleware.AuthenticateToken(Middleware.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
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

		meta := map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": map[string]interface{}{
				"products": productsRes,
				"meta":     meta,
			},
		})
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

		meta := map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": map[string]interface{}{
				"products": productsRes,
				"meta":     meta,
			},
		})
	}
}

func DeleteProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")
	var product Model.Product

	db.DB.First(&product, productId)
	if product.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
			"error":   map[string]interface{}{},
		})
	}

	db.DB.Delete(&product)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
	})
}
