package main

import (
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
	"github.com/SonicRoshan/Velocity/global/utils"
	"github.com/SonicRoshan/Velocity/services/users-srv/handler"
	proto "github.com/SonicRoshan/Velocity/services/users-srv/proto"
)

var log = logger.GetLogger("users_service.log")

func main() {

	defer utils.HandlePanic(log)

	service := utils.CreateService(config.UsersService)

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
