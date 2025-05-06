package app

import (
	"log/slog"

	"github.com/DenisKokorin/Wish-List/internal/storage/postgres"
	grpcapp "github.com/DenisKokorin/Wish-List/internal/wishlist/app/gRPC"
	wishlistservice "github.com/DenisKokorin/Wish-List/internal/wishlist/services"
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
