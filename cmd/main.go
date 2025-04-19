package main

import (
	"log/slog"
	"os"

	config "github.com/DenisKokorin/Wish-List/internal"
	"github.com/DenisKokorin/Wish-List/internal/app"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("starting application")

	application := app.New(log, cfg.GRPC.Port)

	application.GRPCserver.MustRun()

	log.Info("end")
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
