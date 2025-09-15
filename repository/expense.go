package repository

import (
	"context"
	"expense-tracker/entity"
)

type ExpenseRepoItf interface {
	SelectCategories(ctx context.Context) ([]entity.GetCategoriesRes, error)
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
