package dto

type Pagination struct {
	Total int64 `json:"total"`
	Limit int   `json:"limit"`
	Skip  int   `json:"skip"`
}
