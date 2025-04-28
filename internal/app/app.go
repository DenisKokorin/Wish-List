package app

import (
	"log/slog"

	grpcapp "github.com/DenisKokorin/Wish-List/internal/app/gRPC"
	wishlistservice "github.com/DenisKokorin/Wish-List/internal/services/wishlist"
	"github.com/DenisKokorin/Wish-List/internal/storage/postgres"
)

type App struct {
	GRPCserver *grpcapp.App
}

func New(log *slog.Logger, port int, path string) *App {
	storage, err := postgres.New(path)
	if err != nil {
		panic(err)
	}

	wishListService := wishlistservice.New(log, storage, storage)

	grpcApp := grpcapp.New(log, port, wishListService)
	return &App{
		GRPCserver: grpcApp,
	}
}
