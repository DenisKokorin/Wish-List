package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	config "github.com/DenisKokorin/Wish-List/internal"
	clients "github.com/DenisKokorin/Wish-List/internal/clients/auth"
	groupapp "github.com/DenisKokorin/Wish-List/internal/group/app"
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

	log.Info("starting group application")

	authClient, err := clients.New(context.Background(), cfg.Clients.Auth.Address)
	if err != nil {
		log.Error("failed to init auth client: %w", err)
		os.Exit(1)
	}

	path := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	application := groupapp.New(log, cfg.GrpcGroup.Port, path, authClient)
	go application.GRPCserver.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Info("stopping group application", slog.String("sign", sign.String()))

	application.GRPCserver.Stop()

	log.Info("application group stoped")
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
