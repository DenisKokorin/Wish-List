package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	config "github.com/DenisKokorin/Wish-List/internal"
	"github.com/DenisKokorin/Wish-List/internal/wishlist/app"
	"github.com/joho/godotenv"
)

func main() {
	cfg := config.MustLoad()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := "postgres"
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	log := setupLogger()

	log.Info("starting wishlist application")

	path := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	application := app.New(log, cfg.GrpcWishList.Port, path)

	go application.GRPCserver.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Info("stopping wishlist application", slog.String("sign", sign.String()))

	application.GRPCserver.Stop()

	log.Info("application wishlist stoped")
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
