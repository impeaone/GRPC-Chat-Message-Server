package app

import (
	grpcapp "GrpcMessangerMsgServer/internal/app/grpc"
	"GrpcMessangerMsgServer/pkg/logger/logger"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(grpcPort int, storagepath string, logs *logger.Log) *App {
	//TODO : логи все внутри
	grpcAPP := grpcapp.NewApp(grpcPort, logs)
	return &App{
		GRPCServer: grpcAPP,
	}
}
