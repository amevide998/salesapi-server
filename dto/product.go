package dto

import "sales-api/Model"

type Products struct {
	ProductId int `json:"productId"`
	Quantity  int `json:"qty"`
}

type ProductsCategory struct {
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

type ProductList struct {
	ProductRes []*Model.ProductResult `json:"product_res"`
	Meta       Pagination             `json:"meta"`
}
