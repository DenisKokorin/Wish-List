package main

import (
	"database/sql"
	"fmt"

	config "github.com/DenisKokorin/Wish-List/internal"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.MustLoad()

	path := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)
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
