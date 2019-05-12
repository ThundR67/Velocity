package main

import (
	proto "github.com/SonicRoshan/Velocity/user-data-srv/proto"
	userDataManager "github.com/SonicRoshan/Velocity/user-data-srv/user-data-manager"
	logger "github.com/jex-lin/golang-logger"
	micro "github.com/micro/go-micro"
)

func main() {

	//Log Loding Logger
	var Log = logger.NewLogFile("logs/main.log")

	//Creating Service
	service := micro.NewService(
		micro.Name("Velocity.basic-auth-srv"),
	)

	service.Init()

	//Initialising service handler
	userDataService := userDataManager.UserDataService{}
	userDataService.Init()

	//Registering Service
	proto.RegisterUserDataManagerHandler(service.Server(), userDataService)

	// Run the server
	if err := service.Run(); err != nil {
		Log.Criticalf("Service Failed With Error %s", err)
	}
}
