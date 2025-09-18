package usecase

import (
	"context"
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/entity"
	"expense-tracker/repository"
)

type ExpenseUseCaseItf interface {
	GetCategories(ctx context.Context) ([]entity.GetCategoriesRes, error)
	InsertExpense(ctx context.Context, expense entity.AddExpense, userId int) (*int, error)
	UpdateExpense(ctx context.Context, expense entity.UpdateExpense, expenseId int) error
	DeleteExpense(ctx context.Context, addressId int) error
}

type ExpenseUseCaseStruct struct {
	er repository.ExpenseRepoItf
	tx repository.TxRepoItf
}

func NewExpenseUseCase(er repository.ExpenseRepoItf, tx repository.TxRepoItf) ExpenseUseCaseStruct {
	return ExpenseUseCaseStruct{
		er: er,
		tx: tx,
	}
}

func (euc ExpenseUseCaseStruct) GetCategories(ctx context.Context) ([]entity.GetCategoriesRes, error) {
	var categories []entity.GetCategoriesRes

	err := euc.tx.WithTx(ctx, func(ctx context.Context) error {
		allCategories, err := euc.er.SelectCategories(ctx)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorGetCategories.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}

		categories = allCategories
		return nil
	})

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (euc ExpenseUseCaseStruct) InsertExpense(ctx context.Context, expense entity.AddExpense, userId int) (*int, error) {
	var expenseId int

	err := euc.tx.WithTx(ctx, func(ctx context.Context) error {
		id, err := euc.er.InsertNewExpense(ctx, expense, userId)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorAddExpense.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}
		expenseId = *id
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &expenseId, nil
}

func (euc ExpenseUseCaseStruct) UpdateExpense(ctx context.Context, expense entity.UpdateExpense, expenseId int) error {
	err := euc.tx.WithTx(ctx, func(ctx context.Context) error {
		err := euc.er.UpdateExpense(ctx, expense, expenseId)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorUpdateExpense.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (euc ExpenseUseCaseStruct) DeleteExpense(ctx context.Context, addressId int) error {
	err := euc.tx.WithTx(ctx, func(ctx context.Context) error {
		notExists, err := euc.er.IsExpenseIdNotExists(ctx, addressId)
		if notExists {
			return dto.CustomError{
				ErrorStr:    constant.ErrorExpenseNotFound.Error(),
				InternalErr: constant.ErrorExpenseNotFound.Error(),
				Status:      constant.BadRequest,
			}
		}
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorDeleteExpense.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}
		err = euc.er.DeleteExpense(ctx, addressId)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorDeleteExpense.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
