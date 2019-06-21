package main

import (
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	proto "github.com/SonicRoshan/Velocity/services/jwt-srv/proto"
	handler "github.com/SonicRoshan/Velocity/services/jwt-srv/service-handler"
	micro "github.com/micro/go-micro"
	"go.uber.org/zap"
)

func main() {
	log := logger.GetLogger("jwt_service.log")

	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Service Paniced", zap.Any("Panic", r))
		}
	}()

	service := micro.NewService(
		micro.Name(config.JWTService),
	)

	service.Init()

	serviceHandler := handler.ServiceHandler{}
	serviceHandler.Init()

	proto.RegisterJWTHandler(service.Server(), serviceHandler)

	if err := service.Run(); err != nil {
		log.Fatal("Service Failed With Error", zap.Error(err))
	}
}
