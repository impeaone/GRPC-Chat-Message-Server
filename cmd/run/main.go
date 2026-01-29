package main

import (
	"GrpcMessangerMsgServer/internal/app"
	config2 "GrpcMessangerMsgServer/pkg/config"
	"GrpcMessangerMsgServer/pkg/logger/logger"
	"fmt"
)

func main() {

	logs := logger.NewLog()
	StoragePath := ""

	config, err := config2.ReadConfig()
	if err != nil {
		logs.Error(fmt.Sprintf("Reading config file error: %v", err), logger.GetPlace())
		return
	}
	application := app.NewApp(config.Port, StoragePath, logs)
	// стартуем
	application.GRPCServer.MustRun()
}
