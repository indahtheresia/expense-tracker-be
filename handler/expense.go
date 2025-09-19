package handler

import (
	"expense-tracker/constant"
	"expense-tracker/dto"
	"expense-tracker/entity"
	"expense-tracker/usecase"
	"expense-tracker/util"
	"fmt"
	"strconv"
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

func (eh ExpenseHandler) UpdateExpense(ctx *gin.Context) {
	var editExpense dto.UpdateExpense
	err := ctx.ShouldBindJSON(&editExpense)
	if err != nil {
		ctx.Error(err)
		return
	}

	expenseId := ctx.Param("id")
	id, err := strconv.Atoi(expenseId)
	if err != nil {
		ctx.Error(err)
		return
	}

	date := "0001-01-01"
	if editExpense.Date != "" {
		date = editExpense.Date
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		ctx.Error(dto.CustomError{
			ErrorStr:    constant.ErrorParsingDate.Error(),
			InternalErr: constant.ErrorParsingDate.Error(),
			Status:      constant.BadRequest,
		})
		return
	}

	data := entity.UpdateExpense{
		Title:      editExpense.Title,
		CategoryId: editExpense.CategoryId,
		Amount:     editExpense.Amount,
		Date:       parsedDate,
	}

	err = eh.euc.UpdateExpense(ctx, data, id)
	if err != nil {
		ctx.Error(err)
		return
	}

	util.ResponseMsg(ctx, true, nil, "Success update expense data", constant.OK)
}

func (eh ExpenseHandler) DeleteExpense(ctx *gin.Context) {
	expenseId := ctx.Param("id")
	id, err := strconv.Atoi(expenseId)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = eh.euc.DeleteExpense(ctx, id)
	if err != nil {
		ctx.Error(err)
		return
	}

	respMsg := fmt.Sprintf("Expense with id %d has been successfully deleted", id)

	util.ResponseMsg(ctx, true, nil, respMsg, constant.NoContent)
}

func (eh ExpenseHandler) GetExpenses(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	userId := sub.(entity.SubAuth).Id
	allExpenses, err := eh.euc.GetExpenses(ctx, userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	resData := []dto.GetExpenseRes{}
	for _, val := range allExpenses {
		resData = append(resData, dto.GetExpenseRes{
			Id:           val.Id,
			Title:        val.Title,
			Amount:       val.Amount,
			CategoryId:   val.CategoryId,
			CategoryName: val.CategoryName,
			Date:         val.Date,
		})
	}

	util.ResponseMsg(ctx, true, nil, resData, constant.OK)
}
