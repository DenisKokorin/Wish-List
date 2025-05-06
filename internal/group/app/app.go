package groupapp

import (
	"log/slog"

	clients "github.com/DenisKokorin/Wish-List/internal/clients/auth"
	groupgrpcapp "github.com/DenisKokorin/Wish-List/internal/group/app/gRPC"
	groupservice "github.com/DenisKokorin/Wish-List/internal/group/services"
	"github.com/DenisKokorin/Wish-List/internal/storage/postgres"
)

type App struct {
	GRPCserver *groupgrpcapp.App
}

func New(log *slog.Logger, port int, path string, authClient *clients.CLient) *App {
	storage, err := postgres.New(path)
	if err != nil {
		panic(err)
	}

	groupService := groupservice.New(log, storage, authClient)

	grpcApp := groupgrpcapp.New(log, port, groupService)
	return &App{
		GRPCserver: grpcApp,
	}
}
