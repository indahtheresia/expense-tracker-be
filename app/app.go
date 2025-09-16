package app

import (
	"database/sql"
	"expense-tracker/constant"
	"expense-tracker/handler"
	"expense-tracker/middleware"
	"expense-tracker/middleware/logger"
	"expense-tracker/repository"
	"expense-tracker/usecase"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func App(db *sql.DB) {
	logger.SetLogrusLogger()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorMiddleware())

	transactionRepo := repository.NewTx(db)

	userRepo := repository.NewUserRepo(db)
	userUseCase := usecase.NewUserUseCase(userRepo, transactionRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	expenseRepo := repository.NewExpenseRepo(db)
	expenseUseCase := usecase.NewExpenseUseCase(expenseRepo, transactionRepo)
	expenseHandler := handler.NewExpenseHandler(expenseUseCase)

	r.POST("/users/register", userHandler.Register)
	r.POST("/users/login", userHandler.Login)

	r.POST("/users/expense", middleware.AuthenticationMiddleware, expenseHandler.InsertExpense)

	r.GET("/categories", expenseHandler.GetCategories)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(constant.ServerPort); err != nil {
		log.Fatalf("Failed starting server: %v\n", err)
	}
}
