package main

import (
	"net/http"

	"github.com/SonicRoshan/Velocity/gateway/auth-server/router"
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
)

//Loding Logger
var log = logger.GetLogger("auth_server.log")

func main() {
	router := router.GetRouter()
	http.Handle("/", router)
	log.Critical(http.ListenAndServe(config.AuthServerConfigIPAddress, nil))
}
