package handler

import (
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/usecase"
	"expense-tracker/util"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	euc usecase.ExpenseUseCaseItf
}

func NewExpenseHandler(euc usecase.ExpenseUseCaseItf) ExpenseHandler {
	return ExpenseHandler{
		euc: euc,
	}
}

func (eh ExpenseHandler) GetCategories(ctx *gin.Context) {
	categories, err := eh.euc.GetCategories(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	resData := []dto.GetCategoriesRes{}
	for _, val := range categories {
		resData = append(resData, dto.GetCategoriesRes{
			Id:   val.Id,
			Name: val.Name,
		})
	}

	util.ResponseMsg(ctx, true, nil, resData, constant.OK)
}
