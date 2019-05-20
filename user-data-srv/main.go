package main

import (
	logger "github.com/SonicRoshan/Velocity/global/logs"
	proto "github.com/SonicRoshan/Velocity/user-data-srv/proto"
	userDataService "github.com/SonicRoshan/Velocity/user-data-srv/service"
	micro "github.com/micro/go-micro"
)

func main() {
	//Loading logger
	log := logger.GetLogger("user_data_service_log.log")

	//Handling a panic
	defer func() {
		if r := recover(); r != nil {
			log.Criticalf("Service Paniced Due To Reason %s", r)
		}
	}()

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
		log.Critical(msg)
		//Panicing As Connection To MongoDB Failed
		panic(msg)
	}

	//Registering Service
	proto.RegisterUserDataManagerHandler(service.Server(), userDataService)

	// Run the server
	if err := service.Run(); err != nil {
		log.Criticalf("Service Failed With Error %s", err)
	}
}
