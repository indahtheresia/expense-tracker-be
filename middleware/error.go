package middleware

import (
	"errors"
	"expense-tracker/dto"
	"expense-tracker/middleware/logger"
	"expense-tracker/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors[0]

		logger.Log.Printf("Request ID: %s, Error: %s\n", GetRequestID(ctx), err.Error())

		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			validationErrors := make([]string, 0)
			var errMsg []string
			for i, fe := range ve {
				field := strings.ToLower(fe.Field())

				switch fe.Tag() {
				case "required":
					{
						errMsg = append(errMsg, fmt.Sprintf("%s is required", field))
					}
				case "min":
					{
						errMsg = append(errMsg, fmt.Sprintf("%s must be greater than or equal %s", field, fe.Param()))
					}
				case "max":
					{
						errMsg = append(errMsg, fmt.Sprintf("%s cannot be longer than %s characters.", field, fe.Param()))
					}
				case "lte":
					{
						errMsg = append(errMsg, fmt.Sprintf("%s should be less than %s", field, fe.Param()))
					}
				case "gte":
					{
						errMsg = append(errMsg, fmt.Sprintf("%s should be greater than %s", field, fe.Param()))
					}
				}
				validationErrors = append(validationErrors, errMsg[i])
			}

			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Error:   validationErrors,
				Data:    nil,
			})
			return
		}

		var customError dto.CustomError
		if errors.As(err, &customError) {
			util.ResponseMsg(ctx, false, customError.ErrorStr, nil, customError.Status)
			return
		}

		logger.Log.Error(err.Error())
		errMsg := fmt.Errorf("internal server error. please try again")
		util.ResponseMsg(ctx, false, errMsg.Error(), nil, http.StatusInternalServerError)
	}
}
