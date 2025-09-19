package dto

import "time"

type GetCategoriesRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AddExpense struct {
	Title      string  `json:"title" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	CategoryId int     `json:"category_id" binding:"required"`
	Date       string  `json:"date"`
}

type UpdateExpense struct {
	Title      string  `json:"title,omitempty"`
	Amount     float64 `json:"amount,omitempty"`
	CategoryId int     `json:"category_id,omitempty"`
	Date       string  `json:"date,omitempty"`
}

type GetExpenseRes struct {
	Id           int       `json:"expense_id"`
	Title        string    `json:"title"`
	Amount       float64   `json:"amount"`
	CategoryId   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Date         time.Time `json:"date"`
}
