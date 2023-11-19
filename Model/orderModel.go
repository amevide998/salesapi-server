package Model

import "time"

type Order struct {
	Id             int       `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	CashierID      int       `json:"cashier_id"`
	PaymentTypesId int       `json:"payment_types_id"`
	TotalPrice     int       `json:"total_price"`
	TotalPaid      int       `json:"total_paid"`
	TotalReturn    int       `json:"total_return"`
	ReceiptId      string    `json:"receipt_id"`
	IsDownload     int       `json:"is_download"`
	ProductId      string    `json:"product_id"`
	Quantities     string    `json:"quantities"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProductResponseOrder struct {
	ProductId        int      `json:"productId" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name             string   `json:"name"`
	Price            int      `json:"price"`
	Qty              int      `json:"qty"`
	Discount         Discount `json:"discount"`
	TotalNormalPrice int      `json:"total_normal_price"`
	TotalFinalPrice  int      `json:"total_final_price"`
}

type ProductOrder struct {
	Id         int    `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku        string `json:"sku"`
	Name       string `json:"name"`
	Stock      int    `json:"stock"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryId int    `json:"category_id"`
	DiscountId int    `json:"discount_id"`
}

type RevenueResponse struct {
	PaymentTypeId int    `json:"payment_type_id"`
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	TotalAmount   int    `json:"total_amount"`
}

type SoldResponse struct {
	ProductId   int    `json:"product_id"`
	Name        string `json:"name"`
	TotalQty    int    `json:"total_qty"`
	TotalAmount int    `json:"total_amount"`
}
