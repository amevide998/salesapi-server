package dto

import (
	"sales-api/Model"
	"time"
)

type OrderResponse struct {
	Order    Model.Order                   `json:"order"`
	Products []*Model.ProductResponseOrder `json:"products"`
}

type SubTotalResponse struct {
	SubTotal int                           `json:"subTotal"`
	Products []*Model.ProductResponseOrder `json:"products"`
}

type CheckOrder struct {
	IsDownloaded bool `json:"is_downloaded"`
}

type OrderDetail struct {
	OrderId        int               `json:"order_id"`
	CashierId      int               `json:"cashier_id"`
	PaymentTypesId int               `json:"payment_types_id"`
	TotalPrice     int               `json:"total_price"`
	TotalPaid      int               `json:"total_paid"`
	TotalReturn    int               `json:"total_return"`
	ReceiptId      string            `json:"receipt_id"`
	CreatedAt      time.Time         `json:"created_at"`
	Cashier        Model.Cashier     `json:"cashier"`
	PaymentType    Model.PaymentType `json:"payment_type"`
}

type OrderDetailResponse struct {
	Orders   OrderDetail `json:"orders"`
	Products []*Model.Product
}

type OrderList struct {
	OrderId        int               `json:"orderId"`
	CashierID      int               `json:"cashiersId"`
	PaymentTypesId int               `json:"paymentTypesId"`
	TotalPrice     int               `json:"totalPrice"`
	TotalPaid      int               `json:"totalPaid"`
	TotalReturn    int               `json:"totalReturn"`
	ReceiptId      string            `json:"receiptId"`
	CreatedAt      time.Time         `json:"createdAt"`
	Payments       Model.PaymentType `json:"payment_type"`
	Cashiers       Model.Cashier     `json:"cashier"`
}

type OrderListResponse struct {
	Order []*OrderList `json:"order"`
	Meta  Pagination   `json:"meta"`
}
