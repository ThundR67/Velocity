package main

import (
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/gateway/auth-server/router"
	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
)

//Loding Logger
var log = logger.GetLogger("auth_server.log")

var waitGroup sync.WaitGroup

func main() {
	for _, ipAddress := range config.AuthServerConfigIPAddresses {
		waitGroup.Add(1)
		go func(addr string) {
			runServer(addr)
		}(ipAddress)
	}
	waitGroup.Wait()
}

//runServer runs a single instance of the auth server
func runServer(address string) {
	router := router.GetRouter()
	http.Handle("/", router)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(
			"Auth Server Returned Error While Listening And Serving",
			zap.String("IpAddress", address),
			zap.Error(err),
		)
	}
}
