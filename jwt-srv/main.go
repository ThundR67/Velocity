package main

import (
	proto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	jwtservice "github.com/SonicRoshan/Velocity/jwt-srv/service"
	logger "github.com/jex-lin/golang-logger"
	micro "github.com/micro/go-micro"
)

func main() {

	//Log Loding Logger
	var log = logger.NewLogFile("logs/main.log")

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
