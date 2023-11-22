package dto

import "sales-api/Model"

type Revenue struct {
	TotalRevenue int64 `json:"total_revenue"`
	PaymentTypes []*Model.RevenueResponse
}

type Sold struct {
	TotalSold []*Model.SoldResponse `json:"total_sold"`
}
