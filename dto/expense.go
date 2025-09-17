package dto

type GetCategoriesRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AddExpense struct {
	Title      string `json:"title" binding:"required"`
	Amount     int    `json:"amount" binding:"required"`
	CategoryId int    `json:"category_id" binding:"required"`
	Date       string `json:"date"`
}

type UpdateExpense struct {
	Title      string `json:"title,omitempty"`
	Amount     int    `json:"amount,omitempty"`
	CategoryId int    `json:"category_id,omitempty"`
	Date       string `json:"date,omitempty"`
}
