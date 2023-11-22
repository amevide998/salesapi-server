package dto

// Category struct with two values
type Category struct {
	Id   uint   `json:"categoryId"`
	Name string `json:"name"`
}

type CategoryList struct {
	Categories []Category `json:"categories"`
	Meta       Pagination `json:"page"`
}
