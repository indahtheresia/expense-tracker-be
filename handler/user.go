package handler

import (
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/entity"
	"expense-tracker/usecase"
	"expense-tracker/util"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	uuc usecase.UserUseCaseItf
}

func NewUserHandler(uuc usecase.UserUseCaseItf) UserHandler {
	return UserHandler{
		uuc: uuc,
	}
}

func (uh UserHandler) Register(ctx *gin.Context) {
	var user dto.RegisterUserReq
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.Error(err)
		return
	}

	data := entity.InsertUserReq{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	id, err := uh.uuc.InsertUser(ctx, data)
	if err != nil {
		ctx.Error(err)
		return
	}

	userId := *id
	resMsg := fmt.Sprintf("Registered! User Id: %d.", userId)

	util.ResponseMsg(ctx, true, nil, resMsg, constant.Created)
}

func (uh UserHandler) Login(ctx *gin.Context) {
	var user dto.LoginUserReq
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.Error(err)
		return
	}

	data := entity.LoginReq{
		Email:    user.Email,
		Password: user.Password,
	}

	token, role, err := uh.uuc.LoginUser(ctx, data)
	if err != nil {
		ctx.Error(err)
		return
	}

	resData := dto.LoginUserRes{
		Token: *token,
		Role:  *role,
	}

	util.ResponseMsg(ctx, true, nil, resData, constant.OK)
}
