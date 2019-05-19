package main

import (
	proto "github.com/SonicRoshan/Velocity/user-data-srv/proto"
	userDataService "github.com/SonicRoshan/Velocity/user-data-srv/service"
	logger "github.com/jex-lin/golang-logger"
	micro "github.com/micro/go-micro"
)

func main() {

	//Log Loding Logger
	var Log = logger.NewLogFile("logs/main.log")

	//Creating Service
	service := micro.NewService(
		micro.Name("user-data-srv"),
	)

	service.Init()

	//Initialising service handler
	userDataService := userDataService.UserDataService{}
	err := userDataService.Init()
	if err != nil {
		msg := "Not Able To Connect To MongoDB Server Due To Error " + err.Error()
		Log.Critical(msg)
		//Panicing As Connection To MongoDB Failed
		panic(msg)
	}

	//Registering Service
	proto.RegisterUserDataManagerHandler(service.Server(), userDataService)

	// Run the server
	if err := service.Run(); err != nil {
		Log.Criticalf("Service Failed With Error %s", err)
	}
}
