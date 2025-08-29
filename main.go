package main

import (
	"expense-tracker/app"
	"expense-tracker/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("error loading .env file", err)
	}
	db := db.InitDB()
	defer db.Close()

	app.App(db)
}
