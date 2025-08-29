package app

import (
	"database/sql"
	"expense-tracker/constant"
	"log"

	"github.com/gin-gonic/gin"
)

func App(db *sql.DB) {
	r := gin.New()
	r.Use(gin.Recovery())

	if err := r.Run(constant.ServerPort); err != nil {
		log.Fatalf("Failed starting server: %v\n", err)
	}
}
