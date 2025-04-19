package main

import (
	"log/slog"
	"os"
	config "wishlistservice/internal"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("starting application", slog.Any("cfg", cfg))

}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
