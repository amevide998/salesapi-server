package dto

// Payment struct with two values
type Payment struct {
	Id            uint   `json:"paymentId"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	PaymentTypeId int    `json:"payment_type_id"`
	Logo          string `json:"logo"`
}

type CreatePayment struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Logo string `json:"logo"`
}

type ListPayment struct {
	Payments []Payment  `json:"payment"`
	Meta     Pagination `json:"meta"`
}
