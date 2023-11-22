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
)

// category controller
// @Description create category
// @Summary create category
// @Tags category
// @Produce json
// @Param request body Model.Category true "request"
// @Success 200 {object} Response.WebResponse[Model.Category]
// @Router /categories [post]
func CreateCategory(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("category error in post request %v", err)
	}
	if data["name"] == "" {
		response := Response.NewWebErrorResponse("category name is required")
		return c.Status(400).JSON(response)
	}
	category := Model.Category{
		Name: data["name"],
	}
	db.DB.Create(&category)

	response := Response.NewWebResponse(category)
	return c.Status(200).JSON(response)
}

// category controller
// @Description get category details
// @Summary get category details
// @Tags category
// @Produce json
// @Param categoryId path string true "category id"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[Model.Category]
// @Router /categories/{categoryId} [get]
func GetCategoryDetails(c *fiber.Ctx) error {
	categoryId := c.Params("categoryId")

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

	var category Model.Category
	db.DB.Select("id ,name").Where("id=?", categoryId).First(&category)

	if category.Name == "" {
		response := Response.NewWebErrorResponse("No category found")
		return c.Status(404).JSON(response)
	}

	categoryData := dto.Category{
		Id:   uint(category.Id),
		Name: category.Name,
	}

	response := Response.NewWebResponse(categoryData)
	return c.JSON(response)
}

// category controller
// @Description delete category
// @Summary get delete category
// @Tags category
// @Produce json
// @Param categoryId path string true "category id"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[string]
// @Router /categories/{categoryId} [delete]
func DeleteCategory(c *fiber.Ctx) error {

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

	categoryId := c.Params("categoryId")
	var category Model.Category
	db.DB.Where("id=?", categoryId).First(&category)

	if category.Id == 0 {
		response := Response.NewWebErrorResponse("category not found")
		return c.Status(404).JSON(response)
	}

	db.DB.Where("id = ?", categoryId).Delete(&category)

	response := Response.NewWebResponse(categoryId)
	return c.Status(200).JSON(response)
}

// category controller
// @Description update category
// @Summary get update category
// @Tags category
// @Produce json
// @Param categoryId path string true "category id"
// @Param request body Model.Category true "request"
// @Param Authorization header string true "authorization"
// @Success 200 {object} Response.WebResponse[Model.Category]
// @Router /categories/{categoryId} [put]
func UpdateCategory(c *fiber.Ctx) error {
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

	categoryId := c.Params("categoryId")
	var category Model.Category

	db.DB.Find(&category, "id = ?", categoryId)

	if category.Name == "" {
		response := Response.NewWebErrorResponse("Category not exist against this id")
		return c.Status(404).JSON(response)
	}

	var updateCashierData Model.Category
	c.BodyParser(&updateCashierData)
	if updateCashierData.Name == "" {
		response := Response.NewWebErrorResponse("Category name is required")
		return c.Status(400).JSON(response)
	}
	category.Name = updateCashierData.Name
	db.DB.Save(&category)

	response := Response.NewWebResponse(category)
	return c.Status(200).JSON(response)

}

// category controller
// @Description get category list
// @Summary get category list
// @Tags category
// @Produce json
// @Param Authorization header string true "authorization"
// @Param limit query string true "limit"
// @Param skip query string true "skip"
// @Success 200 {object} Response.WebResponse[dto.CategoryList]
// @Router /categories [get]
func CategoryList(c *fiber.Ctx) error {
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

	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var category []dto.Category
	db.DB.Select("id ,name").Limit(limit).Offset(skip).Find(&category).Count(&count)
	metaMap := dto.Pagination{
		Total: count,
		Limit: limit,
		Skip:  skip,
	}

	categoriesData := dto.CategoryList{
		Categories: category,
		Meta:       metaMap,
	}

	response := Response.NewWebResponse(categoriesData)
	return c.JSON(response)

}
