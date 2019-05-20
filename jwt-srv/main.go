package main

import (
	logger "github.com/SonicRoshan/Velocity/global/logs"
	proto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	jwtservice "github.com/SonicRoshan/Velocity/jwt-srv/service"
	micro "github.com/micro/go-micro"
)

func main() {
	//Loading logger
	log := logger.GetLogger("jwt_service_log.log")

	//Handling a panic
	defer func() {
		if r := recover(); r != nil {
			log.Criticalf("Service Paniced Due To Reason %s", r)
		}
	}()

	//Creating Service
	service := micro.NewService(
		micro.Name("jwt-srv"),
	)

	service.Init()

	//Initialising service handler
	jwtService := jwtservice.Service{}
	jwtService.Init()

	//Registering Service
	proto.RegisterJWTManagerHandler(service.Server(), jwtService)

	// Run the server
	if err := service.Run(); err != nil {
		log.Criticalf("Service Failed With Error %s", err)
	}
}
