package repository

import (
	"context"
	"expense-tracker/entity"
	"fmt"
	"reflect"
	"time"
)

type ExpenseRepoItf interface {
	SelectCategories(ctx context.Context) ([]entity.GetCategoriesRes, error)
	InsertNewExpense(ctx context.Context, expense entity.AddExpense, userId int) (*int, error)
	UpdateExpense(ctx context.Context, expense entity.UpdateExpense, expenseId int) error
	DeleteExpense(ctx context.Context, addressId int) error
	IsExpenseIdNotExists(ctx context.Context, addressId int) (bool, error)
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

func (er ExpenseRepoStruct) UpdateExpense(ctx context.Context, expense entity.UpdateExpense, expenseId int) error {
	var fieldAssignments string

	allowedEdits := map[string]string{
		"Title":      "title",
		"Amount":     "amount",
		"CategoryId": "category_id",
		"Date":       "date",
	}

	val := reflect.ValueOf(expense)
	typ := reflect.TypeOf(expense)
	first := true
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if !field.CanInterface() {
			continue
		}

		column, exists := allowedEdits[fieldType.Name]
		if !exists {
			continue
		}

		editVal := field.Interface()
		switch v := editVal.(type) {
		case string:
			if v != "" {
				if !first {
					fieldAssignments += ", "
				}
				fieldAssignments += fmt.Sprintf("%s = '%s'", column, v)
				first = false
			}
		case int:
			if v != 0 {
				if !first {
					fieldAssignments += ", "
				}
				fieldAssignments += fmt.Sprintf("%s = %d", column, v)
				first = false
			}
		case time.Time:
			if !v.IsZero() {
				if !first {
					fieldAssignments += ", "
				}
				fieldAssignments += fmt.Sprintf("%s = '%v'", column, v.Format("2006-01-02"))
				first = false
			}
		}
	}

	fieldAssignments += ", updated_at = NOW()"
	sql := fmt.Sprintf("UPDATE expenses SET %s WHERE id = %d", fieldAssignments, expenseId)

	db := ChooseDbOrTx(ctx, er.db)

	_, err := db.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

func (er ExpenseRepoStruct) DeleteExpense(ctx context.Context, addressId int) error {
	sql := `UPDATE expenses
					SET deleted_at = $1
					WHERE id = $2
	`

	db := ChooseDbOrTx(ctx, er.db)
	_, err := db.ExecContext(ctx, sql, time.Now(), addressId)
	if err != nil {
		return err
	}

	return nil
}

func (er ExpenseRepoStruct) IsExpenseIdNotExists(ctx context.Context, addressId int) (bool, error) {
	var idExists bool
	sql := `SELECT NOT EXISTS(SELECT 1 FROM expenses WHERE id = $1 AND deleted_at IS NULL)`
	db := ChooseDbOrTx(ctx, er.db)
	err := db.QueryRowContext(ctx, sql, addressId).Scan(&idExists)
	if err != nil {
		return true, err
	}
	return idExists, nil
}
