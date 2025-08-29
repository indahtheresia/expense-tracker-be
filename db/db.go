package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB() *sql.DB {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, name, password, dbname)
	db, err := sql.Open("pgx", connString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
