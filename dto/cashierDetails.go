package dto

import "time"

type CashierDetails struct {
	CashierId int       `json:"cashier_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
