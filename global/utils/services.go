package utils

import (
	"time"

	micro "github.com/micro/go-micro"
	"go.uber.org/zap"
)

//HandlePanic is used to handle panic when it occurs in a microservice
func HandlePanic(log *zap.Logger) {
	if r := recover(); r != nil {
		log.Fatal("Service Paniced", zap.Any("Panic", r))
	}
}

//CreateService is used to create a service defination
func CreateService(name string) micro.Service {
	service := micro.NewService(
		micro.Name(name),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
	)

	service.Init()

	return service
}
