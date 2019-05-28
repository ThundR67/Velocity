package main

import (
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	proto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	handler "github.com/SonicRoshan/Velocity/jwt-srv/service-handler"
	micro "github.com/micro/go-micro"
)

func main() {
	log := logger.GetLogger("jwt_service.log")

	defer func() {
		if r := recover(); r != nil {
			log.Criticalf("Service Paniced Due To Reason %s", r)
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
		log.Criticalf("Service Failed With Error %s", err)
	}
}
