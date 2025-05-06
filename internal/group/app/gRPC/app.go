package groupgrpcapp

import (
	"fmt"
	"log/slog"
	"net"

	gr "github.com/DenisKokorin/Wish-List/internal/group/gRPC"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCserver *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, group gr.Group) *App {
	gRPCserver := grpc.NewServer()
	gr.Register(gRPCserver, group)

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
	log := a.log.With(slog.String("op", "grpcGRapp.Run"), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", "grpcGRapp.Run", err)
	}

	log.Info("GRgrpc server is running", slog.String("addr", l.Addr().String()))
	if err := a.gRPCserver.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", "grpcGRapp.Run", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping WLgrpc server")
	a.gRPCserver.GracefulStop()
}
