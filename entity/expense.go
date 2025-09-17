package entity

import "time"

type GetCategoriesRes struct {
	Id   int
	Name string
}

type AddExpense struct {
	Title      string
	Amount     int
	CategoryId int
	Date       time.Time
}

type UpdateExpense struct {
	Title      string
	Amount     int
	CategoryId int
	Date       time.Time
}
