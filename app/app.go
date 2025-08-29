package app

import (
	"database/sql"
	"expense-tracker/constant"
	"expense-tracker/middleware"
	"expense-tracker/middleware/logger"
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

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(constant.ServerPort); err != nil {
		log.Fatalf("Failed starting server: %v\n", err)
	}
}
