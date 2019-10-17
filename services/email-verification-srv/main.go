package main

import (
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
	"github.com/SonicRoshan/Velocity/global/utils"
	"github.com/SonicRoshan/Velocity/services/email-verification-srv/handler"
	proto "github.com/SonicRoshan/Velocity/services/email-verification-srv/proto"
)

func main() {
	//Loading logger
	log := logger.GetLogger("email_verification_service.log")

	//Handling a panic
	defer utils.HandlePanic(log)

	service := utils.CreateService(config.EmailVerificationSrv)

	emailVerificationService := handler.EmailVerification{}
	err := emailVerificationService.Init()

	if err != nil {
		msg := "Not Able To Initialize Handler " + err.Error()
		log.Fatal(msg)
		panic(msg)
	}

	proto.RegisterEmailVerificationHandler(service.Server(), emailVerificationService)

	if err := service.Run(); err != nil {
		log.Fatal("Service Failed With Error", zap.Error(err))
	}
}
