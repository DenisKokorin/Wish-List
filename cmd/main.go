package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	config "github.com/DenisKokorin/Wish-List/internal"
	"github.com/DenisKokorin/Wish-List/internal/app"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("starting application")

	application := app.New(log, cfg.GRPC.Port)

	go application.GRPCserver.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Info("stopping application", slog.String("sign", sign.String()))

	application.GRPCserver.Stop()

	log.Info("application stoped")
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
