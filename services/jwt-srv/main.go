package main

import (
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
	"github.com/SonicRoshan/Velocity/global/utils"

	proto "github.com/SonicRoshan/Velocity/services/jwt-srv/proto"
	handler "github.com/SonicRoshan/Velocity/services/jwt-srv/service-handler"
)

var log = logger.GetLogger("jwt_service.log")

func main() {

	defer utils.HandlePanic(log)

	service := utils.CreateService(config.JWTService)

	serviceHandler := handler.ServiceHandler{}
	serviceHandler.Init()

	proto.RegisterJWTHandler(service.Server(), serviceHandler)

	if err := service.Run(); err != nil {
		log.Fatal("Service Failed With Error", zap.Error(err))
	}
}
