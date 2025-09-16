package handler

import (
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/entity"
	"expense-tracker/usecase"
	"expense-tracker/util"
	"fmt"
	"time"

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

func (eh ExpenseHandler) InsertExpense(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	userId := sub.(entity.SubAuth).Id

	var expense dto.AddExpense
	err := ctx.ShouldBindJSON(&expense)
	if err != nil {
		ctx.Error(err)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", expense.Date)
	if err != nil {
		ctx.Error(dto.CustomError{
			ErrorStr:    constant.ErrorParsingDate.Error(),
			InternalErr: constant.ErrorParsingDate.Error(),
			Status:      constant.BadRequest,
		})
		return
	}

	data := entity.AddExpense{
		Title:      expense.Title,
		CategoryId: expense.CategoryId,
		Amount:     expense.Amount,
		Date:       parsedDate,
	}

	expenseId, err := eh.euc.InsertExpense(ctx, data, userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	respMsg := fmt.Sprintf("success create a new expense with id %d", *expenseId)

	util.ResponseMsg(ctx, true, nil, respMsg, constant.Created)
}
