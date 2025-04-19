package app

import (
	"log/slog"

	grpcapp "github.com/DenisKokorin/Wish-List/internal/app/gRPC"
)

type App struct {
	GRPCserver *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	grpcApp := grpcapp.New(log, port)
	return &App{
		GRPCserver: grpcApp,
	}
}
