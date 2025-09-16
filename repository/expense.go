package repository

import (
	"context"
	"expense-tracker/entity"
	"time"
)

type ExpenseRepoItf interface {
	SelectCategories(ctx context.Context) ([]entity.GetCategoriesRes, error)
	InsertNewExpense(ctx context.Context, expense entity.AddExpense, userId int) (*int, error)
}

type ExpenseRepoStruct struct {
	db DBTX
}

func NewExpenseRepo(db DBTX) ExpenseRepoStruct {
	return ExpenseRepoStruct{
		db: db,
	}
}

func (er ExpenseRepoStruct) SelectCategories(ctx context.Context) ([]entity.GetCategoriesRes, error) {
	sql := `SELECT id, name
					FROM categories`
	db := ChooseDbOrTx(ctx, er.db)
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []entity.GetCategoriesRes
	for rows.Next() {
		var category entity.GetCategoriesRes
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (er ExpenseRepoStruct) InsertNewExpense(ctx context.Context, expense entity.AddExpense, userId int) (*int, error) {
	sql := `INSERT INTO expenses
					(user_id, title, category_id, amount, date, created_at, updated_at)
					VALUES
					($1, $2, $3, $4, $5, $6, $7)
					RETURNING (id);
	`
	var expenseId int
	db := ChooseDbOrTx(ctx, er.db)
	err := db.QueryRowContext(ctx, sql, &userId, &expense.Title, &expense.CategoryId, &expense.Amount, &expense.Date, time.Now(), time.Now()).Scan(&expenseId)
	if err != nil {
		return nil, err
	}
	return &expenseId, nil
}
