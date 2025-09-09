package usecase

import (
	"context"
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/entity"
	"expense-tracker/repository"
	"expense-tracker/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserUseCaseItf interface {
	InsertUser(ctx context.Context, user entity.InsertUserReq) (*int, error)
	LoginUser(ctx context.Context, user entity.LoginReq) (*string, *string, error)
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
		emailExists, err := uuc.ur.IsEmailExists(ctx, user.Email)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorInternalServer.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}
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

func (uuc UserUseCaseStruct) LoginUser(ctx context.Context, user entity.LoginReq) (*string, *string, error) {
	var token string
	var role string
	err := uuc.tx.WithTx(ctx, func(ctx context.Context) error {
		emailNotExists, err := uuc.ur.IsEmailNotExists(ctx, user.Email)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorInternalServer.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}
		if emailNotExists {
			return dto.CustomError{
				ErrorStr:    constant.ErrorUserEmailNotExists.Error(),
				InternalErr: constant.ErrorUserEmailNotExists.Error(),
				Status:      constant.BadRequest,
			}
		}

		hashPassword, err := uuc.ur.SelectHashPasswordByEmail(ctx, user.Email)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorUserEmailNotExists.Error(),
				InternalErr: err.Error(),
				Status:      constant.BadRequest,
			}
		}
		hashPass := []byte(*hashPassword)
		err = util.CompareHashPassword(hashPass, []byte(user.Password))
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorInternalServer.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}

		userId, err := uuc.ur.SelectIdByEmail(ctx, user.Email)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorInternalServer.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}

		id := strconv.Itoa(*userId)
		now := time.Now()
		registeredClaims := dto.CustomClaims{
			Role:       constant.UserRole,
			Permission: []string{"read", "edit"},
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:  "expense_tracker_app",
				Subject: id,
				IssuedAt: &jwt.NumericDate{
					Time: now,
				},
				ExpiresAt: &jwt.NumericDate{
					Time: now.Add(24 * time.Hour),
				},
			},
		}

		jwtToken, err := util.GenerateJWTToken(&registeredClaims)
		if err != nil {
			return dto.CustomError{
				ErrorStr:    constant.ErrorInternalServer.Error(),
				InternalErr: err.Error(),
				Status:      constant.InternalServerError,
			}
		}
		token = *jwtToken
		role = constant.UserRole
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return &token, &role, nil
}
