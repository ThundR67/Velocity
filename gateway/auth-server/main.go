package main

import (
	"net/http"
	"sync"

	"github.com/SonicRoshan/Velocity/gateway/auth-server/router"
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"go.uber.org/zap"
)

//Loding Logger
var log = logger.GetLogger("auth_server.log")

var waitGroup sync.WaitGroup

func main() {
	router := router.GetRouter()

	http.Handle("/", router)

	//Running Auth Server On All Addresses Provided
	for _, ipAddress := range config.AuthServerConfigIPAddresses {
		waitGroup.Add(1)
		go func(address string) {
			log.Fatal(
				"Auth Server Returned Error While Listening And Serving",
				zap.String("IpAddress", address),
				zap.Error(http.ListenAndServe(address, nil)),
			)
			waitGroup.Done()
		}(ipAddress)
	}

	waitGroup.Wait()

}
