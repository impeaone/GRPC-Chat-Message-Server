package grpcapp

import (
	authgrpc "GrpcMessangerMsgServer/internal/grpc/auth"
	"GrpcMessangerMsgServer/pkg/logger/logger"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
	logs       *logger.Log
}

func NewApp(port int, logs *logger.Log) *App {
	GRPCServer := grpc.NewServer()
	authgrpc.Register(GRPCServer, logs)
	return &App{GRPCServer, port, logs}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		app.logs.Error(fmt.Sprintf("App run error: %s", err.Error()), logger.GetPlace())
		return
	}
}

func (app *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		app.logs.Error(fmt.Sprintf("Failed to listen: %s", err.Error()), logger.GetPlace())
		return err
	}
	app.logs.Info(fmt.Sprintf("Server listening on port %d", app.port), logger.GetPlace())

	if errServe := app.gRPCServer.Serve(l); errServe != nil {
		app.logs.Error(fmt.Sprintf("Failed to serve: %s", errServe.Error()), logger.GetPlace())
		return errServe
	}
	app.logs.Info("App is running", logger.GetPlace())
	return nil
}

func (app *App) Stop() {
	app.gRPCServer.GracefulStop()
	app.logs.Info("App is stopping...", logger.GetPlace())
	return
}
