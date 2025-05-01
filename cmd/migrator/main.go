package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := "postgres"
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	path := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println(path)
	db, err := sql.Open("postgres", path)
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	log.Info("migrations applied")

	db.Close()
}
