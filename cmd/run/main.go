package main

import (
	"GrpcMessangerMsgServer/internal/app"
	config2 "GrpcMessangerMsgServer/pkg/config"
	"GrpcMessangerMsgServer/pkg/logger/logger"
	"fmt"
)

func main() {
	//TODO: какую-нибудь бд для действий пользователей или т.д. хз кароче

	logs := logger.NewLog()
	StoragePath := "" // TODO: хз что это должно было быть

	//Файл с конфигом, но там маловато, конечно
	config, err := config2.ReadConfig()
	if err != nil {
		logs.Error(fmt.Sprintf("Reading config file error: %v", err), logger.GetPlace())
		return
	}
	//логи все глубже
	application := app.NewApp(config.Port, StoragePath, logs)
	// стартуем
	application.GRPCServer.MustRun()
}
