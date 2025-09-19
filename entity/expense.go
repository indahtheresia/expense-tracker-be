package entity

import "time"

type GetCategoriesRes struct {
	Id   int
	Name string
}

type AddExpense struct {
	Title      string
	Amount     float64
	CategoryId int
	Date       time.Time
}

type UpdateExpense struct {
	Title      string
	Amount     float64
	CategoryId int
	Date       time.Time
}

type GetExpenseRes struct {
	Id           int
	Title        string
	Amount       float64
	CategoryId   int
	CategoryName string
	Date         time.Time
}
