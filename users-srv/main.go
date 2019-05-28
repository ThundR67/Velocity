package main

import (
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/SonicRoshan/Velocity/users-srv/handler"
	proto "github.com/SonicRoshan/Velocity/users-srv/proto"
	micro "github.com/micro/go-micro"
)

func main() {
	//Loading logger
	log := logger.GetLogger("users_service.log")

	//Handling a panic
	defer func() {
		if r := recover(); r != nil {
			log.Criticalf("Service Paniced Due To Reason %s", r)
		}
	}()

	service := micro.NewService(
		micro.Name(config.UsersService),
	)

	service.Init()

	usersService := handler.UsersService{}
	err := usersService.Init()

	if err != nil {
		msg := "Not Able To Connect To MongoDB Server Due To Error " + err.Error()
		log.Critical(msg)
		panic(msg)
	}

	proto.RegisterUsersHandler(service.Server(), usersService)

	if err := service.Run(); err != nil {
		log.Criticalf("Service Failed With Error %s", err)
	}
}
