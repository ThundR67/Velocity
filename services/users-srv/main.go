package main

import (
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/SonicRoshan/Velocity/services/users-srv/handler"
	proto "github.com/SonicRoshan/Velocity/services/users-srv/proto"
	micro "github.com/micro/go-micro"
	"go.uber.org/zap"
)

func main() {
	//Loading logger
	log := logger.GetLogger("users_service.log")

	//Handling a panic
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Service Paniced Due", zap.Any("Panic", r))
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
		log.Fatal(msg)
		panic(msg)
	}

	proto.RegisterUsersHandler(service.Server(), usersService)

	if err := service.Run(); err != nil {
		log.Fatal("Service Failed With Error", zap.Error(err))
	}
}
