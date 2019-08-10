package main

import (
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/SonicRoshan/Velocity/services/email-verification-srv/handler"
	proto "github.com/SonicRoshan/Velocity/services/email-verification-srv/proto"
	micro "github.com/micro/go-micro"
	"go.uber.org/zap"
)

func main() {
	//Loading logger
	log := logger.GetLogger("email_verification_service.log")

	//Handling a panic
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Service Paniced Due", zap.Any("Panic", r))
		}
	}()

	service := micro.NewService(
		micro.Name(config.EmailVerificationSrv),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
	)

	service.Init()

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
