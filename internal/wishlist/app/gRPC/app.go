package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	wishlist "github.com/DenisKokorin/Wish-List/internal/wishlist/gRPC"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCserver *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, wishList wishlist.WishList) *App {
	gRPCserver := grpc.NewServer()
	wishlist.Register(gRPCserver, wishList)

	return &App{
		log:        log,
		gRPCserver: gRPCserver,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	log := a.log.With(slog.String("op", "grpcWLapp.Run"), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", "grpcaWLpp.Run", err)
	}

	log.Info("WLgrpc server is running", slog.String("addr", l.Addr().String()))
	if err := a.gRPCserver.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", "grpcWLapp.Run", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping WLgrpc server")
	a.gRPCserver.GracefulStop()
}
