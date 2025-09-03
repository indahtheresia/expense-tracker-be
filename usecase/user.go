package usecase

import (
	"context"
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/entity"
	"expense-tracker/repository"
	"expense-tracker/util"
)

type UserUseCaseItf interface {
	InsertUser(ctx context.Context, user entity.InsertUserReq) (*int, error)
}

type UserUseCaseStruct struct {
	ur repository.UserRepoItf
	tx repository.TxRepoItf
}

func NewUserUseCase(ur repository.UserRepoItf, tx repository.TxRepoItf) UserUseCaseStruct {
	return UserUseCaseStruct{
		ur: ur,
		tx: tx,
	}
}

func (uuc UserUseCaseStruct) InsertUser(ctx context.Context, user entity.InsertUserReq) (*int, error) {
	password, err := util.GenerateBcrypt(user.Password)
	if err != nil {
		return nil, dto.CustomError{
			ErrorStr:    err.Error(),
			InternalErr: err.Error(),
			Status:      constant.InternalServerError,
		}
	}

	user.Password = *password

	var userId *int

	err = uuc.tx.WithTx(ctx, func(ctx context.Context) error {
		emailExists := uuc.ur.IsEmailExists(ctx, user.Email)
		if emailExists {
			return dto.CustomError{
				ErrorStr:    constant.ErrorUserEmailExists.Error(),
				InternalErr: constant.ErrorUserEmailExists.Error(),
				Status:      constant.BadRequest,
			}
		}

		userId, err = uuc.ur.InsertUser(ctx, user)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorInsertNewUser.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return userId, nil
}
